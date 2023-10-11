package maps

func GetFirstOfMap(data map[string]any) (string, any) {
	for k, v := range data {
		if childMap, ok := v.(map[string]any); ok {
			return GetFirstOfMap(childMap)
		} else {
			return k, v
		}
	}
	return "", nil
}
