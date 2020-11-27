package model

import (
	"strings"
)

type Commander struct {
	props map[string]string
}

func NewCommander(name string) *Commander {
	props := make(map[string]string)
	idx := strings.Index(name, ":")
	if idx > 0 {
		values := name[(idx + 1):]
		idx = strings.Index(values, ",")

		if idx > 0 {
			dataList := strings.Split(values, ",")
			for _, data := range dataList {
				idx = strings.Index(data, "=")
				if idx != -1 {
					props[data[:idx]] = data[(idx + 1):]
				}
			}
		} else {
			idx = strings.Index(values, "=")
			if idx != -1 {
				props[values[:idx]] = values[(idx + 1):]
			}
		}
	}

	return &Commander{
		props: props,
	}
}

func (c *Commander) GetProps(name string, defaultValue string) string {
	value, ok := c.props[name]
	if ok {
		return value
	}
	return defaultValue
}
