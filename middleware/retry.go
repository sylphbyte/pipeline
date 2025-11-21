package middleware

import (
	"app/pkg/pipe"
	"fmt"
	"time"
)

// RetryFunc 重试中间件生成函数
// maxRetries: 最大重试次数
// backoff: 退避时间（每次重试会递增）
func RetryFunc[Option any, Payload any, Result any](
	maxRetries int,
	backoff time.Duration,
) func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return func(next pipe.GenericHandler) pipe.GenericHandler {
		return func(ctx pipe.Context, pipeCtx interface{}) error {
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
func Retry[Option any, Payload any, Result any]() func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return RetryFunc[Option, Payload, Result](3, 100*time.Millisecond)
}
