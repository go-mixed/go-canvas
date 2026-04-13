package misc

func Filter[E any](list []E, fn func(item E, index int) bool) []E {
	var result []E
	for i, item := range list {
		if fn(item, i) {
			result = append(result, item)
		}
	}
	return result
}

func Map[E any, R any](list []E, fn func(item E, index int) R) []R {
	var result []R
	for i, item := range list {
		result = append(result, fn(item, i))
	}
	return result
}

func MapFilter[E any, R any](list []E, fn func(item E, index int) (R, bool)) []R {
	var result []R
	for i, item := range list {
		r, ok := fn(item, i)
		if ok {
			result = append(result, r)
		}
	}
	return result
}

func Empty[E any]() E {
	var e E
	return e
}

func First[E any](list []E) E {
	if len(list) > 0 {
		return list[0]
	}
	return Empty[E]()
}

func Last[E any](list []E) E {
	if len(list) > 0 {
		return list[len(list)-1]
	}
	return Empty[E]()
}

func ToPtr[E any](v E) *E {
	return &v
}
