package memo

func KeyedErr[K comparable, V any](f func(key K) (V, error)) func(key K) (V, error) {
	var m = make(map[K]V)

	return func(key K) (V, error) {
		v, ok := m[key]

		if ok {
			return v, nil
		}

		v, err := f(key)

		if err != nil {
			return v, err
		}

		m[key] = v
		return v, nil
	}
}
