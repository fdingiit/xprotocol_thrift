package sofahessian

// IdentityIntMap implements the identity map of java in hessian with Golang.
type IdentityIntMap struct {
	m map[interface{}]int
}

func NewIdentityIntMap() *IdentityIntMap {
	return &IdentityIntMap{
		m: make(map[interface{}]int, 256),
	}
}

func (m *IdentityIntMap) Size() int {
	return len(m.m)
}

func (m *IdentityIntMap) Get(key interface{}) int {
	value, ok := m.m[key]
	if ok {
		return value
	}

	return -1
}

func (m *IdentityIntMap) Put(key interface{}, value int, replaced bool) int {
	oldvalue, ok := m.m[key]
	if ok {
		if replaced {
			m.m[key] = value
			return value
		}
		return oldvalue
	}

	m.m[key] = value
	return value
}
