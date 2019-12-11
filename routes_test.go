package main

import (
	"testing"
)

func TestInsertAndGet(t *testing.T) {
	router := NewRouter()
	router.Add("/haha", 10)
	router.Add("/hello", 20)
	router.Add("/hello/foo", 40)
	router.Add("/foo/:bar", 32)

	haha, ok := router.Get("/haha").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if haha != 10 {
		t.Error("Expected haha to be 10")
	}

	hello, ok := router.Get("/hello").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if hello != 20 {
		t.Error("Expected hello to be 20")
	}

	foo, ok := router.Get("/hello/foo").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if foo != 40 {
		t.Error("Expected haha to be 40")
	}

	bar, ok := router.Get("/foo/20").(int)
	if !ok {
		t.Error("Not an integer")
	}
	if bar != 32 {
		t.Error("Expdected bar to be 32")
	}
}

func TestPartialGet(t *testing.T) {
	router := NewRouter()
	router.Add("/haha", 10)
	router.Add("/foobar/:something", 20)

	var result PartialRouteResult

	result = router.GetPartial("/haha")
	haha, ok := result.Value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if haha != 10 {
		t.Error("Expected haha to be 10")
	}

	result = router.GetPartial("/haha/nothing")
	nothing, ok := result.Value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if nothing != 10 {
		t.Error("Expected haha to be 10")
	}

	result = router.GetPartial("/foobar/haha/lol")
	lol, ok := result.Value.(int)
	if !ok {
		t.Error("Not an integer")
	}
	if lol != 20 {
		t.Error("Expected foobar to be 20")
	}
}
