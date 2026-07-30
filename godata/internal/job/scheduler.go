package job

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Job interface {
	Name() string
	Run(ctx context.Context) error
}

type ScheduledJob struct {
	Job      Job
	Interval time.Duration
}

type Scheduler struct {
	jobs   []ScheduledJob
	logger *zap.Logger
	mu     sync.Mutex
	done   chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		logger: zap.L().Named("scheduler"),
		done:   make(chan struct{}),
	}
}

func (s *Scheduler) AddJob(job Job, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, ScheduledJob{Job: job, Interval: interval})
}

func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, sj := range s.jobs {
		go s.runJob(ctx, sj)
	}

	s.logger.Info("scheduler started", zap.Int("jobs", len(s.jobs)))
}

func (s *Scheduler) Stop() {
	close(s.done)
	s.logger.Info("scheduler stopped")
}

func (s *Scheduler) runJob(ctx context.Context, sj ScheduledJob) {
	ticker := time.NewTicker(sj.Interval)
	defer ticker.Stop()

	s.logger.Info("job scheduled",
		zap.String("name", sj.Job.Name()),
		zap.Duration("interval", sj.Interval),
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		case <-ticker.C:
			start := time.Now()
			s.logger.Info("job started", zap.String("name", sj.Job.Name()))

			if err := sj.Job.Run(ctx); err != nil {
				s.logger.Error("job failed",
					zap.String("name", sj.Job.Name()),
					zap.Duration("duration", time.Since(start)),
					zap.Error(err),
				)
			} else {
				s.logger.Info("job completed",
					zap.String("name", sj.Job.Name()),
					zap.Duration("duration", time.Since(start)),
				)
			}
		}
	}
}
