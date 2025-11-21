package middleware

import (
	"app/pkg/pipe"
	"time"
)

// Logging 日志中间件
// 记录每个 Hook 的执行情况（简化版，不依赖 Logger）
func Logging[Option any, Payload any, Result any]() func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return func(next pipe.GenericHandler) pipe.GenericHandler {
		return func(ctx pipe.Context, pipeCtx interface{}) error {
			start := time.Now()

			// 执行下一个 Handler
			err := next(ctx, pipeCtx)

			_ = time.Since(start) // 记录执行时间（可用于统计）

			return err
		}
	}
}
