package pipeline

import "time"

// ExecutionStats 管道执行统计信息
type ExecutionStats struct {
	PipelineName  string        // 管道名称
	HookStats     []HookStat    // 各个 Hook 的统计
	TotalDuration time.Duration // 总执行时间
	StartTime     time.Time     // 开始时间
	EndTime       time.Time     // 结束时间
	Success       bool          // 是否成功
	Error         error         // 错误信息（如果有）
}

// HookStat Hook 执行统计
type HookStat struct {
	Name      string        // Hook 名称
	Index     int           // Hook 索引
	Duration  time.Duration // 执行时长
	Error     error         // 错误（如果有）
	StartTime time.Time     // 开始时间
	EndTime   time.Time     // 结束时间
}

// AddHookStat 添加 Hook 统计
func (s *ExecutionStats) AddHookStat(stat HookStat) {
	s.HookStats = append(s.HookStats, stat)
}

// MarkStart 标记开始时间
func (s *ExecutionStats) MarkStart() {
	s.StartTime = time.Now()
}

// MarkEnd 标记结束时间
func (s *ExecutionStats) MarkEnd(err error) {
	s.EndTime = time.Now()
	s.TotalDuration = s.EndTime.Sub(s.StartTime)
	s.Success = err == nil
	s.Error = err
}

// NewExecutionStats 创建执行统计
func NewExecutionStats(pipelineName string) *ExecutionStats {
	return &ExecutionStats{
		PipelineName: pipelineName,
		HookStats:    make([]HookStat, 0),
	}
}
