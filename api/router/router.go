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

type Routers []Router

func (r Routers) Routes() []Route {
	var rr []Route
	for _, rt := range r {
		rr = append(rr, rt.Routes()...)
	}
	return rr
}
