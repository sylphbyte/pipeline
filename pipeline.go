// Package pipe 提供了一个功能完善的泛型管道引擎，用于编排复杂的业务流程。
//
// # 主要特性
//
//   - 泛型支持：完全类型安全的管道定义
//   - 中间件系统：灵活的横切关注点处理（日志、超时、重试、恢复）
//   - 生命周期钩子：完善的流程控制（执行前、执行后、错误处理）
//   - 并发安全：使用 sync.RWMutex 保护共享数据
//   - 执行统计：内置性能监控和统计信息
//   - 错误处理：结构化的错误类型，支持错误链
//
// # 基本使用
//
// 定义泛型参数：
//
//	type MyOption struct {
//	    EnableCache bool
//	}
//	type MyPayload struct {
//	    UserID int
//	}
//	type MyResult struct {
//	    Output []string
//	}
//
// 创建管道并添加 Hook：
//
//	pipeline := pipe.NewPipeline[MyOption, MyPayload, MyResult](
//	    "my-pipeline",
//	    func(opt *MyOption) { opt.EnableCache = true },
//	).
//	    AddHook(ValidateHook).
//	    AddNamedHook("process", ProcessHook)
//
// 执行管道：
//
//	result, err := pipeline.Execute(ctx, &MyPayload{UserID: 123})
//
// # 中间件支持
//
// 使用内置中间件：
//
//	import "app/pkg/pipe/middleware"
//
//	pipeline.Use(
//	    middleware.RecoveryWithError[MyOption, MyPayload, MyResult](),
//	    middleware.Logging[MyOption, MyPayload, MyResult](),
//	    middleware.RetryFunc[MyOption, MyPayload, MyResult](3, 100*time.Millisecond),
//	)
//
// # 生命周期钩子
//
// 注册生命周期回调：
//
//	pipeline.
//	    OnBeforeExecute(func(ctx pipe.Context, pipeCtx *pipe.PipeContext[...]) {
//	        // 执行前逻辑
//	    }).
//	    OnAfterExecute(func(ctx pipe.Context, pipeCtx *pipe.PipeContext[...], err error) {
//	        stats := pipeCtx.Stats()
//	        // 查看执行统计
//	    }).
//	    OnError(func(ctx pipe.Context, hookName string, err error) {
//	        // 错误处理
//	    })
//
// 更多示例请参考 examples/example.go 文件。
package pipeline

import (
	"time"
)

// OptionHandler Option 配置处理函数
type OptionHandler[Option any] func(option *Option)

// Pipeline 管道定义
type Pipeline[C Context, Option any, Payload any, Result any] struct {
	Name   string  // 管道名称
	option *Option // 选项数据（指针类型，与 OptionHandler 一致）

	hooks       []*Hook[C, Option, Payload, Result]      // Hook 列表
	middlewares []Middleware[C, Option, Payload, Result] // 中间件列表

	// 生命周期钩子
	beforeExecute []func(ctx C, pipeCtx *PipeContext[Option, Payload, Result])
	afterExecute  []func(ctx C, pipeCtx *PipeContext[Option, Payload, Result], err error)
	onError       []func(ctx C, hookName string, err error)
}

// NewPipeline 创建新的管道
func NewPipeline[C Context, Option any, Payload any, Result any](
	name string,
	opts ...OptionHandler[Option],
) *Pipeline[C, Option, Payload, Result] {
	option := new(Option) // 使用 new() 创建泛型类型的零值指针
	for _, opt := range opts {
		opt(option)
	}

	return &Pipeline[C, Option, Payload, Result]{
		Name:          name,
		option:        option,
		hooks:         make([]*Hook[C, Option, Payload, Result], 0),
		middlewares:   make([]Middleware[C, Option, Payload, Result], 0),
		beforeExecute: make([]func(C, *PipeContext[Option, Payload, Result]), 0),
		afterExecute:  make([]func(C, *PipeContext[Option, Payload, Result], error), 0),
		onError:       make([]func(C, string, error), 0),
	}
}

// AddHook 添加 Hook（简化版，直接使用 Handler）
func (p *Pipeline[C, Option, Payload, Result]) AddHook(
	handlers ...HookHandler[C, Option, Payload, Result],
) *Pipeline[C, Option, Payload, Result] {
	for _, handler := range handlers {
		hook := &Hook[C, Option, Payload, Result]{
			Handler: handler,
		}
		p.hooks = append(p.hooks, hook)
	}
	return p
}

// AddNamedHook 添加命名的 Hook
func (p *Pipeline[C, Option, Payload, Result]) AddNamedHook(
	name string,
	handler HookHandler[C, Option, Payload, Result],
) *Pipeline[C, Option, Payload, Result] {
	hook := &Hook[C, Option, Payload, Result]{
		Name:    name,
		Handler: handler,
	}
	p.hooks = append(p.hooks, hook)
	return p
}

// AddHookWithOptions 添加带配置的 Hook
func (p *Pipeline[C, Option, Payload, Result]) AddHookWithOptions(
	hook *Hook[C, Option, Payload, Result],
) *Pipeline[C, Option, Payload, Result] {
	p.hooks = append(p.hooks, hook)
	return p
}

// Use 使用中间件
func (p *Pipeline[C, Option, Payload, Result]) Use(
	middlewares ...Middleware[C, Option, Payload, Result],
) *Pipeline[C, Option, Payload, Result] {
	p.middlewares = append(p.middlewares, middlewares...)
	return p
}

// OnBeforeExecute 注册执行前钩子
func (p *Pipeline[C, Option, Payload, Result]) OnBeforeExecute(
	fn func(ctx C, pipeCtx *PipeContext[Option, Payload, Result]),
) *Pipeline[C, Option, Payload, Result] {
	p.beforeExecute = append(p.beforeExecute, fn)
	return p
}

// OnAfterExecute 注册执行后钩子
func (p *Pipeline[C, Option, Payload, Result]) OnAfterExecute(
	fn func(ctx C, pipeCtx *PipeContext[Option, Payload, Result], err error),
) *Pipeline[C, Option, Payload, Result] {
	p.afterExecute = append(p.afterExecute, fn)
	return p
}

// OnError 注册错误处理钩子
func (p *Pipeline[C, Option, Payload, Result]) OnError(
	fn func(ctx C, hookName string, err error),
) *Pipeline[C, Option, Payload, Result] {
	p.onError = append(p.onError, fn)
	return p
}

// Execute 执行管道
func (p *Pipeline[C, Option, Payload, Result]) Execute(
	ctx C,
	payload *Payload,
) (*Result, error) {
	// 初始化 Result
	var result Result

	// 创建执行统计
	stats := NewExecutionStats(p.Name)
	stats.MarkStart()

	// 初始化 PipeContext
	pipeCtx := &PipeContext[Option, Payload, Result]{
		Name:    p.Name,
		Option:  p.option, // 指针传递，避免大结构体拷贝
		Payload: payload,
		Result:  &result,              // 指针传递，允许 Hook 修改
		data:    make(map[string]any), // 初始化中间状态
		stats:   stats,
	}

	// 执行 BeforeExecute 钩子
	for _, fn := range p.beforeExecute {
		fn(ctx, pipeCtx)
	}

	var finalErr error

	// 执行所有 Hook
	for i, hook := range p.hooks {
		// 检查是否中断
		if pipeCtx.IsAborted() {
			break
		}

		// 记录 Hook 开始时间
		hookStat := HookStat{
			Name:      hook.Name,
			Index:     i,
			StartTime: time.Now(),
		}

		// 应用中间件
		handler := hook.Handler
		if len(p.middlewares) > 0 {
			handler = applyMiddlewares(handler, p.middlewares)
		}

		// 执行 Hook
		err := handler(ctx, pipeCtx)

		// 记录 Hook 结束时间
		hookStat.EndTime = time.Now()
		hookStat.Duration = hookStat.EndTime.Sub(hookStat.StartTime)
		hookStat.Error = err
		stats.AddHookStat(hookStat)

		// 处理错误
		if err != nil {
			// 调用错误处理钩子
			for _, errFn := range p.onError {
				errFn(ctx, hook.Name, err)
			}

			// 如果设置了 SkipOnError，则跳过错误继续执行
			if hook.SkipOnError {
				continue
			}

			// 否则中断执行并返回错误
			finalErr = newPipeError(p.Name, hook.Name, i, err)
			break
		}
	}

	// 标记执行结束
	stats.MarkEnd(finalErr)

	// 执行 AfterExecute 钩子
	for _, fn := range p.afterExecute {
		fn(ctx, pipeCtx, finalErr)
	}

	if finalErr != nil {
		return nil, finalErr
	}

	return pipeCtx.Result, nil
}
