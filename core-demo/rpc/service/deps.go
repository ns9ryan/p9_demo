package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	ModeOff = "off"
	ModeOn  = "on"
)

const RoleSuperAdmin = "super_admin"
const blacklistPrefix = "rbacx:jwt:bl:"

type Deps struct {
	Client           *ent.Client
	Enforcer         *casbin.Enforcer
	Redis            redis.UniversalClient
	Mode             string
	JWTSecret        string
	JWTExpire        int64
	JWTRefreshSecret string
	JWTRefreshExpire int64
	APIPrefix        string
	InitToken        string
}

func RandomSalt() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashPassword(pw string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h), err
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func NewUserCode() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (d *Deps) DomainFromClaims(c *ctxdata.Claims) string {
	if d.Mode == ModeOff {
		return ""
	}
	if c == nil {
		return ""
	}
	return casbinx.DomainID(c.OperatorID)
}

func (d *Deps) RequireMode(want string) error {
	if d.Mode != want {
		return xerr.BadRequest(i18n.AuthPartnerModeMismatch)
	}
	return nil
}

func (d *Deps) ActiveUserByID(ctx context.Context, id int64) (*model.User, error) {
	u, err := d.Client.User.Query().Where(user.ID(id), user.DeletedAtIsNil()).Only(ctx)
	if err != nil {
		return nil, err
	}
	return userFromEnt(u), nil
}

func (d *Deps) RoleCodesOfUser(ctx context.Context, userID int64) ([]string, error) {
	roles, err := d.Client.User.Query().Where(user.ID(userID)).QueryRoles().
		Where(role.DeletedAtIsNil(), role.StatusEQ(model.StatusNormal)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(roles))
	for _, r := range roles {
		codes = append(codes, r.RoleCode)
	}
	return codes, nil
}

func (d *Deps) RotateSalt(ctx context.Context, userID int64) error {
	salt, err := RandomSalt()
	if err != nil {
		return err
	}
	n, err := d.Client.User.Update().
		Where(user.ID(userID), user.DeletedAtIsNil()).
		SetSalt(salt).
		Save(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return xerr.NotFound(i18n.UserNotFound)
	}
	return nil
}

func (d *Deps) BlacklistToken(ctx context.Context, token string, expUnix int64) error {
	if d.Redis == nil || token == "" {
		return nil
	}
	ttl := time.Until(time.Unix(expUnix, 0))
	if ttl <= 0 {
		ttl = time.Minute
	}
	return d.Redis.Set(ctx, blacklistPrefix+token, "1", ttl).Err()
}

func (d *Deps) TokenBlacklisted(ctx context.Context, token string) bool {
	if d.Redis == nil || token == "" {
		return false
	}
	n, err := d.Redis.Exists(ctx, blacklistPrefix+token).Result()
	return err == nil && n > 0
}

func (d *Deps) SignTokenPair(ctx context.Context, u *model.User, roleCodes []string) (TokenInfo, error) {
	base := d.tokenClaims(ctx, u, roleCodes)
	access, accessExp, err := jwt.Sign(d.JWTSecret, d.JWTExpire, base.WithType(jwt.TokenAccess))
	if err != nil {
		return TokenInfo{}, err
	}
	refresh, refreshExp, err := jwt.Sign(d.JWTRefreshSecret, d.JWTRefreshExpire, base.WithType(jwt.TokenRefresh))
	if err != nil {
		return TokenInfo{}, err
	}
	return TokenInfo{
		AccessToken:   access,
		RefreshToken:  refresh,
		Expire:        accessExp,
		RefreshExpire: refreshExp,
	}, nil
}

func (d *Deps) tokenClaims(ctx context.Context, u *model.User, roleCodes []string) jwt.Claims {
	oid := int64(0)
	code := ""
	if u.OperatorID != nil {
		oid = *u.OperatorID
		if op, err := d.Client.Operator.Get(ctx, oid); err == nil {
			code = op.OperatorCode
		}
	}
	return jwt.Claims{
		UserID:       u.ID,
		UserCode:     u.UserCode,
		Username:     u.Username,
		OperatorID:   oid,
		OperatorCode: code,
		RoleCodes:    roleCodes,
		Salt:         u.Salt,
	}
}

func (d *Deps) AllAPIPolicies(ctx context.Context, domain string) ([][]string, error) {
	apis, err := d.Client.API.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var policies [][]string
	for _, a := range apis {
		if d.Mode == ModeOff && isOperatorAPI(a.Path) {
			continue
		}
		policies = append(policies, []string{RoleSuperAdmin, domain, a.Path, a.Method})
	}
	return policies, nil
}

func isOperatorAPI(path string) bool {
	return strings.HasSuffix(path, "/operator/self") || strings.HasSuffix(path, "/operator/update")
}

func (d *Deps) ReplaceRoleAPIPolicies(roleCode, domain string, policies [][]string) error {
	if _, err := d.Enforcer.RemoveFilteredPolicy(0, roleCode, domain); err != nil {
		return err
	}
	if len(policies) == 0 {
		return nil
	}
	_, err := d.Enforcer.AddPolicies(policies)
	return err
}
