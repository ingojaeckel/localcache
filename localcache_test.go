package localcache 

import (
	"testing"
)

type Point struct {
	X, Y int
}

func Test_Exists(t *testing.T) {
	var p Point
	if Exists("nonExistingValue", &p) {
		t.Errorf("value shouldn't have existed")
	}
}

func Test_Get(t *testing.T) {
	var p Point
	if err := Get("k", &p); err == nil {
		t.Errorf("expected an error")
	}
	Put("k", &Point{2,3})
	if err := Get("k", &p); err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if p.X != 2 || p.Y != 3 {
		t.Errorf("parameter mismatch")
	}
}

func Test_Put(t *testing.T) {
	Put("key1", &Point{1,2})

	var p Point
	
	if !Exists("key1", &p) {
		t.Errorf("value does not exist")
	}
	if err := Get("key1", &p); err != nil {
		t.Errorf("error getting value: %s", err)
	}
	if p.X != 1 {
		t.Errorf("X value mismatch: %d", p.X)
	}
	if p.Y != 2 {
		t.Errorf("Y value mismatch: %d", p.Y)
	}

}
