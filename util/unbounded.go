package util

func MakeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		curVal := func() interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}
		outCh := func() chan interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return out
		}
		for len(inQueue) > 0 || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue = append(inQueue, v)
				}
			case outCh() <- curVal():
				inQueue = inQueue[1:]
			}
		}
		close(out)
	}()
	return in, out
}

func MakeInfinitePriority[T any]() (chan<- Item[T], <-chan Item[T]) {
	in := make(chan Item[T])
	out := make(chan Item[T])
	go func() {
		inQueue := Pqueue_init[T]()
		curVal := func() Item[T] {
			if inQueue.Len() == 0 {
				return Item[T]{}
			}
			retval := inQueue.Pop()
			return *retval
		}
		outCh := func() chan Item[T] {
			if inQueue.Len() == 0 {
				return nil
			}
			return out
		}
		for inQueue.Len() > 0 || in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue.Push(&v)
				}
			case outCh() <- curVal():

			}
		}
		close(out)
	}()
	return in, out
}
