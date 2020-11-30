package router

import (
	mosn "mosn.io/mosn/pkg/types"
)

type Matcher interface {
	Match(headers mosn.HeaderMap) bool
}
