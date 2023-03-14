package crawler

import "time"

type Config struct {
	MaxRedirects    int
	RevisitDelay    time.Duration
	WorkerCount     int
	ExcludePatterns []string
}
