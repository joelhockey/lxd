package ioprogress

import (
	"time"
)

// ProgressTracker provides the stream information needed for tracking
type ProgressTracker struct {
	Length  int64
	Handler func(size, processed, percent, speed int64)

	percentage float64
	processed  int64
	start      *time.Time
	last       *time.Time
}

func (pt *ProgressTracker) update(n int) {
	// Skip the rest if no handler attached
	if pt.Handler == nil {
		return
	}

	// Initialize start time if needed
	if pt.start == nil {
		cur := time.Now()
		pt.start = &cur
		pt.last = pt.start
	}

	// Skip if no data to count
	if n <= 0 {
		return
	}

	// Update interval handling
	var percentage float64
	if pt.Length > 0 {
		// If running in relative mode, check that we increased by at least 1%
		percentage = float64(pt.processed) / float64(pt.Length) * float64(100)
		if percentage-pt.percentage < 0.9 {
			return
		}
	} else {
		// If running in absolute mode, check that at least a second elapsed
		interval := time.Since(*pt.last).Seconds()
		if interval < 1 {
			return
		}
	}

	// Determine speed
	speedInt := int64(0)
	duration := time.Since(*pt.start).Seconds()
	if duration > 0 {
		speed := float64(pt.processed) / duration
		speedInt = int64(speed)
	}

	// Determine progress
	var percentInt int64
	if pt.Length > 0 {
		pt.percentage = percentage
		percentInt = int64(1 - (int(percentage) % 1) + int(percentage))
		if percentInt > 100 {
			percentInt = 100
		}
	} else {
		// Update timestamp
		cur := time.Now()
		pt.last = &cur
	}

	pt.Handler(pt.Length, pt.processed, percentInt, speedInt)
}
