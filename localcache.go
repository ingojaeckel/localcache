package localcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var localcache map[string]CachedValue

const debugMode = false

func initializeIfNecessary() bool {
	if localcache == nil {
		localcache = make(map[string]CachedValue)
		return true
	}
	return false
}

func put(key, value string, ttlSeconds int64) {
	if debugMode {
		fmt.Printf("put(%s,%s)\n", key, value)
	}
	initializeIfNecessary()
	var expiry int64
	if ttlSeconds == 0 {
		expiry = 0 // Disables expiry
	} else {
		expiry = time.Now().Unix() + ttlSeconds
	}
	localcache[key] = CachedValue{value, expiry}
}

func get(key string) (string, bool) {
	if debugMode {
		fmt.Printf("get(%s)\n", key)
	}
	if initializeIfNecessary() {
		return "", false
	}
	exists := exists(key) // This will expire the key if necessary.
	if exists {
		return localcache[key].Value, true
	}
	return "", false
}

func exists(key string) bool {
	if debugMode {
		fmt.Printf("exists(%s)\n", key)
	}
	if initializeIfNecessary() {
		return false
	}
	_, exists := localcache[key]
	if exists && localcache[key].Expired() {
		expire(key)
		return false
	}
	return exists
}

func expire(key string) {
	delete(localcache, key)
	if debugMode {
		fmt.Printf("expired item %s", key)
	}
}

func getKey(someObject interface{}, id string) string {
	return fmt.Sprintf("%s.%s", reflect.TypeOf(someObject), id)
}

func Exists(id string, deserialized interface{}) bool {
	return exists(getKey(deserialized, id))
}

func Get(id string, deserialized interface{}) error {
	key := getKey(deserialized, id)

	item, exists := get(key)
	if !exists {
		return errors.New("Key not found")
	}
	dec := json.NewDecoder(strings.NewReader(string(item)))

	if err := dec.Decode(&deserialized); err != nil {
		return err
	}
	return nil
}

func Put(id string, someObject interface{}) error {
	return PutWithTTL(id, someObject, 0)
}

func PutWithTTL(id string, someObject interface{}, ttlSeconds int64) error {
	key := getKey(someObject, id)
	serializedObject, err := json.Marshal(someObject)
	if err != nil {
		return err
	}
	put(key, string(serializedObject), ttlSeconds)
	return nil
}

func (c CachedValue) Expired() bool {
	if c.ExpiresOn == 0 {
		return false // Expiry is disabled on this CachedValue
	}
	return time.Now().Unix() > c.ExpiresOn
}

type CachedValue struct {
	Value     string
	ExpiresOn int64 // Unix timestamp (seconds) for when this will expire. Or 0 to disable expiry.
}
