package pipe

import "sync"

// PipeContext 管道上下文，包含输入(Payload)和输出(Result)
// 所有的业务数据都放在 Payload 里
type PipeContext[Option any, Payload any, Result any] struct {
	Name    string   // 流程名称
	Option  *Option  // 选项数据（指针，避免大结构体拷贝）
	Payload *Payload // 输入数据：包含 User, Channel, Req, Plans 等
	Result  *Result  // 输出数据：Hook 将结果写到这里

	data  map[string]any // 中间状态：Hook 之间可以共享数据（私有，通过方法访问）
	abort bool           // 控制位：是否中断后续 Hook（私有，通过方法访问）
	mu    sync.RWMutex   // 保护 data 和 abort 的并发访问

	stats *ExecutionStats // 执行统计
}

// Abort 允许 Hook 中断流程（比如参数校验不通过）
func (p *PipeContext[Option, Payload, Result]) Abort() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.abort = true
}

// IsAborted 是否已中断
func (p *PipeContext[Option, Payload, Result]) IsAborted() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.abort
}

// Stats 获取执行统计
func (p *PipeContext[Option, Payload, Result]) Stats() *ExecutionStats {
	return p.stats
}

// Set 设置共享数据（并发安全）
func (p *PipeContext[Option, Payload, Result]) Set(key string, value any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[key] = value
}

// Get 获取共享数据（并发安全）
func (p *PipeContext[Option, Payload, Result]) Get(key string) (any, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.data[key]
	return val, ok
}

// MustGet 获取共享数据（不存在时 panic，并发安全）
func (p *PipeContext[Option, Payload, Result]) MustGet(key string) any {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.data[key]
	if !ok {
		panic("key not found: " + key)
	}
	return val
}
