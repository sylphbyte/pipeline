package middleware

import (
	"fmt"

	pipe "github.com/sylphbyte/pipeline"
)

// Recovery Panic 恢复中间件
// 捕获 Hook 中的 panic（简化版，不依赖 Logger）
func Recovery[C pipe.Context, Option any, Payload any, Result any]() pipe.Middleware[C, Option, Payload, Result] {
	return func(next pipe.HookHandler[C, Option, Payload, Result]) pipe.HookHandler[C, Option, Payload, Result] {
		return func(ctx C, pipeCtx *pipe.PipeContext[Option, Payload, Result]) error {
			defer func() {
				if r := recover(); r != nil {
					// 捕获 panic，但不做处理（仅防止程序崩溃）
					_ = r
				}
			}()

			// 执行下一个 Handler
			err := next(ctx, pipeCtx)

			return err
		}
	}
}

// RecoveryWithError Panic 恢复中间件（将 panic 转换为 error）
func RecoveryWithError[C pipe.Context, Option any, Payload any, Result any]() pipe.Middleware[C, Option, Payload, Result] {
	return func(next pipe.HookHandler[C, Option, Payload, Result]) pipe.HookHandler[C, Option, Payload, Result] {
		return func(ctx C, pipeCtx *pipe.PipeContext[Option, Payload, Result]) error {
			var err error

			defer func() {
				if r := recover(); r != nil {
					// 将 panic 转换为 error
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()

			// 执行下一个 Handler
			err = next(ctx, pipeCtx)

			return err
		}
	}
}
