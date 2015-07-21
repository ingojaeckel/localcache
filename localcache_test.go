package localcache 

import (
	"testing"
	"time"
)

type Point struct {
	X, Y int
}

func Test_InitializeOnGet(t *testing.T) {
	var p Point
	if Get("nonExistingValue", &p) == nil {
		t.Errorf("value shouldn't have existed")
	}
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

func Test_Put_With_TTL(t *testing.T) {
	var p Point
	PutWithTTL("key1", &Point{1,2}, -1)

	if Exists("key1", &p) {
		t.Errorf("value should not exist since it's expired")
	}
	
	PutWithTTL("key1", &Point{1,2}, -1)
	if Get("key1", &p) == nil {
		t.Errorf("should have received an error because this key should have been expired already")
	}
}

func Test_Expiry(t *testing.T) {
	notExpired := CachedValue{"foo", time.Now().Unix() + 60}
	if notExpired.Expired() {
		t.Errorf("This should only expire 60 sec from now.")
	}
	notExpired2 := CachedValue{"foo", 0}
	if notExpired2.Expired() {
		t.Errorf("This should only expire 60 sec from now.")
	}
	notExpired3 := CachedValue{Value:"foo"}
	if notExpired3.Expired() {
		t.Errorf("This should only expire 60 sec from now.")
	}
	expired := CachedValue{"foo", time.Now().Unix() - 1}
	if !expired.Expired() {
		t.Errorf("This should have expired 1 sec ago.")
	}
}