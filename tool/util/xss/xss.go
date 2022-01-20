package xss

import (
	"github.com/dreamlu/gt/tool/type/cmap"
	"html"
)

func XssMap(args cmap.CMap) {
	for _, v := range args {
		v[0] = html.EscapeString(v[0])
	}
	return
}
