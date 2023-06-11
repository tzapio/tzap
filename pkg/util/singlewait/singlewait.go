package singlewait

type SingleWait[T any] struct {
	data T
	done chan bool
}

// New creates a new instance of SingleWait struct and initiates the Done channel.
func New[T any](fn func() T) *SingleWait[T] {
	s := &SingleWait[T]{
		done: make(chan bool),
	}
	s.Process(fn)
	return s
}

// Process simulates a long running process and closes the Done channel when finished.
func (s *SingleWait[T]) Process(fn func() T) {
	go func() {
		// Simulate a long running process.
		res := fn()
		s.data = res
		// Close the Done channel when finished.
		close(s.done)
	}()
}
func (s *SingleWait[T]) GetData() T {
	s.waitForProcess()
	return s.data
}

// waitForProcess waits for the Process method to finish.
func (s *SingleWait[T]) waitForProcess() {
	<-s.done
}
