package main

import (
	"errors"
	"strings"
)

type routeNodeOrRouteNodeMap struct {
	child    *RouteNode
	children *map[string]*RouteNode
}

func (r routeNodeOrRouteNodeMap) IsChildOnly() bool {
	return r.child != nil
}

func newChild(child RouteNode) routeNodeOrRouteNodeMap {
	newChild := child
	return routeNodeOrRouteNodeMap{&newChild, nil}
}

func newChildren(children map[string]*RouteNode) routeNodeOrRouteNodeMap {
	newChildren := children
	return routeNodeOrRouteNodeMap{nil, &newChildren}
}

// RouteNode a single ruote node.
type RouteNode struct {
	childOrChildren routeNodeOrRouteNodeMap
	value           interface{}
}

func NewRouteNode(components []string, value interface{}) RouteNode {
	var r RouteNode
	r.value = newChild
	r.childOrChildren = newChildren(make(map[string]*RouteNode))

	r.Add(components, value)
	return r
}

// Add adds a new sub rooute.
func (r *RouteNode) Add(components []string, value interface{}) {
	if len(components) <= 0 {
		r.value = value
	} else {
		first, remainder := components[0], components[1:]
		if first[0] == ':' {
			if r.childOrChildren.IsChildOnly() {
				r.childOrChildren.child.Add(remainder, value)
			} else {
				r.childOrChildren = newChild(NewRouteNode(remainder, value))
			}
		} else {
			if r.childOrChildren.IsChildOnly() {
				r.childOrChildren = newChildren(make(map[string]*RouteNode))
				node := NewRouteNode(remainder, value)
				(*r.childOrChildren.children)[first] = &node
			} else {
				node, ok := (*r.childOrChildren.children)[first]
				if !ok {
					node := NewRouteNode(remainder, value)
					(*r.childOrChildren.children)[first] = &node
				} else {
					node.Add(remainder, value)
				}
			}
		}
	}
}

func (r *RouteNode) Get(components []string) interface{} {
	if len(components) <= 0 {
		return r.value
	}
	first, remainder := components[0], components[1:]
	if r.childOrChildren.IsChildOnly() {
		return r.childOrChildren.child.Get(remainder)
	}
	node, ok := (*r.childOrChildren.children)[first]
	if !ok {
		return nil
	}
	return node.Get(remainder)
}

type PartialRouteNodeResult struct {
	Retrieved bool
	Value     interface{}
	Remainder []string
}

func (r *RouteNode) GetPartial(components []string) PartialRouteNodeResult {
	if len(components) <= 0 {
		return PartialRouteNodeResult{
			Retrieved: true,
			Value:     r.value,
			Remainder: components,
		}
	}
	first, remainder := components[0], components[1:]
	if r.childOrChildren.IsChildOnly() {
		return r.childOrChildren.child.GetPartial(remainder)
	}
	node, ok := (*r.childOrChildren.children)[first]
	if !ok {
		return PartialRouteNodeResult{
			Retrieved: true,
			Value:     r.value,
			Remainder: remainder,
		}
	}
	return node.GetPartial(remainder)
}

// Routes get the routes.
type Routes struct {
	children map[string]RouteNode
}

// Add adds a new route.
func (r *Routes) Add(route string, value interface{}) error {
	if len(route) <= 0 {
		return errors.New("Route cannot be empty")
	}

	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]

	node, ok := r.children[first]
	if !ok {
		r.children[first] = NewRouteNode(remainder, value)
	} else {
		node.Add(remainder, value)
	}

	return nil
}

func (r Routes) Get(route string) interface{} {
	if len(route) <= 0 {
		return nil
	}
	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]

	node, ok := r.children[first]
	if !ok {
		return nil
	}
	return node.Get(remainder)
}

func NewRouter() Routes {
	return Routes{make(map[string]RouteNode)}
}

type PartialRouteResult struct {
	Retrieved bool
	Value     interface{}
	Remainder string
}

func (r Routes) GetPartial(route string) PartialRouteResult {
	if len(route) <= 0 {
		return PartialRouteResult{}
	}
	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]
	node, ok := r.children[first]
	if !ok {
		return PartialRouteResult{
			Retrieved: false,
			Value:     nil,
			Remainder: route,
		}
	}

	result := node.GetPartial(remainder)
	if !result.Retrieved {
		return PartialRouteResult{
			Retrieved: false,
			Value:     nil,
			Remainder: "/" + strings.Join(result.Remainder, "/"),
		}
	}

	return PartialRouteResult{
		Retrieved: true,
		Value:     result.Value,
		Remainder: "/" + strings.Join(result.Remainder, "/"),
	}
}
