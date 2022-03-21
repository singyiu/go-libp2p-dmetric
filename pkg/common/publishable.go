package common

type Publishable interface {
	Marshalable
	ShouldBePublished() bool
	OnPublished()
}
