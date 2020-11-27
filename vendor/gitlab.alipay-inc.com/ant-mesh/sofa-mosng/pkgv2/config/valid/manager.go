package valid

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
)

var (
	validationManager ValidationManager = &manager{
		vMap: make(map[api.Type][]Validation),
	}
)

type manager struct {
	vMap map[api.Type][]Validation
}

func (f *manager) Register(v Validation) {
	t := v.ResourceType()
	if vs := f.vMap[t]; vs != nil {
		f.vMap[t] = append(f.vMap[t], v)
	} else {
		f.vMap[t] = []Validation{v}
	}
}

func (f *manager) GetValidations(t api.Type) []Validation {
	return f.vMap[t]
}

func (f *manager) DoValid(o api.Object) (error, bool) {
	for _, v := range f.vMap[o.Type()] {
		if err, ok := v.Valid(o); !ok {
			return err, false
		}
	}

	return nil, true
}

func ValidationManagerInstance() ValidationManager {
	return validationManager
}
