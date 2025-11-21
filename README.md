# Go Pipeline - 泛型管道引擎

一个功能完善、生产就绪的 Go 泛型管道引擎，用于优雅地编排复杂的业务流程。

## 特性

✨ **泛型支持** - 完全类型安全的管道定义  
🎯 **中间件系统** - 灵活的横切关注点处理  
📊 **执行统计** - 内置的性能监控和统计  
🔧 **生命周期钩子** - 完善的流程控制  
⚡ **错误处理** - 结构化的错误信息  
🛡️ **生产就绪** - 包含超时、重试、恢复等企业级特性

## 快速开始
### 安装
```
go get github.com/sylphbyte/pipeline
```

### 基本使用

```go
package main

import (
    "app/pkg/pipe"
    "github.com/sylphbyte/sylph"
)

// 定义泛型参数
type MyOption struct {
    EnableCache bool
}

type MyPayload struct {
    UserID int
    Data   string
}

type MyResult struct {
    Output []string
}

// 定义 Hook
func ValidateHook(ctx sylph.Context, pipeCtx *pipe.PipeContext[MyOption, MyPayload, MyResult]) error {
    if pipeCtx.Payload.UserID <= 0 {
        pipeCtx.Abort() // 中断后续流程
        return fmt.Errorf("invalid user ID")
    }
    return nil
}

func ProcessHook(ctx sylph.Context, pipeCtx *pipe.PipeContext[MyOption, MyPayload, MyResult]) error {
    // 处理数据
    pipeCtx.Result.Output = append(pipeCtx.Result.Output, pipeCtx.Payload.Data)
    
    // 在 Hook 之间共享数据
    pipeCtx.Set("processed", true)
    return nil
}

func main() {
    // 创建管道
    pipeline := pipe.NewPipeline[MyOption, MyPayload, MyResult](
        "my-pipeline",
        func(opt *MyOption) {
            opt.EnableCache = true
        },
    ).
        AddHook(ValidateHook).
        AddNamedHook("process", ProcessHook)
    
    // 执行管道
    payload := &MyPayload{UserID: 123, Data: "test"}
    result, err := pipeline.Execute(ctx, payload)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result.Output)
}
```

### 使用中间件

```go
import "app/pkg/pipe/middleware"

pipeline := pipe.NewPipeline[MyOption, MyPayload, MyResult]("my-pipeline").
    // 使用内置中间件
    Use(middleware.RecoveryWithError[MyOption, MyPayload, MyResult]()).
    Use(middleware.Logging[MyOption, MyPayload, MyResult]()).
    Use(middleware.TimeoutFunc[MyOption, MyPayload, MyResult](5 * time.Second)).
    Use(middleware.RetryFunc[MyOption, MyPayload, MyResult](3, 100*time.Millisecond)).
    // 添加 Hook
    AddHook(ValidateHook).
    AddHook(ProcessHook)
```

### 生命周期钩子

```go
pipeline := pipe.NewPipeline[MyOption, MyPayload, MyResult]("my-pipeline").
    OnBeforeExecute(func(ctx sylph.Context, pipeCtx *pipe.PipeContext[...]) {
        ctx.Logger().Info("Pipeline started")
    }).
    OnAfterExecute(func(ctx sylph.Context, pipeCtx *pipe.PipeContext[...], err error) {
        stats := pipeCtx.Stats()
        ctx.Logger().Infof("Pipeline completed in %v", stats.TotalDuration)
    }).
    OnError(func(ctx sylph.Context, hookName string, err error) {
        ctx.Logger().Errorf("Hook '%s' failed: %v", hookName, err)
    }).
    AddHook(ValidateHook).
    AddHook(ProcessHook)
```

### 高级 Hook 配置

```go
// 使用 Hook 构建器
hook := pipe.NewHook(ProcessHook).
    WithName("process-data").
    WithDescription("处理用户数据").
    WithTimeout(5 * time.Second).
    SkipOnError().  // 错误时跳过而非中断
    Build()

pipeline := pipe.NewPipeline[MyOption, MyPayload, MyResult]("my-pipeline").
    AddHookWithOptions(hook)
```

### 执行统计

```go
result, err := pipeline.Execute(ctx, payload)
if err != nil {
    if pipeErr, ok := err.(*pipe.PipeError); ok {
        fmt.Printf("Pipeline: %s\n", pipeErr.PipelineName)
        fmt.Printf("Failed at: %s (index %d)\n", pipeErr.HookName, pipeErr.HookIndex)
    }
}

// 获取执行统计
stats := pipeCtx.Stats()
fmt.Printf("Total duration: %v\n", stats.TotalDuration)
for _, hookStat := range stats.HookStats {
    fmt.Printf("Hook '%s': %v\n", hookStat.Name, hookStat.Duration)
}
```

## 内置中间件

### Logging
记录每个 Hook 的执行时间和错误

```go
middleware.Logging[Option, Payload, Result]()
```

### Timeout
为 Hook 添加超时控制

```go
middleware.TimeoutFunc[Option, Payload, Result](30 * time.Second)
middleware.Timeout[Option, Payload, Result]() // 默认 30秒
```

### Retry
失败时自动重试

```go
middleware.RetryFunc[Option, Payload, Result](3, 100*time.Millisecond)
middleware.Retry[Option, Payload, Result]() // 默认重试3次
```

### Recovery
捕获 panic 并记录堆栈

```go
middleware.Recovery[Option, Payload, Result]()           // 仅记录
middleware.RecoveryWithError[Option, Payload, Result]()  // 转换为错误
```

## 最佳实践

### 1. 清晰的职责分离
每个 Hook 只负责一个具体的业务步骤

### 2. 使用命名 Hook
使用 `AddNamedHook` 为 Hook 命名，便于调试和监控

### 3. 合理使用中间件
将通用的横切逻辑（日志、超时、重试等）放在中间件中

### 4. 利用共享数据
使用 `pipeCtx.Set()` 和 `pipeCtx.Get()` 在 Hook 之间共享数据

### 5. 监控执行统计
在生产环境中收集和分析执行统计，优化性能

## 架构

```
pkg/pipe/
├── pipeline.go      # Pipeline 核心定义
├── context.go       # PipeContext 定义
├── hook.go          # Hook 相关定义
├── middleware.go    # 中间件系统
├── error.go         # 错误类型
├── stats.go         # 执行统计
└── middleware/      # 内置中间件
    ├── logging.go
    ├── timeout.go
    ├── retry.go
    └── recovery.go
```

## 实际应用场景

- **API 请求处理链** - 验证、鉴权、业务处理、响应
- **数据处理管道** - 清洗、转换、验证、存储
- **业务流程编排** - 订单处理、支付流程、审批流程
- **微服务集成** - 服务调用编排、错误处理、降级

## License

MIT
