package queue

type Element interface {
	ID() string
	GetVersion() int64
}

type Queue interface {
	Push(i Element) error
	Pop() (Element, error)
	Len() int
}
