package sofahessian

type EncodeTypeRefs struct {
	refs []string
}

func NewEncodeTyperefs() *EncodeTypeRefs {
	return &EncodeTypeRefs{}
}

func (e *EncodeTypeRefs) Reset() {
	e.refs = e.refs[:0]
}

func (e *EncodeTypeRefs) Get(typ string) (int, bool) {
	for i := range e.refs {
		if e.refs[i] == typ {
			return i, true
		}
	}
	return -1, false
}

func (e *EncodeTypeRefs) Set(typ string) {
	for i := range e.refs {
		if e.refs[i] == typ {
			return
		}
	}
	e.refs = append(e.refs, typ)
}

type EncodeClassrefs struct {
	refs []interface{}
}

func NewEncodeClassrefs() *EncodeClassrefs {
	return &EncodeClassrefs{}
}

func (c *EncodeClassrefs) Reset() {
	c.refs = c.refs[:0]
}

func (c *EncodeClassrefs) Add(cls interface{}) (ref int, referenced bool) {
	for i := range c.refs {
		if cls == c.refs[i] {
			return i, true
		}
	}
	n := len(c.refs)
	c.refs = append(c.refs[:n], cls)
	return n, false
}

func (c *EncodeClassrefs) Get(cls interface{}) int {
	for i := range c.refs {
		if c.refs[i] == cls {
			return i
		}
	}
	return -1
}

type EncodeObjectrefs struct {
	refs []interface{}
}

func NewEncodeObjectrefs() *EncodeObjectrefs {
	return &EncodeObjectrefs{}
}

func (c *EncodeObjectrefs) Reset() {
	c.refs = c.refs[:0]
}

func (c *EncodeObjectrefs) Add(obj interface{}) int {
	for i := range c.refs {
		if safeEqual(c.refs[i], obj) {
			c.refs[i] = obj
			return i

		}
	}
	c.refs = append(c.refs, obj)
	return len(c.refs) - 1
}

func (c *EncodeObjectrefs) Get(obj interface{}) int {
	for i := range c.refs {
		if safeEqual(obj, c.refs[i]) {
			return i
		}
	}
	return -1
}

type DecodeObjectRefs struct {
	refs []interface{}
}

func NewDecodeObjectRefs() *DecodeObjectRefs {
	return &DecodeObjectRefs{}
}

func (t *DecodeObjectRefs) Reset() {
	t.refs = t.refs[:0]
}

func (t DecodeObjectRefs) Len() int {
	return len(t.refs)
}

func (t DecodeObjectRefs) Get(id int) (interface{}, bool) {
	if id < 0 || id >= len(t.refs) {
		return "", false
	}

	return t.refs[id], true
}

func (t *DecodeObjectRefs) Append(obj interface{}) {
	t.refs = append(t.refs, obj)
}

type DecodeTypeRefs struct {
	refs []string
}

func NewDecodeTypeRefs() *DecodeTypeRefs {
	return &DecodeTypeRefs{
		refs: []string{},
	}
}

func (t *DecodeTypeRefs) Reset() {
	t.refs = t.refs[:0]
}

func (t DecodeTypeRefs) Len() int {
	return len(t.refs)
}

func (t DecodeTypeRefs) Get(id int) (string, bool) {
	if id < 0 || id >= len(t.refs) {
		return "", false
	}

	return t.refs[id], true
}

func (t *DecodeTypeRefs) Append(typ string) {
	t.refs = append(t.refs, typ)
}

type DecodeClassRefs struct {
	refs []ClassDefinition
}

func NewDecodeClassRefs() *DecodeClassRefs {
	return &DecodeClassRefs{}
}

func (j *DecodeClassRefs) Reset() {
	j.refs = j.refs[:0]
}

func (j *DecodeClassRefs) Get(id int) (ClassDefinition, bool) {
	if id < 0 || id >= len(j.refs) {
		return ClassDefinition{}, false
	}

	return j.refs[id], true
}

func (j *DecodeClassRefs) Append(cf ClassDefinition) {
	j.refs = append(j.refs, cf)
}

type ClassDefinition struct {
	class  string
	fields []string
}
