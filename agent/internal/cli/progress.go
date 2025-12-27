// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// ProgressBar represents a progress bar
type ProgressBar struct {
	total    int64
	current  int64
	width    int
	startTime time.Time
	lastUpdate time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:     int64(total),
		current:   0,
		width:     50, // Progress bar width in characters
		startTime: time.Now(),
		lastUpdate: time.Now(),
	}
}

// Update increments the progress counter and updates the display
func (p *ProgressBar) Update(delta int64) {
	atomic.AddInt64(&p.current, delta)
	p.Refresh()
}

// Refresh updates the progress bar display
func (p *ProgressBar) Refresh() {
	now := time.Now()
	// Throttle updates to avoid too many screen refreshes (max 10 times per second)
	if now.Sub(p.lastUpdate) < 100*time.Millisecond && p.current < p.total {
		return
	}
	p.lastUpdate = now

	current := atomic.LoadInt64(&p.current)
	total := p.total

	if total == 0 {
		return
	}

	percentage := float64(current) * 100.0 / float64(total)
	width := p.width
	filled := int(float64(width) * percentage / 100.0)
	
	// Calculate elapsed time and estimate remaining time
	elapsed := now.Sub(p.startTime)
	var remaining time.Duration
	var speed float64
	if current > 0 {
		speed = float64(current) / elapsed.Seconds()
		if speed > 0 {
			remaining = time.Duration(float64(total-current)/speed) * time.Second
		}
	}

	// Build progress bar
	bar := make([]byte, width)
	for i := 0; i < width; i++ {
		if i < filled {
			bar[i] = '='
		} else {
			bar[i] = ' '
		}
	}

	// Format time
	elapsedStr := formatDuration(elapsed)
	remainingStr := formatDuration(remaining)
	if remaining <= 0 {
		remainingStr = "--:--"
	}

	// Print progress bar (use \r to return to beginning of line, \033[K to clear to end)
	// Output to stderr so it doesn't interfere with stdout
	fmt.Fprintf(os.Stderr, "\r[%s] %.1f%% (%d/%d) | Elapsed: %s | Remaining: %s | Speed: %.1f files/s\033[K",
		string(bar), percentage, current, total, elapsedStr, remainingStr, speed)
}

// Finish completes the progress bar and prints final summary
func (p *ProgressBar) Finish() {
	current := atomic.LoadInt64(&p.current)
	total := p.total
	elapsed := time.Since(p.startTime)

	// Clear progress bar line and print final summary
	fmt.Fprintf(os.Stderr, "\r\033[K") // Clear line
	fmt.Fprintf(os.Stderr, "Completed: %d/%d files in %s\n", current, total, formatDuration(elapsed))
}

// formatDuration formats a duration as MM:SS or HH:MM:SS
func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

