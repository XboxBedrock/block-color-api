package util

// This is a ternary utility to help do simpler math
// It acts like the :? syntax in JavaScript and is functionally the same, simply returning a if true or b if the condition is false
func Ternary(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
