package worker

import "context"

type IWorker interface {
	Start(ctx context.Context)
	Stop()
}

type Runner struct {
	workers []IWorker
}

func (r *Runner) Launch(ctx context.Context) {
	for _, w := range r.workers {
		go func() {
			w.Start(ctx)
		}()
	}
}

func (r *Runner) Stop() {

}

func NewRunner(workers ...IWorker) *Runner {
	return &Runner{
		workers: workers,
	}
}
