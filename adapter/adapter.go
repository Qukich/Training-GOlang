package adapter

type Adapter interface {
	WriteRateToDatabase() error
}