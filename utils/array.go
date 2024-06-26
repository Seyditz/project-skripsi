package utils

import "github.com/lib/pq"

func RemoveInt64FromArray(array pq.Int64Array, element int64) pq.Int64Array {
	index := -1
	for i, id := range array {
		if id == element {
			index = i
			break
		}
	}

	if index != -1 {
		return append(array[:index], array[index+1:]...)
	}

	return array
}
