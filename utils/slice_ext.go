package utils

func JKRemoveInt(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func JKRemoveInterface(slice []interface{}, s int) []interface {} {
	return append(slice[:s], slice[s+1:]...)
}
