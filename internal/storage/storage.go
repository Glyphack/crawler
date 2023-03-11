package storage

type Storage interface {
	Get(path string) (string, error)
	Set(path string, value string) error
	Delete(path string) error
}
