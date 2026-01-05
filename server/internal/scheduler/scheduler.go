// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduler

import (
	"context"
	"log"
	"time"
)

// Task represents a scheduled task
type Task interface {
	Run(ctx context.Context) error
	Name() string
}

// Scheduler manages scheduled tasks
type Scheduler struct {
	tasks      []Task
	stop       chan struct{}
	lastRunDay int // Track last run day to prevent multiple runs per day
}

// New creates a new scheduler
func New() *Scheduler {
	return &Scheduler{
		tasks:      []Task{},
		stop:       make(chan struct{}),
		lastRunDay: -1, // Initialize to -1 so first run can happen
	}
}

// AddTask adds a task to the scheduler
func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	// Run immediately on start
	s.runTasks(ctx)

	for {
		select {
		case <-ticker.C:
			s.runTasks(ctx)
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) runTasks(ctx context.Context) {
	now := time.Now()
	currentDay := now.YearDay() // Day of year (1-365/366)
	
	// Run at 3 AM (03:00-03:09) and only once per day
	if now.Hour() == 3 && now.Minute() < 10 && s.lastRunDay != currentDay {
		s.lastRunDay = currentDay
		
		log.Printf("Running scheduled cleanup tasks at %s", now.Format("2006-01-02 15:04:05"))
		for _, task := range s.tasks {
			log.Printf("Running scheduled task: %s", task.Name())
			if err := task.Run(ctx); err != nil {
				// Always log task failures
				log.Printf("Task %s failed: %v", task.Name(), err)
			} else {
				log.Printf("Task %s completed successfully", task.Name())
			}
		}
	}
}

