package datapermctx

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/suyuan32/simple-admin-common/config"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type DataPermKey string

const (
	// ScopeKey is the key to store data scope
	ScopeKey DataPermKey = "data-perm-scope"

	// CustomDeptKey is the key to store custom department ids
	CustomDeptKey DataPermKey = "data-perm-custom-dept"

	// FilterFieldKey is the key to store filter field
	FilterFieldKey DataPermKey = "data-perm-filter-field"
)

// WithScopeContext returns context with data scope
func WithScopeContext(ctx context.Context, scope string) context.Context {
	ctx = metadata.AppendToOutgoingContext(ctx, string(ScopeKey), scope)
	ctx = context.WithValue(ctx, ScopeKey, scope)
	return ctx
}

// GetScopeFromCtx returns data scope from context.
func GetScopeFromCtx(ctx context.Context) (uint8, error) {
	var scope string
	var ok bool

	if scope, ok = ctx.Value(ScopeKey).(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			logx.Error("failed to get data scope from context", logx.Field("detail", ctx))
			return 0, errorx.NewInvalidArgumentError("failed to get data scope")
		} else {
			if data := md.Get(string(ScopeKey)); len(data) > 0 {
				scope = data[0]
			} else {
				return 0, errorx.NewInvalidArgumentError("failed to get data scope")
			}
		}
	}

	id, err := strconv.Atoi(scope)
	if err != nil {
		logx.Error("failed to convert data scope", logx.Field("detail", err))
		return 0, errorx.NewInvalidArgumentError("failed to get data scope")
	}
	return uint8(id), nil
}

// WithCustomDeptContext returns context with custom department ids
func WithCustomDeptContext(ctx context.Context, deptIds string) context.Context {
	ctx = metadata.AppendToOutgoingContext(ctx, string(CustomDeptKey), deptIds)
	ctx = context.WithValue(ctx, CustomDeptKey, deptIds)
	return ctx
}

// GetCustomDeptFromCtx returns custom department ids from context
func GetCustomDeptFromCtx(ctx context.Context) ([]uint64, error) {
	var customDept string
	var ok bool
	var customDeptIds []uint64

	if customDept, ok = ctx.Value(CustomDeptKey).(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			logx.Error("failed to get custom departmrnt ids from context", logx.Field("detail", ctx))
			return nil, errorx.NewInvalidArgumentError("failed to get custom departmrnt ids")
		} else {
			if data := md.Get(string(CustomDeptKey)); len(data) > 0 {
				customDept = data[0]
			} else {
				return nil, errorx.NewInvalidArgumentError("failed to get custom departmrnt ids")
			}
		}
	}

	for _, v := range strings.Split(customDept, ",") {
		id, err := strconv.Atoi(v)
		if err != nil {
			logx.Error("failed to convert custom departmrnt ids", logx.Field("detail", err), logx.Field("data", v))
			return nil, errorx.NewInvalidArgumentError("failed to get custom departmrnt ids")
		}
		customDeptIds = append(customDeptIds, uint64(id))
	}

	return customDeptIds, nil
}

// WithFilterFieldContext returns context with filter field
func WithFilterFieldContext(ctx context.Context, filterField string) context.Context {
	ctx = metadata.AppendToOutgoingContext(ctx, string(FilterFieldKey), filterField)
	ctx = context.WithValue(ctx, FilterFieldKey, filterField)
	return ctx
}

// GetFilterFieldFromCtx returns filter field from context
func GetFilterFieldFromCtx(ctx context.Context) (string, error) {
	if filterField, ok := ctx.Value(FilterFieldKey).(string); !ok {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			logx.Error("failed to get filter field from context", logx.Field("detail", ctx))
			return "", errorx.NewInvalidArgumentError("failed to get filter field")
		} else {
			if data := md.Get(string(FilterFieldKey)); len(data) > 0 {
				return data[0], nil
			} else {
				return "", errorx.NewInvalidArgumentError("failed to get filter field")
			}
		}
	} else {
		return filterField, nil
	}
}

// GetRoleCustomDeptDataPermRedisKey returns the key to store role custom department data into redis
func GetRoleCustomDeptDataPermRedisKey(roleCodes []string) string {
	return fmt.Sprintf("%s:ROLE:%s:CustomDept", config.RedisDataPermissionPrefix, strings.Join(roleCodes, ","))
}

// GetSubDeptDataPermRedisKey returns the key to store sub department data into redis
func GetSubDeptDataPermRedisKey(departmentId uint64) string {
	return fmt.Sprintf("%s:DEPT:%d:SubDept", config.RedisDataPermissionPrefix, departmentId)
}
