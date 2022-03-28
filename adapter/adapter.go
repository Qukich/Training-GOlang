package adapter

type Adapter interface {
	WriteRateToFile() error
	CloseFile() error
	GetCode() string
	GetRateFromFile(ticker string) (*Departure, error)
	GetRateByTimestampFromFile(ticker string, timestamp int64) (*Departure, error)
}
