// Package metrics - Caching infrastructure with TTL support for metrics
package metrics

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MetricsCache provides caching functionality for calculated metrics
type MetricsCache struct {
	store     map[string]*CacheEntry
	mu        sync.RWMutex
	defaultTTL time.Duration
	maxSize   int
	stats     CacheStats
}

// CacheEntry represents a cached metrics entry
type CacheEntry struct {
	Key        string      `json:"key"`
	Data       interface{} `json:"data"`
	CachedAt   time.Time   `json:"cached_at"`
	ExpiresAt  time.Time   `json:"expires_at"`
	TTL        time.Duration `json:"ttl"`
	AccessCount int        `json:"access_count"`
	LastAccessed time.Time `json:"last_accessed"`
	DataSize    int        `json:"data_size"`
	ComputeTime time.Duration `json:"compute_time"`
	Tags       []string   `json:"tags"`
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits           int64         `json:"hits"`
	Misses         int64         `json:"misses"`
	Evictions      int64         `json:"evictions"`
	TotalEntries   int           `json:"total_entries"`
	TotalSize      int           `json:"total_size"`
	HitRate        float64       `json:"hit_rate"`
	AverageLatency time.Duration `json:"average_latency"`
	LastReset      time.Time     `json:"last_reset"`
}

// CacheConfig configures the metrics cache
type CacheConfig struct {
	DefaultTTL       time.Duration            `json:"default_ttl"`
	MaxSize          int                      `json:"max_size"`
	CleanupInterval  time.Duration            `json:"cleanup_interval"`
	MetricTypeTTLs   map[string]time.Duration `json:"metric_type_ttls"`
	EnableStats      bool                     `json:"enable_stats"`
	PersistentCache  bool                     `json:"persistent_cache"`
}

// NewMetricsCache creates a new metrics cache
func NewMetricsCache(config CacheConfig) *MetricsCache {
	if config.DefaultTTL == 0 {
		config.DefaultTTL = 15 * time.Minute
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1000
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute
	}

	cache := &MetricsCache{
		store:      make(map[string]*CacheEntry),
		defaultTTL: config.DefaultTTL,
		maxSize:    config.MaxSize,
		stats: CacheStats{
			LastReset: time.Now(),
		},
	}

	// Start cleanup goroutine
	go cache.startCleanupWorker(config.CleanupInterval)

	return cache
}

// GenerateCacheKey generates a cache key for metrics request
func (mc *MetricsCache) GenerateCacheKey(metricType string, request MetricsRequest) string {
	keyData := struct {
		Type        string        `json:"type"`
		Repository  string        `json:"repository"`
		TimeRange   TimeRange     `json:"time_range"`
		Granularity string        `json:"granularity"`
	}{
		Type:        metricType,
		Repository:  fmt.Sprintf("%s/%s", request.Repository.Owner, request.Repository.Name),
		TimeRange:   request.TimeRange,
		Granularity: request.Granularity,
	}

	keyJSON, _ := json.Marshal(keyData)
	hash := md5.Sum(keyJSON)
	return fmt.Sprintf("%s_%x", metricType, hash)
}

// Get retrieves a cached metric
func (mc *MetricsCache) Get(ctx context.Context, key string) (*CacheEntry, bool) {
	start := time.Now()
	defer func() {
		mc.updateStats(time.Since(start))
	}()

	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.store[key]
	if !exists {
		mc.stats.Misses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		mc.stats.Misses++
		// Remove expired entry (will be done in cleanup)
		return nil, false
	}

	// Update access info
	entry.AccessCount++
	entry.LastAccessed = time.Now()
	mc.stats.Hits++

	return entry, true
}

// Set stores a metric in cache
func (mc *MetricsCache) Set(ctx context.Context, key string, data interface{}, ttl time.Duration, computeTime time.Duration, tags []string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if ttl == 0 {
		ttl = mc.defaultTTL
	}

	// Calculate data size
	dataSize := 0
	if dataJSON, err := json.Marshal(data); err == nil {
		dataSize = len(dataJSON)
	}

	now := time.Now()
	entry := &CacheEntry{
		Key:          key,
		Data:         data,
		CachedAt:     now,
		ExpiresAt:    now.Add(ttl),
		TTL:          ttl,
		AccessCount:  0,
		LastAccessed: now,
		DataSize:     dataSize,
		ComputeTime:  computeTime,
		Tags:         tags,
	}

	// Check if we need to evict entries
	if len(mc.store) >= mc.maxSize {
		mc.evictLRU()
	}

	mc.store[key] = entry
	mc.stats.TotalEntries = len(mc.store)
	mc.stats.TotalSize += dataSize

	return nil
}

// Delete removes an entry from cache
func (mc *MetricsCache) Delete(ctx context.Context, key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if entry, exists := mc.store[key]; exists {
		delete(mc.store, key)
		mc.stats.TotalEntries = len(mc.store)
		mc.stats.TotalSize -= entry.DataSize
	}
}

// Clear removes all entries from cache
func (mc *MetricsCache) Clear(ctx context.Context) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.store = make(map[string]*CacheEntry)
	mc.stats.TotalEntries = 0
	mc.stats.TotalSize = 0
}

// InvalidateByTags removes all entries with matching tags
func (mc *MetricsCache) InvalidateByTags(ctx context.Context, tags []string) int {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	invalidated := 0
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}

	for key, entry := range mc.store {
		for _, entryTag := range entry.Tags {
			if tagSet[entryTag] {
				delete(mc.store, key)
				mc.stats.TotalSize -= entry.DataSize
				invalidated++
				break
			}
		}
	}

	mc.stats.TotalEntries = len(mc.store)
	return invalidated
}

// GetStats returns cache performance statistics
func (mc *MetricsCache) GetStats() CacheStats {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	stats := mc.stats
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total)
	}
	stats.TotalEntries = len(mc.store)

	return stats
}

// ResetStats resets cache statistics
func (mc *MetricsCache) ResetStats() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.stats = CacheStats{
		TotalEntries: len(mc.store),
		TotalSize:    mc.stats.TotalSize,
		LastReset:    time.Now(),
	}
}

// GetTTLForMetricType returns the TTL for a specific metric type
func (mc *MetricsCache) GetTTLForMetricType(metricType string) time.Duration {
	// Default TTLs for different metric types
	ttlMap := map[string]time.Duration{
		"dora":       15 * time.Minute, // DORA metrics update frequently
		"chi":        1 * time.Hour,    // Code health changes slowly
		"ai":         10 * time.Minute, // AI metrics update regularly
		"aggregated": 30 * time.Minute, // Aggregated metrics are expensive to compute
		"trends":     2 * time.Hour,    // Trend analysis can be cached longer
	}

	if ttl, exists := ttlMap[metricType]; exists {
		return ttl
	}
	return mc.defaultTTL
}

// Private methods

// startCleanupWorker runs periodic cleanup of expired entries
func (mc *MetricsCache) startCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.cleanup()
		}
	}
}

// cleanup removes expired entries
func (mc *MetricsCache) cleanup() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	for key, entry := range mc.store {
		if now.After(entry.ExpiresAt) {
			delete(mc.store, key)
			mc.stats.TotalSize -= entry.DataSize
			mc.stats.Evictions++
		}
	}
	mc.stats.TotalEntries = len(mc.store)
}

// evictLRU evicts the least recently used entry
func (mc *MetricsCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range mc.store {
		if oldestKey == "" || entry.LastAccessed.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastAccessed
		}
	}

	if oldestKey != "" {
		entry := mc.store[oldestKey]
		delete(mc.store, oldestKey)
		mc.stats.TotalSize -= entry.DataSize
		mc.stats.Evictions++
	}
}

// updateStats updates performance statistics
func (mc *MetricsCache) updateStats(latency time.Duration) {
	// Update running average of latency
	totalRequests := mc.stats.Hits + mc.stats.Misses + 1
	if totalRequests == 1 {
		mc.stats.AverageLatency = latency
	} else {
		// Running average calculation
		currentAvg := mc.stats.AverageLatency
		mc.stats.AverageLatency = time.Duration(
			(int64(currentAvg)*int64(totalRequests-1) + int64(latency)) / int64(totalRequests),
		)
	}
}

// CacheMiddleware provides caching for metrics calculators
type CacheMiddleware struct {
	cache *MetricsCache
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(cache *MetricsCache) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
	}
}

// CacheOrCompute executes computation with caching
func (cm *CacheMiddleware) CacheOrCompute(
	ctx context.Context,
	metricType string,
	request MetricsRequest,
	computeFunc func() (interface{}, error),
) (interface{}, CacheInfo, error) {
	start := time.Now()

	// Check if caching is disabled
	if !request.UseCache {
		result, err := computeFunc()
		return result, CacheInfo{
			CacheHit:      false,
			ComputeTimeMs: time.Since(start).Milliseconds(),
		}, err
	}

	// Generate cache key
	cacheKey := cm.cache.GenerateCacheKey(metricType, request)

	// Try to get from cache
	if entry, found := cm.cache.Get(ctx, cacheKey); found {
		return entry.Data, CacheInfo{
			CacheHit:      true,
			CacheKey:      cacheKey,
			CachedAt:      &entry.CachedAt,
			TTL:           entry.TTL,
			ExpiresAt:     &entry.ExpiresAt,
			ComputeTimeMs: time.Since(start).Milliseconds(),
		}, nil
	}

	// Compute the result
	computeStart := time.Now()
	result, err := computeFunc()
	computeTime := time.Since(computeStart)

	if err != nil {
		return nil, CacheInfo{
			CacheHit:      false,
			ComputeTimeMs: time.Since(start).Milliseconds(),
		}, err
	}

	// Cache the result
	ttl := request.CacheTTL
	if ttl == 0 {
		ttl = cm.cache.GetTTLForMetricType(metricType)
	}

	tags := []string{
		metricType,
		fmt.Sprintf("repo:%s", request.Repository.FullName),
		fmt.Sprintf("granularity:%s", request.Granularity),
	}

	now := time.Now()
	expiresAt := now.Add(ttl)

	err = cm.cache.Set(ctx, cacheKey, result, ttl, computeTime, tags)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cache result: %v\n", err)
	}

	return result, CacheInfo{
		CacheHit:      false,
		CacheKey:      cacheKey,
		CachedAt:      &now,
		TTL:           ttl,
		ExpiresAt:     &expiresAt,
		ComputeTimeMs: time.Since(start).Milliseconds(),
	}, nil
}

// InvalidateRepositoryCache invalidates all cache entries for a repository
func (cm *CacheMiddleware) InvalidateRepositoryCache(ctx context.Context, repoFullName string) int {
	return cm.cache.InvalidateByTags(ctx, []string{fmt.Sprintf("repo:%s", repoFullName)})
}

// InvalidateMetricTypeCache invalidates all cache entries for a metric type
func (cm *CacheMiddleware) InvalidateMetricTypeCache(ctx context.Context, metricType string) int {
	return cm.cache.InvalidateByTags(ctx, []string{metricType})
}

// GetCacheStats returns cache performance statistics
func (cm *CacheMiddleware) GetCacheStats() CacheStats {
	return cm.cache.GetStats()
}

// WarmupCache pre-computes and caches common metrics
func (cm *CacheMiddleware) WarmupCache(ctx context.Context, repositories []string, metricTypes []string) error {
	// This would be implemented to pre-compute common metric combinations
	// for better user experience
	return nil
}