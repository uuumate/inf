package value

// ClearSlice
// clear all elem in slice
func ClearSlice(slice *[]interface{}) {
	*slice = (*slice)[:0]
}
