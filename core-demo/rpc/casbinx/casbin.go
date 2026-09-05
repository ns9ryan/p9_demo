package casbinx

import (
	"log"
	"strconv"

	"oa.98ent.com/p9/core/rpc/ent"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
)

const modelText = `
[request_definition]
r = sub, dom, obj, act
[policy_definition]
p = sub, dom, obj, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && r.dom == p.dom && keyMatch2(r.obj,p.obj) && r.act == p.act
`

func New(client *ent.Client, rds redis.UniversalClient) (*casbin.Enforcer, error) {
	adapter := NewAdapter(client)
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	attachWatcher(e, rds)
	return e, nil
}

func attachWatcher(e *casbin.Enforcer, rds redis.UniversalClient) {
	if rds == nil {
		return
	}
	opt, ok := redisOptions(rds)
	if !ok {
		return
	}
	w, err := rediswatcher.NewWatcher(opt.Addr, rediswatcher.WatcherOptions{
		Options:    opt,
		Channel:    "/casbin",
		IgnoreSelf: true,
	})
	if err != nil {
		log.Printf("casbin redis watcher disabled: %v", err)
		return
	}
	_ = e.SetWatcher(w)
	_ = w.SetUpdateCallback(func(string) {
		_ = e.LoadPolicy()
	})
}

func redisOptions(rds redis.UniversalClient) (redis.Options, bool) {
	switch c := rds.(type) {
	case *redis.Client:
		opt := c.Options()
		if opt == nil {
			return redis.Options{}, false
		}
		return *opt, true
	default:
		return redis.Options{}, false
	}
}

func Domain(operatorID *int64) string {
	if operatorID == nil || *operatorID == 0 {
		return ""
	}
	return strconv.FormatInt(*operatorID, 10)
}

func DomainID(id int64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(id, 10)
}
