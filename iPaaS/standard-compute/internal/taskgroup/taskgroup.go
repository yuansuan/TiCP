package taskgroup

import (
	"context"
	"fmt"
	"sync"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

type TaskGroup struct {
	mu    sync.Mutex
	tasks map[string]task
	wg    sync.WaitGroup
	ctx   context.Context
}

func New(ctx context.Context) *TaskGroup {
	return &TaskGroup{
		mu:    sync.Mutex{},
		tasks: make(map[string]task),
		wg:    sync.WaitGroup{},
		ctx:   ctx,
	}
}

func (tg *TaskGroup) Add(t task) error {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	if _, exist := tg.tasks[t.Name()]; exist {
		return fmt.Errorf("task %s already exist", t.Name())
	}

	tg.tasks[t.Name()] = t
	return nil
}

func (tg *TaskGroup) StartAll() {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	for _, t := range tg.tasks {
		tg.wg.Add(1)
		go func(t task) {
			log.Infof("starting task %s ...", t.Name())
			if err := t.Start(tg.ctx); err != nil {
				log.Errorf("start task %s failed, %v", t.Name(), err)
			}

			log.Infof("task %s stopped", t.Name())
			tg.wg.Done()
		}(t)
	}
}

func (tg *TaskGroup) Wait() {
	tg.wg.Wait()
}
