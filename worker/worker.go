// A convenient way to run jobs from beanstalkd.
//
// This package is still in flux.
//
//   import "beanstalk/worker"
//   worker.Handle("resize", my_resizer)
//   worker.Handle("sync", my_syncer)
//   worker.DialAndWork("localhost:11300", nil)
//
package worker

import (
	"beanstalk"
	"os"
)

type Reserver interface {
	Reserve() (beanstalk.Job, os.Error)
}

type Runner interface {
	Run(beanstalk.Job)
}

type Worker interface {
	Work(Reserver)
	TubeNames() []string
}

type TubeMux struct {
	DefaultRunner Runner
	handlers map[string]Runner
}

func (w *TubeMux) Handle(tubeName string, r Runner) {
	w.handlers[tubeName] = r
}

func (w *TubeMux) Work(r Reserver) os.Error {
	for {
		j := r.Reserve()
		h, ok := w.handlers[j.Tube]
		if !ok {
			h = w.DefaultRunner
		}
		h.Run(j)
	}
}

// Install a runner in DefaultTubeMux.
func Handle(tubeName string, r Runner) {
	DefaultTubeMux.Handle(tubeName, r)
}

// Open a connection to beanstalkd and run jobs. If w is nil, use
// DefaultTubeMux.
func DialAndWork(addr string, w Worker) os.Error {
	if w == nil {
		w = DefaultTubeMux
	}
	c, err := beanstalk.Dial(addr)
	if err != nil {
		return err
	}
	ts := c.NewTubeSet(w.tubeNames())
	err = w.Work(ts)
	//c.Close()
	return err
}
