package middleware

import (
	"app/pkg/pipe"
	"context"
	"fmt"
	"time"

)

// TimeoutFunc 超时中间件生成函数
// 为每个 Hook 添加超时控制
func TimeoutFunc[Option any, Payload any, Result any](timeout time.Duration) func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return func(next pipe.GenericHandler) pipe.GenericHandler {
		return func(ctx pipe.Context, pipeCtx interface{}) error {
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
func Timeout[Option any, Payload any, Result any]() func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return TimeoutFunc[Option, Payload, Result](30 * time.Second)
}
