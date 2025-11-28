package metrics

import (
	"sync"
	"time"
)

// Metrics 指标收集器接口
type Metrics interface {
	// 记录请求
	RecordRequest(route string, statusCode int, duration time.Duration)
	// 记录错误
	RecordError(route string, errorType string)
	// 获取统计信息（返回 interface{} 以兼容 MetricsProvider）
	GetStats() interface{}
	// 重置统计信息
	Reset()
}

// Stats 统计信息
type Stats struct {
	// 总请求数
	TotalRequests int64 `json:"total_requests"`
	// 成功请求数 (2xx)
	SuccessRequests int64 `json:"success_requests"`
	// 客户端错误数 (4xx)
	ClientErrors int64 `json:"client_errors"`
	// 服务端错误数 (5xx)
	ServerErrors int64 `json:"server_errors"`
	// 平均响应时间（毫秒）
	AvgResponseTime float64 `json:"avg_response_time_ms"`
	// P95 响应时间（毫秒）
	P95ResponseTime float64 `json:"p95_response_time_ms"`
	// P99 响应时间（毫秒）
	P99ResponseTime float64 `json:"p99_response_time_ms"`
	// 按路由分类的统计
	RouteStats map[string]*RouteStats `json:"route_stats"`
	// 错误统计
	ErrorStats map[string]int64 `json:"error_stats"`
	// 统计开始时间
	StartTime time.Time `json:"start_time"`
	// 最后更新时间
	LastUpdate time.Time `json:"last_update"`
}

// RouteStats 路由统计信息
type RouteStats struct {
	// 请求数
	Count int64 `json:"count"`
	// 平均响应时间（毫秒）
	AvgDuration float64 `json:"avg_duration_ms"`
	// 最小响应时间（毫秒）
	MinDuration int64 `json:"min_duration_ms"`
	// 最大响应时间（毫秒）
	MaxDuration int64 `json:"max_duration_ms"`
	// 成功数
	SuccessCount int64 `json:"success_count"`
	// 错误数
	ErrorCount int64 `json:"error_count"`
}

// memoryMetrics 内存指标收集器实现
type memoryMetrics struct {
	mu sync.RWMutex

	// 总请求数
	totalRequests int64
	// 成功请求数
	successRequests int64
	// 客户端错误数
	clientErrors int64
	// 服务端错误数
	serverErrors int64

	// 响应时间记录（用于计算百分位数）
	responseTimes []int64
	// 响应时间总和（用于计算平均值）
	totalDuration int64

	// 按路由分类的统计
	routeStats map[string]*routeStatsInternal

	// 错误统计
	errorStats map[string]int64

	// 统计开始时间
	startTime time.Time
	// 最后更新时间
	lastUpdate time.Time

	// 配置
	maxResponseTimeSamples int // 最多保留的响应时间样本数
}

// routeStatsInternal 路由统计内部结构
type routeStatsInternal struct {
	count         int64
	totalDuration int64
	minDuration   int64
	maxDuration   int64
	successCount  int64
	errorCount    int64
}

// Config 指标配置
type Config struct {
	// 最多保留的响应时间样本数（用于计算百分位数）
	MaxResponseTimeSamples int
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		MaxResponseTimeSamples: 10000, // 保留最近 10000 个样本
	}
}

// New 创建新的指标收集器
func New(config Config) Metrics {
	if config.MaxResponseTimeSamples <= 0 {
		config.MaxResponseTimeSamples = 10000
	}

	return &memoryMetrics{
		responseTimes:          make([]int64, 0, config.MaxResponseTimeSamples),
		routeStats:             make(map[string]*routeStatsInternal),
		errorStats:             make(map[string]int64),
		startTime:              time.Now(),
		lastUpdate:             time.Now(),
		maxResponseTimeSamples: config.MaxResponseTimeSamples,
	}
}

// RecordRequest 记录请求
// 需求: 7.4 - 统计各类请求的数量和平均响应时间
func (m *memoryMetrics) RecordRequest(route string, statusCode int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	durationMs := duration.Milliseconds()

	// 更新总计数
	m.totalRequests++
	m.totalDuration += durationMs
	m.lastUpdate = time.Now()

	// 更新状态码统计
	if statusCode >= 200 && statusCode < 300 {
		m.successRequests++
	} else if statusCode >= 400 && statusCode < 500 {
		m.clientErrors++
	} else if statusCode >= 500 {
		m.serverErrors++
	}

	// 记录响应时间（用于百分位数计算）
	if len(m.responseTimes) < m.maxResponseTimeSamples {
		m.responseTimes = append(m.responseTimes, durationMs)
	} else {
		// 如果超过最大样本数，使用循环缓冲区
		// 简单实现：移除最旧的样本
		m.responseTimes = append(m.responseTimes[1:], durationMs)
	}

	// 更新路由统计
	if _, exists := m.routeStats[route]; !exists {
		m.routeStats[route] = &routeStatsInternal{
			minDuration: durationMs,
			maxDuration: durationMs,
		}
	}

	routeStat := m.routeStats[route]
	routeStat.count++
	routeStat.totalDuration += durationMs

	if durationMs < routeStat.minDuration {
		routeStat.minDuration = durationMs
	}
	if durationMs > routeStat.maxDuration {
		routeStat.maxDuration = durationMs
	}

	if statusCode >= 200 && statusCode < 400 {
		routeStat.successCount++
	} else {
		routeStat.errorCount++
	}
}

// RecordError 记录错误
func (m *memoryMetrics) RecordError(route string, errorType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := route + ":" + errorType
	m.errorStats[key]++
	m.lastUpdate = time.Now()
}

// GetStats 获取统计信息
// 需求: 7.5 - 返回系统状态和关键指标快照
func (m *memoryMetrics) GetStats() interface{} {
	return m.getStatsInternal()
}

// getStatsInternal 内部方法，返回具体类型
func (m *memoryMetrics) getStatsInternal() *Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := &Stats{
		TotalRequests:   m.totalRequests,
		SuccessRequests: m.successRequests,
		ClientErrors:    m.clientErrors,
		ServerErrors:    m.serverErrors,
		RouteStats:      make(map[string]*RouteStats),
		ErrorStats:      make(map[string]int64),
		StartTime:       m.startTime,
		LastUpdate:      m.lastUpdate,
	}

	// 计算平均响应时间
	if m.totalRequests > 0 {
		stats.AvgResponseTime = float64(m.totalDuration) / float64(m.totalRequests)
	}

	// 计算百分位数
	if len(m.responseTimes) > 0 {
		stats.P95ResponseTime = calculatePercentile(m.responseTimes, 0.95)
		stats.P99ResponseTime = calculatePercentile(m.responseTimes, 0.99)
	}

	// 复制路由统计
	for route, routeStat := range m.routeStats {
		avgDuration := float64(0)
		if routeStat.count > 0 {
			avgDuration = float64(routeStat.totalDuration) / float64(routeStat.count)
		}

		stats.RouteStats[route] = &RouteStats{
			Count:        routeStat.count,
			AvgDuration:  avgDuration,
			MinDuration:  routeStat.minDuration,
			MaxDuration:  routeStat.maxDuration,
			SuccessCount: routeStat.successCount,
			ErrorCount:   routeStat.errorCount,
		}
	}

	// 复制错误统计
	for key, count := range m.errorStats {
		stats.ErrorStats[key] = count
	}

	return stats
}

// Reset 重置统计信息
func (m *memoryMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests = 0
	m.successRequests = 0
	m.clientErrors = 0
	m.serverErrors = 0
	m.totalDuration = 0
	m.responseTimes = make([]int64, 0, m.maxResponseTimeSamples)
	m.routeStats = make(map[string]*routeStatsInternal)
	m.errorStats = make(map[string]int64)
	m.startTime = time.Now()
	m.lastUpdate = time.Now()
}

// calculatePercentile 计算百分位数
func calculatePercentile(values []int64, percentile float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// 复制并排序
	sorted := make([]int64, len(values))
	copy(sorted, values)

	// 简单的冒泡排序（对于小数据集足够）
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 计算百分位数索引
	index := int(float64(len(sorted)-1) * percentile)
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return float64(sorted[index])
}
