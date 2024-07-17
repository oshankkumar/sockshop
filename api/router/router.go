package router

import "github.com/oshankkumar/sockshop/api/httpkit"

type Route struct {
	Method  string
	Pattern string
	Handler httpkit.Handler
}

type Router interface {
	Routes() []Route
}

type RouterFunc func() []Route

func (r RouterFunc) Routes() []Route {
	return r()
}

func ComposeRouters(rr ...Router) Router {
	return RouterFunc(func() []Route {
		var rtt []Route
		for _, rt := range rr {
			rtt = append(rtt, rt.Routes()...)
		}
		return rtt
	})
}
