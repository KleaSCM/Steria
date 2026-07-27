package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMetrics tracks performance metrics
type PerformanceMetrics struct {
	// Atomic counters
	FilesProcessed  int64
	BytesProcessed  int64
	CommitsCreated  int64
	BranchesCreated int64
	CacheHits       int64
	CacheMisses     int64

	// Timing metrics
	mu             sync.RWMutex
	operationTimes map[string][]time.Duration
	lastOperation  time.Time

	// Memory usage
	MemoryUsage int64
	CacheSize   int64
}

// NewPerformanceMetrics creates new performance metrics
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		operationTimes: make(map[string][]time.Duration),
	}
}

// StartOperation starts timing an operation
func (pm *PerformanceMetrics) StartOperation(operation string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		pm.RecordOperation(operation, duration)
	}
}

// RecordOperation records operation timing
func (pm *PerformanceMetrics) RecordOperation(operation string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.operationTimes[operation] = append(pm.operationTimes[operation], duration)
	pm.lastOperation = time.Now()

	// Keep only last 100 measurements per operation
	if len(pm.operationTimes[operation]) > 100 {
		pm.operationTimes[operation] = pm.operationTimes[operation][1:]
	}
}

// IncrementFilesProcessed increments file counter
func (pm *PerformanceMetrics) IncrementFilesProcessed(count int64) {
	atomic.AddInt64(&pm.FilesProcessed, count)
}

// IncrementBytesProcessed increments byte counter
func (pm *PerformanceMetrics) IncrementBytesProcessed(bytes int64) {
	atomic.AddInt64(&pm.BytesProcessed, bytes)
}

// IncrementCommitsCreated increments commit counter
func (pm *PerformanceMetrics) IncrementCommitsCreated() {
	atomic.AddInt64(&pm.CommitsCreated, 1)
}

// IncrementBranchesCreated increments branch counter
func (pm *PerformanceMetrics) IncrementBranchesCreated() {
	atomic.AddInt64(&pm.BranchesCreated, 1)
}

// IncrementCacheHit increments cache hit counter
func (pm *PerformanceMetrics) IncrementCacheHit() {
	atomic.AddInt64(&pm.CacheHits, 1)
}

// IncrementCacheMiss increments cache miss counter
func (pm *PerformanceMetrics) IncrementCacheMiss() {
	atomic.AddInt64(&pm.CacheMisses, 1)
}

// GetAverageTime gets average time for an operation
func (pm *PerformanceMetrics) GetAverageTime(operation string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.operationTimes[operation]
	if !exists || len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	return total / time.Duration(len(times))
}

// GetCacheHitRate gets cache hit rate
func (pm *PerformanceMetrics) GetCacheHitRate() float64 {
	hits := atomic.LoadInt64(&pm.CacheHits)
	misses := atomic.LoadInt64(&pm.CacheMisses)
	total := hits + misses

	if total == 0 {
		return 0
	}

	return float64(hits) / float64(total) * 100
}

// GetStats returns formatted statistics
func (pm *PerformanceMetrics) GetStats() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	stats := fmt.Sprintf("=== Steria Performance Stats ===\n")
	stats += fmt.Sprintf("Files Processed: %d\n", atomic.LoadInt64(&pm.FilesProcessed))
	stats += fmt.Sprintf("Bytes Processed: %d MB\n", atomic.LoadInt64(&pm.BytesProcessed)/(1024*1024))
	stats += fmt.Sprintf("Commits Created: %d\n", atomic.LoadInt64(&pm.CommitsCreated))
	stats += fmt.Sprintf("Branches Created: %d\n", atomic.LoadInt64(&pm.BranchesCreated))
	stats += fmt.Sprintf("Cache Hit Rate: %.2f%%\n", pm.GetCacheHitRate())
	stats += fmt.Sprintf("Last Operation: %s\n", pm.lastOperation.Format(time.RFC3339))

	// Operation timing stats
	stats += "\nOperation Timings:\n"
	for operation, times := range pm.operationTimes {
		if len(times) > 0 {
			var total time.Duration
			var min, max time.Duration = times[0], times[0]

			for _, t := range times {
				total += t
				if t < min {
					min = t
				}
				if t > max {
					max = t
				}
			}

			avg := total / time.Duration(len(times))
			stats += fmt.Sprintf("  %s: avg=%v, min=%v, max=%v, count=%d\n",
				operation, avg, min, max, len(times))
		}
	}

	return stats
}

// Reset resets all metrics
func (pm *PerformanceMetrics) Reset() {
	atomic.StoreInt64(&pm.FilesProcessed, 0)
	atomic.StoreInt64(&pm.BytesProcessed, 0)
	atomic.StoreInt64(&pm.CommitsCreated, 0)
	atomic.StoreInt64(&pm.BranchesCreated, 0)
	atomic.StoreInt64(&pm.CacheHits, 0)
	atomic.StoreInt64(&pm.CacheMisses, 0)

	pm.mu.Lock()
	pm.operationTimes = make(map[string][]time.Duration)
	pm.mu.Unlock()
}

// Global metrics instance
var GlobalMetrics = NewPerformanceMetrics()

// PerformanceProfiler provides profiling capabilities
type PerformanceProfiler struct {
	metrics *PerformanceMetrics
	start   time.Time
}

// StartProfiling starts performance profiling
func StartProfiling() *PerformanceProfiler {
	return &PerformanceProfiler{
		metrics: GlobalMetrics,
		start:   time.Now(),
	}
}

// EndProfiling ends profiling and returns results
func (pp *PerformanceProfiler) EndProfiling() string {
	duration := time.Since(pp.start)
	return fmt.Sprintf("Profiling completed in %v\n%s", duration, pp.metrics.GetStats())
}
