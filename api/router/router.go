package router

import (
	"net/http"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
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
