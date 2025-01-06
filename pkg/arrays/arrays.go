package arrays

func HasStrItem(arr []string, item any) bool {
	for _, element := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func HasIntItem(arr []int, item any) bool {
	for _, element := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func HasAnyItem(arr []any, item any) bool {
	for _, element := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func IndexOf(collection []any, el any) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}
