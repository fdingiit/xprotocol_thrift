package map_helper

func GetValues(dataMap map[string]string) []string {
	values := make([]string, len(dataMap))

	i := 0
	for _, v := range dataMap {
		values[i] = v
		i++
	}

	return values
}

func Filp(dataMap map[string]string) map[string]string {
	var result map[string]string
	for k, v := range dataMap {
		result[v] = k
	}
	return result
}

// use map to simulate Set, using byte as map value data structure to save memory
func GetValuesMap(dataMap *map[string]string) map[string]byte {
	var result map[string]byte = make(map[string]byte)
	if dataMap == nil || len(*dataMap) == 0 {
		return result
	}
	for _, v := range *dataMap {
		result[v] = 0
	}
	return result
}
