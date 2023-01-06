package worker

type Worker interface {
	Do(interface{}) error
}

type WorkerFunc func(interface{}) error

func (s WorkerFunc) Do(a interface{}) error {
	return s(a)
}

func FromFunc(f func(interface{}) error) Worker {
	return WorkerFunc(f)
}
