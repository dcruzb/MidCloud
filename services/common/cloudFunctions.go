package common

type CloudFunctions interface {
	Price(size float64) (price float64, err error)
	Availability() (available bool, err error)
}
