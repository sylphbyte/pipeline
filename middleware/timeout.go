package middleware

import (
	"context"
	"fmt"
	"time"

	pipe "github.com/sylphbyte/pipeline"
)

// TimeoutFunc 超时中间件生成函数
// 为每个 Hook 添加超时控制
func TimeoutFunc[C pipe.Context, Option any, Payload any, Result any](timeout time.Duration) pipe.Middleware[C, Option, Payload, Result] {
	return func(next pipe.HookHandler[C, Option, Payload, Result]) pipe.HookHandler[C, Option, Payload, Result] {
		return func(ctx C, pipeCtx *pipe.PipeContext[Option, Payload, Result]) error {
			done := make(chan error, 1)

			// 在 goroutine 中执行
			go func() {
				done <- next(ctx, pipeCtx)
			}()

			// 使用 context 超时控制
			timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			select {
			case err := <-done:
				return err
			case <-timeoutCtx.Done():
				return fmt.Errorf("hook timeout after %v", timeout)
			}
		}
	}
}

// Timeout 超时中间件（默认 30 秒）
func Timeout[C pipe.Context, Option any, Payload any, Result any]() pipe.Middleware[C, Option, Payload, Result] {
	return TimeoutFunc[C, Option, Payload, Result](30 * time.Second)
}
