package api

import (
	ewa "github.com/egovorukhin/egowebapi"
)

type Messages struct {
}

func (u *Messages) Get(route *ewa.Route) {
	route.EmptyHandler()
}
