package main

import (
	"errors"
	"strings"
)

// All of this

type routeNodeOrRouteNodeMap struct {
	child    *routeNode
	children *map[string]*routeNode
}

func (r routeNodeOrRouteNodeMap) IsChildOnly() bool {
	return r.child != nil
}

func newChild(child routeNode) *routeNodeOrRouteNodeMap {
	newChild := child
	return &routeNodeOrRouteNodeMap{&newChild, nil}
}

func newChildren(children map[string]*routeNode) *routeNodeOrRouteNodeMap {
	newChildren := children
	return &routeNodeOrRouteNodeMap{nil, &newChildren}
}

// A children of nodes, or just a single node.
type routeNode struct {
	childOrChildren *routeNodeOrRouteNodeMap
	value           interface{}
}

func newRouteNode(components []string, value interface{}) routeNode {
	var r routeNode
	r.value = newChild
	r.childOrChildren = newChildren(make(map[string]*routeNode))

	r.add(components, value)
	return r
}

// add adds a new sub rooute.
func (r *routeNode) add(components []string, value interface{}) {
	if len(components) <= 0 {
		r.value = value
	} else {
		first, remainder := components[0], components[1:]
		if first[0] == ':' {
			if r.childOrChildren.IsChildOnly() {
				r.childOrChildren.child.add(remainder, value)
			} else {
				r.childOrChildren = newChild(newRouteNode(remainder, value))
			}
		} else {
			if r.childOrChildren.IsChildOnly() {
				r.childOrChildren = newChildren(make(map[string]*routeNode))
				node := newRouteNode(remainder, value)
				(*r.childOrChildren.children)[first] = &node
			} else {
				node, ok := (*r.childOrChildren.children)[first]
				if !ok {
					node := newRouteNode(remainder, value)
					(*r.childOrChildren.children)[first] = &node
				} else {
					node.add(remainder, value)
				}
			}
		}
	}
}

func (r *routeNode) get(components []string) interface{} {
	if len(components) <= 0 {
		return r.value
	}
	first, remainder := components[0], components[1:]
	if r.childOrChildren.IsChildOnly() {
		return r.childOrChildren.child.get(remainder)
	}
	node, ok := (*r.childOrChildren.children)[first]
	if !ok {
		return nil
	}
	return node.get(remainder)
}

type PartialRouteNodeResult struct {
	Retrieved bool
	Value     interface{}
	Remainder []string
}

func (r *routeNode) getPartial(components []string) PartialRouteNodeResult {
	if len(components) <= 0 {
		return PartialRouteNodeResult{
			Retrieved: true,
			Value:     r.value,
			Remainder: components,
		}
	}
	first, remainder := components[0], components[1:]
	if r.childOrChildren.IsChildOnly() {
		return r.childOrChildren.child.getPartial(remainder)
	}
	node, ok := (*r.childOrChildren.children)[first]
	if !ok {
		return PartialRouteNodeResult{
			Retrieved: true,
			Value:     r.value,
			Remainder: remainder,
		}
	}
	return node.getPartial(remainder)
}

// routes get the routes.
type routes struct {
	children map[string]routeNode
}

// add adds a new route.
func (r *routes) add(route string, value interface{}) error {
	if len(route) <= 0 {
		return errors.New("Route cannot be empty")
	}

	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]
	node, ok := r.children[first]
	if !ok {
		r.children[first] = newRouteNode(remainder, value)
	} else {
		node.add(remainder, value)
	}

	return nil
}

func (r routes) get(route string) interface{} {
	if len(route) <= 0 {
		return nil
	}
	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]

	node, ok := r.children[first]
	if !ok {
		return nil
	}
	return node.get(remainder)
}

func newRouter() routes {
	return routes{make(map[string]routeNode)}
}

type partialRouteResult struct {
	retrieved bool
	value     interface{}
	remainder string
}

// Let's say we only have a handler registered at /foo/bar, but we request a
// handler at /foo/bar/baz, then we will still get the handler at /foo/bar.
func (r routes) getShortCircuited(route string) partialRouteResult {
	if len(route) <= 0 {
		return partialRouteResult{}
	}
	components := strings.Split(route, "/")
	first, remainder := components[0], components[1:]
	node, ok := r.children[first]
	if !ok {
		return partialRouteResult{
			retrieved: false,
			value:     nil,
			remainder: route,
		}
	}

	result := node.getPartial(remainder)
	if !result.Retrieved {
		return partialRouteResult{
			retrieved: false,
			value:     nil,
			remainder: "/" + strings.Join(result.Remainder, "/"),
		}
	}

	return partialRouteResult{
		retrieved: true,
		value:     result.Value,
		remainder: "/" + strings.Join(result.Remainder, "/"),
	}
}
