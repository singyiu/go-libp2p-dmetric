package common

type Marshalable interface {
	ToJsonBytes() ([]byte, error)
}

type Publishable interface {
	Marshalable
	ShouldBePublished() bool
	OnPublished()
}
