package common

func IsNotNull(value interface{}) bool {
	return value != nil
}

func Ternary(cond bool, val1 interface{}, val2 interface{}) interface{} {
	if cond {
		return val1
	}
	return val2
}
