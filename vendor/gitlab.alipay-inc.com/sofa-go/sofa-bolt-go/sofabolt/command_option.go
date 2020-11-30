package sofabolt

type ReadOption struct{}

func NewReadOption() *ReadOption { return &ReadOption{} }

type WriteOption struct{}

func NewWriteOption() *WriteOption { return &WriteOption{} }
