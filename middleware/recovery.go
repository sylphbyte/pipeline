package middleware

import (
	"app/pkg/pipe"
	"fmt"

)

// Recovery Panic 恢复中间件
// 捕获 Hook 中的 panic（简化版，不依赖 Logger）
func Recovery[Option any, Payload any, Result any]() func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return func(next pipe.GenericHandler) pipe.GenericHandler {
		return func(ctx pipe.Context, pipeCtx interface{}) error {
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
func RecoveryWithError[Option any, Payload any, Result any]() func(
	next pipe.GenericHandler,
) pipe.GenericHandler {
	return func(next pipe.GenericHandler) pipe.GenericHandler {
		return func(ctx pipe.Context, pipeCtx interface{}) error {
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
