package storage

type Storage interface {
	Add(key, value string) bool
	Get(key string) (string, bool)
}

type Message struct {
	Key   string `json:"key"`
	Value string `json:"val"`
}
