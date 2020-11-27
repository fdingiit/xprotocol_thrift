package valid

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
)

func init() {
	ValidationManagerInstance().Register(&routeValid{})
	ValidationManagerInstance().Register(&routeGroupValid{})
}

type routeValid struct{}

func (*routeValid) ResourceType() api.Type {
	return api.ROUTER
}

func (*routeValid) Valid(o api.Object) (error, bool) {

	if router, ok := o.(*v1.Router); ok {

		if router.GroupBind == "" {
			return errors.Error("no group name"), false
		}

	} else {
		return errors.Error("cast exception"), false
	}

	return nil, true
}

type routeGroupValid struct{}

func (*routeGroupValid) ResourceType() api.Type {
	return api.ROUTER_GROUP
}

func (*routeGroupValid) Valid(o api.Object) (error, bool) {
	if rg, ok := o.(*v1.RouterGroup); ok {
		for _, r := range rg.Routers {
			r.GroupBind = rg.GroupName
		}
		// todo more
	} else {
		return errors.Error("cast exception"), false
	}

	return nil, true
}
