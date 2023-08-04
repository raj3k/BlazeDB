package utils

import "reflect"

// Function to check if a variable is of type interface{}
func IsInterface(variable interface{}) bool {
	return reflect.TypeOf(variable).Kind() == reflect.Interface
}

// Function to check if a variable is of type []interface{}
func IsSliceOfInterface(variable interface{}) bool {
	t := reflect.TypeOf(variable)
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Interface
}
