package operator

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/operator-base/rpc/ent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InitializeDomain operator 初始化域名
type InitializeDomain struct {
	DomainName string // 域名
	DomainType int64  // 域名类型: 1分站后台, 2代理后台, 3会员H5
}

// normalizeDomains 整理 operator 初始化域名
func normalizeDomains(domains []InitializeDomain) {
	for i := range domains {
		domains[i].DomainName = strings.ToLower(strings.TrimSpace(domains[i].DomainName))
	}
}

// validateDomains 校验 operator 初始化域名
func validateDomains(domains []InitializeDomain) error {
	// 初始化时三种域名必须全部存在
	if len(domains) != 3 {
		return status.Error(codes.InvalidArgument, "three operator domains are required")
	}

	domainTypes := make(map[int64]struct{}, 3)
	domainNames := make(map[string]struct{}, 3)

	for _, item := range domains {
		// 校验域名
		if item.DomainName == "" {
			return status.Error(codes.InvalidArgument, "domain name is required")
		}

		// 校验域名类型
		if item.DomainType < 1 || item.DomainType > 3 {
			return status.Error(codes.InvalidArgument, "domain type is invalid")
		}

		// 每种域名类型只能存在一条
		if _, exists := domainTypes[item.DomainType]; exists {
			return status.Error(codes.InvalidArgument, "duplicate domain type")
		}
		domainTypes[item.DomainType] = struct{}{}

		// 同一个 operator 初始化时不能出现重复域名
		if _, exists := domainNames[item.DomainName]; exists {
			return status.Error(codes.InvalidArgument, "duplicate domain name")
		}
		domainNames[item.DomainName] = struct{}{}
	}

	return nil
}

// createDomains 创建 operator 当前有效域名
func createDomains(ctx context.Context, tx *ent.Tx, operatorID int64, domains []InitializeDomain) error {
	builders := make([]*ent.OperatorDomainCreate, 0, len(domains))

	for _, item := range domains {
		builders = append(builders,
			tx.OperatorDomain.
				Create().
				SetOperatorID(operatorID).      // operator 本地主键
				SetDomainName(item.DomainName). // 域名
				SetDomainType(item.DomainType), // 域名类型
		)
	}

	if err := tx.OperatorDomain.CreateBulk(builders...).Exec(ctx); err != nil {
		return fmt.Errorf("创建 operator 域名失败: %w", err)
	}

	return nil
}
