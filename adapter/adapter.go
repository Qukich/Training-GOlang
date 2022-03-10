package adapter

type Adapter interface {
	WriteRateToDatabase() error
	CloseDB() error
}
