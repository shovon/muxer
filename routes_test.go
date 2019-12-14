package main

import (
	"testing"
)

func TestInsertAndGet(t *testing.T) {
	router := newRouter()
	router.add("/haha", 10)
	router.add("/hello", 20)
	router.add("/hello/foo", 40)
	router.add("/foo/:bar", 32)

	haha, ok := router.get("/haha").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if haha != 10 {
		t.Error("Expected haha to be 10")
	}

	hello, ok := router.get("/hello").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if hello != 20 {
		t.Error("Expected hello to be 20")
	}

	foo, ok := router.get("/hello/foo").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if foo != 40 {
		t.Error("Expected haha to be 40")
	}

	bar, ok := router.get("/foo/20").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if bar != 32 {
		t.Error("Expdected bar to be 32")
	}
}

func TestPartialGet(t *testing.T) {
	router := newRouter()
	router.add("/haha", 10)
	router.add("/foobar/:something", 20)

	var result partialRouteResult

	result = router.getShortCircuited("/haha")
	haha, ok := result.value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if haha != 10 {
		t.Error("Expected haha to be 10")
	}

	result = router.getShortCircuited("/haha/nothing")
	nothing, ok := result.value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if nothing != 10 {
		t.Error("Expected haha to be 10")
	}

	result = router.getShortCircuited("/foobar/haha/lol")
	lol, ok := result.value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if lol != 20 {
		t.Error("Expected foobar to be 20")
	}
}
