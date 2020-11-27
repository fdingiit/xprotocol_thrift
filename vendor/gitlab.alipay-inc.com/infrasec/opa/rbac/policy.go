package rbac

import (
	"context"
	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
	"mosn.io/api"
)

type InheritPolicy struct {
	// The set of permissions that define a role.
	// Each permission is matched with OR semantics.
	// To match all actions for this policy, a single Permission with the `any` field set to true should be used.
	InheritPermissions []InheritPermission
	// The set of principals that are assigned/denied the role based on “action”.
	// Each principal is matched with OR semantics.
	// To match all downstreams for this policy, a single Principal with the `any` field set to true should be used.
	InheritPrincipals []InheritPrincipal
}

// Receive the rbactypes.RBACPolicyConf input and convert it to mosn rbac policy
func NewInheritPolicy(policy *rbactypes.RBACPolicyConf) (*InheritPolicy, error) {
	inheritPolicy := new(InheritPolicy)

	// fill permission
	inheritPolicy.InheritPermissions = make([]InheritPermission, len(policy.GetPermissions()))
	for idx, permission := range policy.GetPermissions() {
		if inheritPermission, err := NewInheritPermission(permission); err != nil {
			return nil, err
		} else {
			inheritPolicy.InheritPermissions[idx] = inheritPermission
		}
	}

	// fill principal
	inheritPolicy.InheritPrincipals = make([]InheritPrincipal, len(policy.GetPrincipals()))
	for idx, principal := range policy.GetPrincipals() {
		if inheritPrincipal, err := NewInheritPrincipal(principal); err != nil {
			return nil, err
		} else {
			inheritPolicy.InheritPrincipals[idx] = inheritPrincipal
		}
	}

	return inheritPolicy, nil
}

// A policy matches if and only if at least one of its permissions match the action taking place
// AND at least one of its principals match the downstream.
func (inheritPolicy *InheritPolicy) Match(ctx context.Context, cb api.StreamReceiverFilterHandler, headers api.HeaderMap) bool {
	permissionMatch, principalMatch := false, false
	for _, permission := range inheritPolicy.InheritPermissions {
		if permission.Match(ctx, cb, headers) {
			permissionMatch = true
			break
		}
	}
	if permissionMatch == false {
		return false
	}

	for _, principal := range inheritPolicy.InheritPrincipals {
		if principal.Match(cb, headers) {
			principalMatch = true
			break
		}
	}
	return principalMatch
}
