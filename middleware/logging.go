package middleware

import (
	"time"

	pipe "github.com/sylphbyte/pipeline"
)

// Logging 日志中间件
// 记录每个 Hook 的执行情况（简化版，不依赖 Logger）
func Logging[C pipe.Context, Option any, Payload any, Result any]() pipe.Middleware[C, Option, Payload, Result] {
	return func(next pipe.HookHandler[C, Option, Payload, Result]) pipe.HookHandler[C, Option, Payload, Result] {
		return func(ctx C, pipeCtx *pipe.PipeContext[Option, Payload, Result]) error {
			start := time.Now()

			// 执行下一个 Handler
			err := next(ctx, pipeCtx)

			_ = time.Since(start) // 记录执行时间（可用于统计）

			return err
		}
	}
}
