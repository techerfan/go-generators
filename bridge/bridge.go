package bridge

import "github.com/techerfan/go-generators/ordone"

func Bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	// this is the channel that will return all values from bridge.
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		// this loop is responsible for pulling channels off of
		// chanStream and providing them to a nested loop for use.
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range ordone.OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
