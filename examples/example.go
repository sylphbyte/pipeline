package examples

import (
	"fmt"

	pipe "github.com/sylphbyte/pipeline"
	"github.com/sylphbyte/sylph"
)

// 示例定义泛型参数
type ExampleOption struct {
	EnableCache bool
	MaxRetries  int
}

type ExamplePayload struct {
	UserID int
	Data   string
}

type ExampleResult struct {
	Output   []string
	Metadata map[string]any
}

// 示例 Hook 1: 验证
func ValidateHook(ctx sylph.Context, pipeCtx *pipe.PipeContext[ExampleOption, ExamplePayload, ExampleResult]) error {
	if pipeCtx.Payload.UserID <= 0 {
		pipeCtx.Abort() // 中断后续流程
		return fmt.Errorf("invalid user ID: %d", pipeCtx.Payload.UserID)
	}
	return nil
}

// 示例 Hook 2: 处理数据
func ProcessHook(ctx sylph.Context, pipeCtx *pipe.PipeContext[ExampleOption, ExamplePayload, ExampleResult]) error {
	// 处理数据
	pipeCtx.Result.Output = append(pipeCtx.Result.Output, pipeCtx.Payload.Data)

	// 在 Hook 之间共享数据
	pipeCtx.Set("processed", true)
	pipeCtx.Set("processedCount", 1)

	return nil
}

// 示例 Hook 3: 生成元数据
func MetadataHook(ctx sylph.Context, pipeCtx *pipe.PipeContext[ExampleOption, ExamplePayload, ExampleResult]) error {
	// 从共享数据中获取信息
	processed, _ := pipeCtx.Get("processed")
	processedCount, _ := pipeCtx.Get("processedCount")

	// 生成元数据
	pipeCtx.Result.Metadata = map[string]any{
		"processed":      processed,
		"processedCount": processedCount,
		"cacheEnabled":   pipeCtx.Option.EnableCache,
	}

	return nil
}

// 基本使用示例
func BasicExample(ctx sylph.Context) {
	// 创建管道
	pipeline := pipe.NewPipeline[sylph.Context, ExampleOption, ExamplePayload, ExampleResult](
		"example-pipeline",
		func(opt *ExampleOption) {
			opt.EnableCache = true
			opt.MaxRetries = 3
		},
	).
		AddHook(ValidateHook).
		AddNamedHook("process", ProcessHook).
		AddNamedHook("metadata", MetadataHook)

	// 执行管道
	payload := &ExamplePayload{
		UserID: 123,
		Data:   "test data",
	}

	result, err := pipeline.Execute(ctx, payload)
	if err != nil {
		fmt.Printf("Pipeline failed: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
}

// 使用中间件示例
func MiddlewareExample(ctx sylph.Context) {
	// 导入中间件包
	// import "app/pkg/pipe/middleware"

	pipeline := pipe.NewPipeline[sylph.Context, ExampleOption, ExamplePayload, ExampleResult](
		"middleware-pipeline",
	).
		// 使用内置中间件
		// Use(middleware.RecoveryWithError[ExampleOption, ExamplePayload, ExampleResult]()).
		// Use(middleware.Logging[ExampleOption, ExamplePayload, ExampleResult]()).
		// Use(middleware.RetryFunc[ExampleOption, ExamplePayload, ExampleResult](3, 100*time.Millisecond)).
		AddHook(ValidateHook).
		AddHook(ProcessHook)

	payload := &ExamplePayload{UserID: 123, Data: "test"}
	result, err := pipeline.Execute(ctx, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Success: %+v\n", result)
}

// 生命周期钩子示例
func LifecycleExample(ctx sylph.Context) {
	pipeline := pipe.NewPipeline[sylph.Context, ExampleOption, ExamplePayload, ExampleResult](
		"lifecycle-pipeline",
	).
		OnBeforeExecute(func(ctx sylph.Context, pipeCtx *pipe.PipeContext[ExampleOption, ExamplePayload, ExampleResult]) {
			fmt.Println("Pipeline started")
		}).
		OnAfterExecute(func(ctx sylph.Context, pipeCtx *pipe.PipeContext[ExampleOption, ExamplePayload, ExampleResult], err error) {
			stats := pipeCtx.Stats()
			fmt.Printf("Pipeline completed in %v\n", stats.TotalDuration)

			// 打印每个 Hook 的执行时间
			for _, hookStat := range stats.HookStats {
				fmt.Printf("Hook '%s': %v\n", hookStat.Name, hookStat.Duration)
			}
		}).
		OnError(func(ctx sylph.Context, hookName string, err error) {
			fmt.Printf("Hook '%s' failed: %v\n", hookName, err)
		}).
		AddHook(ValidateHook).
		AddHook(ProcessHook)

	payload := &ExamplePayload{UserID: 123, Data: "test"}
	_, _ = pipeline.Execute(ctx, payload)
}

// 高级 Hook 配置示例
func AdvancedHookExample(ctx sylph.Context) {
	// 使用 Hook 构建器
	processHook := pipe.NewHook[sylph.Context, ExampleOption, ExamplePayload, ExampleResult](ProcessHook).
		WithName("process-data").
		WithDescription("处理用户数据").
		// WithTimeout(5 * time.Second).
		SkipOnError(). // 错误时跳过而非中断
		Build()

	pipeline := pipe.NewPipeline[sylph.Context, ExampleOption, ExamplePayload, ExampleResult](
		"advanced-pipeline",
	).
		AddHook(ValidateHook).
		AddHookWithOptions(processHook).
		AddHook(MetadataHook)

	payload := &ExamplePayload{UserID: 123, Data: "test"}
	result, err := pipeline.Execute(ctx, payload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
}
