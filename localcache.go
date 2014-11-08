package localcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var localcache map[string]string
const debugMode = false

func initializeIfNecessary() bool {
	if localcache == nil {
		localcache = make(map[string]string)
		return true
	}
	return false
}

func put(key, value string) {
	if debugMode {
		fmt.Printf("put(%s,%s)\n", key, value)
	}
	initializeIfNecessary()
	localcache[key] = value
}

func get(key string) (string, bool) {
	if debugMode {
		fmt.Printf("get(%s)\n", key)
	}
	if initializeIfNecessary() {
		return "", false
	}
	return localcache[key], exists(key)
}

func exists(key string) bool {
	if debugMode {
		fmt.Printf("exists(%s)\n", key)
	}
	if initializeIfNecessary() {
		return false
	}
	_, exists := localcache[key]
	return exists
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
	key := getKey(someObject, id)

	serializedObject, err := json.Marshal(someObject)
	if err != nil {
		return err
	}

	put(key, string(serializedObject))
	return nil
}
