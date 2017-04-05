package storage

// Adapter is an interface for pluggable components that will store mystack apps log messages
type Adapter interface {
	Start()
	Write(string, string) error
	Read(string, int) ([]string, error)
	Destroy(string) error
	Stop()
}
