package storage

var storage = map[string]string{}

func SetStorage(key, value string) {
	if storage == nil {
		storage = map[string]string{}
	}
	storage[key] = value
}

func GetStorage(key string) (string, bool) {
	value, ok := storage[key]
	return value, ok
}

func DeleteStorage(key string) {
	delete(storage, key)
}
