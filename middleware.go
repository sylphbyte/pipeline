package pipe

// GenericHandler 是一个通用的处理器函数类型别名
// 用于简化 middleware 的函数签名
type GenericHandler func(ctx Context, pipeCtx interface{}) error

// Middleware 中间件定义
// 中间件可以包装 Hook 的执行，添加额外的功能（如日志、超时、重试等）
type Middleware[Option any, Payload any, Result any] func(next HookHandler[Option, Payload, Result]) HookHandler[Option, Payload, Result]

// applyMiddlewares 应用中间件到 Handler
// 中间件按照从左到右的顺序应用，最后一个中间件最先执行
func applyMiddlewares[Option any, Payload any, Result any](
	handler HookHandler[Option, Payload, Result],
	middlewares []Middleware[Option, Payload, Result],
) HookHandler[Option, Payload, Result] {
	// 从后往前应用中间件
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
