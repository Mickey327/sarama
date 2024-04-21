package breaker

type Breaker interface {
	Run(work func() error) error
	Go(work func() error) error
}

func NewDummy() Breaker {
	return &dummyBreaker{}
}

type dummyBreaker struct{}

func (d dummyBreaker) Run(work func() error) error { return work() }
func (d dummyBreaker) Go(work func() error) error  { return work() }
