package metadata

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"strconv"
	"strings"
)

func init() {
	event.EventListenerManagerInstance().Register(event.ResourceEventListenerFuncs{
		Type: api.METADATA,
		AddFunc: func(o api.Object) (e error, b bool) {
			config := o.(v1.Metadata)
			AddOrUpdateConfigTree(config.Key, config.Value, nil)

			return nil, true
		},
		UpdateFunc: func(o api.Object) (e error, b bool) {
			config := o.(*v1.Metadata)
			AddOrUpdateConfigTree(config.Key, config.Value, nil)

			return nil, true
		},
	})
}

func Get(configPath string) (interface{}, *node) {
	if strings.Index(configPath, "$") == 0 {
		configPath = strings.Replace(configPath, "[", ".", -1)
		configPath = strings.Replace(configPath, "]", "", -1)
		split := strings.Split(configPath, ".")
		var node = GetConfigTreeInstance().root
		for i := 1; i < len(split); i++ {
			node = node.Get(split[i])
		}
		return node.Value(), node
	}

	return configPath, nil
}

// todo 干掉递归
func AddOrUpdateConfigTree(key string, conf interface{}, n *node) {
	if n == nil {
		n = GetConfigTreeInstance().root
	}

	c, _ := addOrUpdateNode(n, key, conf)

	if m, ok := conf.(map[interface{}]interface{}); ok {

		for k, v := range m {
			if a, ok := v.([]interface{}); ok {
				AddOrUpdateArrayTree(k.(string), a, c)
			} else {
				AddOrUpdateConfigTree(k.(string), v, c)
			}
		}

	} else if a, ok := conf.([]interface{}); ok {
		AddOrUpdateArrayTree(key, a, n)
	} else {
		if c, ok := hasKey(n, key); ok {
			c.set(conf)
		} else {
			n.addChildren(newNode(key, conf))
		}

	}
}

func AddOrUpdateArrayTree(key string, arr []interface{}, n *node) {
	c, _ := addOrUpdateNode(n, key, arr)

	for i, a := range arr {
		addOrUpdateArray(i, a, c)
	}
}

func addOrUpdateArray(index int, conf interface{}, n *node) {
	c, _ := addOrUpdateNode(n, strconv.Itoa(index), conf)
	if m, ok := conf.(map[interface{}]interface{}); ok {
		// obj( k-v )
		for k, v := range m {
			AddOrUpdateConfigTree(k.(string), v, c)
		}
	} else if a, ok := conf.([]interface{}); ok {
		AddOrUpdateArrayTree(strconv.Itoa(index), a, c)
	} else {
		// v
		if c, ok = hasKey(n, strconv.Itoa(index)); ok {
			c.set(conf)
		} else {
			n.addChildren(newNode(strconv.Itoa(index), conf))
		}
	}

}

func addOrUpdateNode(n *node, key string, conf interface{}) (*node, bool) {
	var c *node
	var ok bool
	if c, ok = hasKey(n, key); ok {
		c.set(conf)
	} else {
		c = newNode(key, conf)
		n.addChildren(c)
		c.parent = n
	}
	return c, ok
}

func hasKey(node *node, key string) (*node, bool) {
	for _, c := range node.children {
		if c.key == key {
			return c, true
		}
	}

	return nil, false
}
