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

func LastOne(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[len(arr)-1]
}