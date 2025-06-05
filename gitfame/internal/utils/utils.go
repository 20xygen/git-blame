package utils

func Contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func AnyFn(list []string, fn func(string) bool) bool {
	for _, item := range list {
		if fn(item) {
			return true
		}
	}
	return false
}

func AnyError(errors ...error) bool {
	for _, item := range errors {
		if item != nil {
			return true
		}
	}
	return false
}

func AllFn(list []string, fn func(string) bool) bool {
	for _, item := range list {
		if !fn(item) {
			return false
		}
	}
	return true
}
