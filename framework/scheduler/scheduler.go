package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Job interface {
	GetName() string
	GetTimeSpec() string
	Run()
}

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler(
	jobs ...Job,
) *Scheduler {
	z := zap.S().With("context", "NewScheduler")

	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		z.Fatal("Scheduler parse timezone failed: %v", err)
	}

	c := cron.New(cron.WithLocation(l))

	for _, j := range jobs {
		id, err := c.AddJob(j.GetTimeSpec(), j)
		if err != nil {
			z.Fatal("Scheduler adding job failed: %v", err)
		}

		z.Infow("Registered job", "job", j.GetName(), "timeSpec", j.GetTimeSpec(), "job_id", id)
	}

	return &Scheduler{
		c: c,
	}
}

func (sc *Scheduler) Start() {
	sc.c.Start()
}

func (sc *Scheduler) Stop() context.Context {
	return sc.c.Stop()
}
