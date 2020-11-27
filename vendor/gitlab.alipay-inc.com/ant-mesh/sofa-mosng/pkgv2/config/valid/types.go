package valid

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
)

type Validation interface {
	ResourceType() api.Type
	Valid(api.Object) (error, bool)
}

type ValidationManager interface {
	Register(validation Validation)
	GetValidations(t api.Type) []Validation
	DoValid(api.Object) (error, bool)
}

type ValidationFunc struct {
	Type    api.Type
	DoValid func(api.Object) (error, bool)
}

func (v *ValidationFunc) ResourceType() api.Type {
	return v.Type
}

func (v *ValidationFunc) Valid(o api.Object) (error, bool) {
	if v.DoValid == nil {
		return errors.Error(" not found valid method"), false
	}
	return v.DoValid(o)
}
