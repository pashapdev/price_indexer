package indexer

import "time"

type timeWindow struct {
	min, max time.Time
}

func (w *timeWindow) Next(current time.Time) {
	w.min = current
	w.max = w.min.Add(time.Minute)
}

func (w *timeWindow) Include(current time.Time) bool {
	if current.Before(w.min) {
		return false
	}

	if current.After(w.max) || current.Equal(w.max) {
		return false
	}
	return true
}
