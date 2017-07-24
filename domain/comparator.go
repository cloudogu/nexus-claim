package domain

import "reflect"

func IsEqual(left interface{}, right interface{}) bool {
	if isSliceOfInterfaces(left) && isSliceOfInterfaces(right) {
		leftSliceOfMap, leftCastOk := castToSliceOfMaps(left.([]interface{}))
		if leftCastOk {
			rightSliceOfMap, rightCastOk := castToSliceOfMaps(right.([]interface{}))
			if rightCastOk {
				return isSliceEqual(leftSliceOfMap, rightSliceOfMap)
			}
		}
	}
	return reflect.DeepEqual(left, right)
}

func isSliceOfInterfaces(value interface{}) bool {
	typeOf := reflect.TypeOf(value)
	return typeOf.Kind() == reflect.Slice && typeOf.Elem().Kind() == reflect.Interface
}

func castToSliceOfMaps(value []interface{}) ([]map[string]interface{}, bool) {
	sliceOfMaps := []map[string]interface{}{}
	for _, item := range value {
		mapItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, false
		}
		sliceOfMaps = append(sliceOfMaps, mapItem)
	}
	return sliceOfMaps, true
}

func isSliceEqual(left []map[string]interface{}, right []map[string]interface{}) bool {
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

func contains(item map[string]interface{}, slice []map[string]interface{}) bool {
	for _, sliceItem := range slice {
		if containsAllOf(item, sliceItem) {
			return true
		}
	}
	return false
}

func containsAllOf(left map[string]interface{}, right map[string]interface{}) bool {
	for key, value := range left {
		if IsEqual(value, right[key]) {
			return true
		}
	}
	return false
}
