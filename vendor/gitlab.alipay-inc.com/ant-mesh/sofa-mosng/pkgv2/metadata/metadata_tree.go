package metadata

import (
	"sync"

	"encoding/json"
	"reflect"
)

type Update interface {
	Update(key string, conf interface{}) bool
}

type listener struct {
	Pointer interface{}
	Factory func(conf interface{}) interface{}
}

type node struct {
	key       string
	cMux      sync.RWMutex
	conf      interface{}
	bytes     []byte
	parent    *node
	nMux      sync.RWMutex
	children  []*node
	lMux      sync.RWMutex
	listeners []*listener
}

var emptyNode = &node{
	children: []*node{},
}

func (n *node) Get(key string) *node {
	n.nMux.RLock()
	defer n.nMux.RUnlock()
	for _, c := range n.children {
		if c.key == key {
			return c
		}
	}
	// support n.Get().Get().Get()
	return emptyNode
}

func (n *node) Value() interface{} {
	n.cMux.RLock()
	defer n.cMux.RUnlock()
	return n.conf
}

func (n *node) set(conf interface{}) {
	n.cMux.Lock()
	n.conf = conf
	n.cMux.Unlock()

	bytes, _ := json.Marshal(n.conf)

	if string(bytes) == string(n.bytes) {
		return
	}

	n.bytes = bytes

	n.update()
}

func (n *node) Register(pointer interface{}, factory func(conf interface{}) interface{}) {
	n.lMux.Lock()
	defer n.lMux.Unlock()
	n.listeners = append(n.listeners, &listener{Pointer: pointer, Factory: factory})
}

func (n *node) addChildren(c *node) {
	n.nMux.Lock()
	defer n.nMux.Unlock()
	n.children = append(n.children, c)
}

func (n *node) update() {
	n.lMux.RLock()
	defer n.lMux.RUnlock()

	for _, l := range n.listeners {
		factory := l.Factory(n.conf)
		v := reflect.ValueOf(l.Pointer).Elem()
		v.Set(reflect.ValueOf(factory).Elem())
	}
}

type Tree struct {
	root *node
}

var (
	once sync.Once
	tree *Tree
)

func GetConfigTreeInstance() *Tree {
	once.Do(func() {
		tree = newTree()
	})

	return tree
}

func newTree() *Tree {
	return &Tree{
		newNode("$", nil),
	}
}

func newNode(key string, conf interface{}) *node {
	return &node{
		key:       key,
		conf:      conf,
		children:  []*node{},
		listeners: []*listener{},
		bytes:     []byte{},
	}
}
