package common

type Marshalable interface {
	ToJsonBytes() ([]byte, error)
}
