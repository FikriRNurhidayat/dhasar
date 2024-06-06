package dhasar

type Adapter[Opt any, Result any] interface {
	Connect(Opt) (Result, error)
	Close() error
}
