package infoprovider

type InfoProvider interface {
	GetInfo() (any, error)
}
