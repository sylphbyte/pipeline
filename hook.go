package pipeline

import (
	"time"
)

// HookHandler Hook 处理函数
type HookHandler[C Context, Option any, Payload any, Result any] func(
	ctx C,
	pipeCtx *PipeContext[Option, Payload, Result],
) error

// Hook Hook 定义，包含元数据
type Hook[C Context, Option any, Payload any, Result any] struct {
	Name        string                                  // Hook 名称
	Description string                                  // Hook 描述
	Handler     HookHandler[C, Option, Payload, Result] // 处理函数
	Timeout     time.Duration                           // 超时时间（0 表示无超时）
	SkipOnError bool                                    // 错误时是否跳过而非中断整个管道
}

// Execute 执行 Hook
func (h *Hook[C, Option, Payload, Result]) Execute(
	ctx C,
	pipeCtx *PipeContext[Option, Payload, Result],
) error {
	return h.Handler(ctx, pipeCtx)
}

// HookBuilder Hook 构建器
type HookBuilder[C Context, Option any, Payload any, Result any] struct {
	hook *Hook[C, Option, Payload, Result]
}

// WithName 设置 Hook 名称
func (b *HookBuilder[C, Option, Payload, Result]) WithName(name string) *HookBuilder[C, Option, Payload, Result] {
	b.hook.Name = name
	return b
}

// WithDescription 设置 Hook 描述
func (b *HookBuilder[C, Option, Payload, Result]) WithDescription(desc string) *HookBuilder[C, Option, Payload, Result] {
	b.hook.Description = desc
	return b
}

// WithTimeout 设置超时时间
func (b *HookBuilder[C, Option, Payload, Result]) WithTimeout(timeout time.Duration) *HookBuilder[C, Option, Payload, Result] {
	b.hook.Timeout = timeout
	return b
}

// SkipOnError 设置错误时跳过
func (b *HookBuilder[C, Option, Payload, Result]) SkipOnError() *HookBuilder[C, Option, Payload, Result] {
	b.hook.SkipOnError = true
	return b
}

// Build 构建 Hook
func (b *HookBuilder[C, Option, Payload, Result]) Build() *Hook[C, Option, Payload, Result] {
	return b.hook
}

// NewHook 创建 Hook
func NewHook[C Context, Option any, Payload any, Result any](
	handler HookHandler[C, Option, Payload, Result],
) *HookBuilder[C, Option, Payload, Result] {
	return &HookBuilder[C, Option, Payload, Result]{
		hook: &Hook[C, Option, Payload, Result]{
			Handler: handler,
		},
	}
}
