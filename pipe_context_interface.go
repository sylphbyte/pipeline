package pipeline

import "context"

// Context 定义管道执行所需的最小上下文接口
// 用户可以传入任何实现了该接口的类型，包括 sylph.Context
type Context interface {
	// 嵌入标准 context.Context，支持超时、取消、传值等标准功能
	context.Context

	// 日志方法：支持结构化日志记录
	Info(pkg, action string, data any)
	Warn(pkg, action string, data any)
	Error(pkg, action string, err error, data any)
	Debug(pkg, action string, data any)
}

// stdContextAdapter 是标准 context.Context 的适配器
// 为不需要日志功能的场景提供默认实现
type stdContextAdapter struct {
	context.Context
}

func (s *stdContextAdapter) Info(pkg, action string, data any)             {}
func (s *stdContextAdapter) Warn(pkg, action string, data any)             {}
func (s *stdContextAdapter) Error(pkg, action string, err error, data any) {}
func (s *stdContextAdapter) Debug(pkg, action string, data any)            {}

// WrapContext 将标准 context.Context 包装为 pipe.Context
// 日志方法默认为空操作（no-op）
func WrapContext(ctx context.Context) Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &stdContextAdapter{Context: ctx}
}
