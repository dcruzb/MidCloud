package common

type CloudFunctions interface {
	Price(size float64) float64
	Availability() bool
}
