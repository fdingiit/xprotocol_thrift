package fastsimplemap

import "unsafe"

// nolint
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

//nolint
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//nolint
func s2b(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
