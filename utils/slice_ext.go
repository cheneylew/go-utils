package utils

func JKRemoveInt(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func JKRemoveInterface(slice []interface{}, s int) []interface {} {
	return append(slice[:s], slice[s+1:]...)
}

func InSlice(val string, slice []string) bool {
	for _, value := range slice {
		if val == value {
			return true
		}
	}

	return false
}