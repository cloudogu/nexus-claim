package domain

import "reflect"

type sliceOfProperties []Properties

func IsEqual(left interface{}, right interface{}) bool {
	if isSliceOfInterfaces(left, right) {
		leftSliceOfProperties, rightSliceOfProperties, ok := castToSliceOfProperties(left, right)
		if ok {
			return isSliceEqual(leftSliceOfProperties, rightSliceOfProperties)
		}
	}
	return reflect.DeepEqual(left, right)
}

func castToSliceOfProperties(left interface{}, right interface{}) (sliceOfProperties, sliceOfProperties, bool) {
	leftSliceOfProperties, leftCastOk := castValueToSliceOfProperties(left.([]interface{}))
	if leftCastOk {
		rightSliceOfProperties, rightCastOk := castValueToSliceOfProperties(right.([]interface{}))
		if rightCastOk {
			return leftSliceOfProperties, rightSliceOfProperties, true
		}
	}
	return nil, nil, false
}

func isSliceOfInterfaces(left interface{}, right interface{}) bool {
	return isValueSliceOfInterfaces(left) && isValueSliceOfInterfaces(right)
}

func isValueSliceOfInterfaces(value interface{}) bool {
	typeOf := reflect.TypeOf(value)
	return typeOf.Kind() == reflect.Slice && typeOf.Elem().Kind() == reflect.Interface
}

func castValueToSliceOfProperties(value []interface{}) (sliceOfProperties, bool) {
	maps := make(sliceOfProperties, 0)
	for _, item := range value {
		mapItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, false
		}
		maps = append(maps, Properties(mapItem))
	}
	return maps, true
}

func isSliceEqual(left sliceOfProperties, right sliceOfProperties) bool {
	if len(left) > len(right) {
		return false
	}

	for _, leftItem := range left {
		if !contains(leftItem, right) {
			return false
		}
	}
	return true
}

func contains(item Properties, slice sliceOfProperties) bool {
	for _, sliceItem := range slice {
		if containsAllOf(item, sliceItem) {
			return true
		}
	}
	return false
}

func containsAllOf(left Properties, right Properties) bool {
	for key, value := range left {
		if IsEqual(value, right[key]) {
			return true
		}
	}
	return false
}
