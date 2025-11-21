package pipe

import (
	"time"

	"github.com/sylphbyte/sylph"
)

// HookHandler Hook 处理函数
type HookHandler[Option any, Payload any, Result any] func(
	ctx Context,
	pipeCtx *PipeContext[Option, Payload, Result],
) error

// Hook Hook 定义，包含元数据
type Hook[Option any, Payload any, Result any] struct {
	Name        string                               // Hook 名称
	Description string                               // Hook 描述
	Handler     HookHandler[Option, Payload, Result] // 处理函数
	Timeout     time.Duration                        // 超时时间（0 表示无超时）
	SkipOnError bool                                 // 错误时是否跳过而非中断整个管道
}

// Execute 执行 Hook
func (h *Hook[Option, Payload, Result]) Execute(
	ctx sylph.Context,
	pipeCtx *PipeContext[Option, Payload, Result],
) error {
	return h.Handler(ctx, pipeCtx)
}

// HookBuilder Hook 构建器
type HookBuilder[Option any, Payload any, Result any] struct {
	hook *Hook[Option, Payload, Result]
}

// WithName 设置 Hook 名称
func (b *HookBuilder[Option, Payload, Result]) WithName(name string) *HookBuilder[Option, Payload, Result] {
	b.hook.Name = name
	return b
}

// WithDescription 设置 Hook 描述
func (b *HookBuilder[Option, Payload, Result]) WithDescription(desc string) *HookBuilder[Option, Payload, Result] {
	b.hook.Description = desc
	return b
}

// WithTimeout 设置超时时间
func (b *HookBuilder[Option, Payload, Result]) WithTimeout(timeout time.Duration) *HookBuilder[Option, Payload, Result] {
	b.hook.Timeout = timeout
	return b
}

// SkipOnError 设置错误时跳过
func (b *HookBuilder[Option, Payload, Result]) SkipOnError() *HookBuilder[Option, Payload, Result] {
	b.hook.SkipOnError = true
	return b
}

// Build 构建 Hook
func (b *HookBuilder[Option, Payload, Result]) Build() *Hook[Option, Payload, Result] {
	return b.hook
}

// NewHook 创建 Hook
func NewHook[Option any, Payload any, Result any](
	handler HookHandler[Option, Payload, Result],
) *HookBuilder[Option, Payload, Result] {
	return &HookBuilder[Option, Payload, Result]{
		hook: &Hook[Option, Payload, Result]{
			Handler: handler,
		},
	}
}
