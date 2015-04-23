package tmplnator

import (
	"fmt"
	"sync"

	l "github.com/Sirupsen/logrus"
)

type Gatherer int

func (g Gatherer) String() string {
	return fmt.Sprintf("Gatherer (%d)", g)
}

func (g Gatherer) gatherOutput(wg *sync.WaitGroup, done <-chan struct{}, out chan *File, c <-chan *File) {
	defer wg.Done()
	defer l.Debugf("%v: Exiting", g)
	l.Debugf("%v: Starting", g)
	for f := range c {
		select {
		case out <- f:
			l.Debugf("%v: Got Output: %v", g, f)
		case <-done:
			return
		}
	}
}
