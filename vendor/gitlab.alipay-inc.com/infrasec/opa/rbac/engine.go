package rbac

import (
	"context"
	"fmt"
	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"mosn.io/api"
)

const (
	NotInWhitePolicy = "Not In Allowed Policies"
)

type RoleBasedAccessControlEngine struct {
	// The request is allowed if and only if:
	//   * `action` is "ALLOWED" and at least one policy matches
	//   * `action` is "DENY" and none of the policies match
	// default is ALLOWED
	Action rbactypes.RBAC_Action

	// Maps from policy name to policy. A match occurs when at least one policy matches the request.
	InheritPolicies map[string]*InheritPolicy

	// The services which must use mTLS to invoke
	StrictMTLS *InheritStrictMTLS
}

// Receive the rbactypes.RBAC input and convert it to mosn rbac engine
func NewRoleBasedAccessControlEngine(rbacConfig *rbactypes.RBAC) (*RoleBasedAccessControlEngine, error) {
	engine := new(RoleBasedAccessControlEngine)

	// fill engine action, `RBAC_DENY` by default
	engine.Action = rbacConfig.GetAction()

	// fill engine policies
	engine.InheritPolicies = make(map[string]*InheritPolicy)
	for name, policy := range rbacConfig.GetPolicies() {
		if inheritPolicy, err := NewInheritPolicy(policy); err != nil {
			// skip to the next policy
			continue
		} else {
			engine.InheritPolicies[name] = inheritPolicy
		}
	}

	// fill strict mTLS policy
	inheritStrictMTLS, err := NewInheritStrictMTLS(rbacConfig.GetStrictMtls())
	if err != nil {
		return nil, err
	}
	engine.StrictMTLS = inheritStrictMTLS

	return engine, nil
}

// echo request will be handled in `Allowed` function
func (engine *RoleBasedAccessControlEngine) Allowed(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) (allowed bool, matchPolicyName string) {
	defer func() {
		if err := recover(); err != nil {
			// defer runs after the return statement but before the function is actually returned,
			// so we can use named return values to hack function return
			allowed, matchPolicyName = true, ""
		}
	}()

	// do the strict mTLS check first
	mTLSAllowed, serviceName := engine.StrictMTLS.Allowed(cb, headers)
	if !mTLSAllowed {
		return false, fmt.Sprintf("mTLS check policy, service name: %s", serviceName)
	}

	// If policies is empty, return allowed
	if len(engine.InheritPolicies) == 0 {
		return true, ""
	}

	if engine.Action == rbactypes.RBAC_ALLOW {
		// when engine action is ALLOW, return allowed if matched any policy
		for _, policy := range engine.InheritPolicies {
			if policy.Match(ctx, cb, headers) {
				return true, ""
			}
		}
		return false, NotInWhitePolicy
	} else if engine.Action == rbactypes.RBAC_DENY {
		// when engine action is DENY, return allowed if not matched any policy
		for name, policy := range engine.InheritPolicies {
			if policy.Match(ctx, cb, headers) {
				return false, name
			}
		}
		return true, ""
	}

	return true, ""
}
