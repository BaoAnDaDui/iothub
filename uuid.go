package iothub

type IdProvider interface {
	ID() (string, error)
}
