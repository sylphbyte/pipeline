package middleware

import (
	"fmt"
	"time"

	pipe "github.com/sylphbyte/pipeline"
)

// RetryFunc 重试中间件生成函数
// maxRetries: 最大重试次数
// backoff: 退避时间（每次重试会递增）
func RetryFunc[C pipe.Context, Option any, Payload any, Result any](
	maxRetries int,
	backoff time.Duration,
) pipe.Middleware[C, Option, Payload, Result] {
	return func(next pipe.HookHandler[C, Option, Payload, Result]) pipe.HookHandler[C, Option, Payload, Result] {
		return func(ctx C, pipeCtx *pipe.PipeContext[Option, Payload, Result]) error {
			var err error

			for i := 0; i <= maxRetries; i++ {
				// 执行 Handler
				err = next(ctx, pipeCtx)

				// 如果成功，立即返回
				if err == nil {
					return nil
				}

				// 如果不是最后一次尝试，等待后重试
				if i < maxRetries {
					waitTime := backoff * time.Duration(i+1)
					time.Sleep(waitTime)
				}
			}

			return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
		}
	}
}

// Retry 重试中间件（默认重试 3 次，退避 100ms）
func Retry[C pipe.Context, Option any, Payload any, Result any]() pipe.Middleware[C, Option, Payload, Result] {
	return RetryFunc[C, Option, Payload, Result](3, 100*time.Millisecond)
}
