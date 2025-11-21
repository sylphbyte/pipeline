package pipeline

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sylphbyte/sylph"
)

// 测试用的泛型参数定义
type TestOption struct {
	EnableCache bool
	MaxRetries  int
}

type TestPayload struct {
	UserID int
	Data   string
}

type TestResult struct {
	Output   []string
	Metadata map[string]any
}

// 测试 Hook
func validateHook(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
	if pipeCtx.Payload.UserID <= 0 {
		return errors.New("invalid user ID")
	}
	return nil
}

func processHook(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
	pipeCtx.Result.Output = append(pipeCtx.Result.Output, pipeCtx.Payload.Data)
	pipeCtx.Set("processed", true)
	return nil
}

func errorHook(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
	return errors.New("test error")
}

// TestBasicPipelineExecution 测试基本管道执行
func TestBasicPipelineExecution(t *testing.T) {
	pipeline := NewPipeline[TestOption, TestPayload, TestResult](
		"test-pipeline",
		func(opt *TestOption) {
			opt.EnableCache = true
			opt.MaxRetries = 3
		},
	).AddHook(validateHook).AddHook(processHook)

	payload := &TestPayload{
		UserID: 123,
		Data:   "test data",
	}

	result, err := pipeline.Execute(newMockContext(), payload)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.Output) != 1 {
		t.Errorf("Expected 1 output, got %d", len(result.Output))
	}

	if result.Output[0] != "test data" {
		t.Errorf("Expected 'test data', got '%s'", result.Output[0])
	}
}

// TestHookExecutionOrder 测试 Hook 执行顺序
func TestHookExecutionOrder(t *testing.T) {
	var executionOrder []string

	hook1 := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executionOrder = append(executionOrder, "hook1")
		return nil
	}

	hook2 := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executionOrder = append(executionOrder, "hook2")
		return nil
	}

	hook3 := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executionOrder = append(executionOrder, "hook3")
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		AddHook(hook1).
		AddHook(hook2).
		AddHook(hook3)

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(executionOrder) != 3 {
		t.Fatalf("Expected 3 hooks executed, got %d", len(executionOrder))
	}

	if executionOrder[0] != "hook1" || executionOrder[1] != "hook2" || executionOrder[2] != "hook3" {
		t.Errorf("Hooks executed in wrong order: %v", executionOrder)
	}
}

// TestAbort 测试 Abort 功能
func TestAbort(t *testing.T) {
	var executed []string

	abortHook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executed = append(executed, "abort")
		pipeCtx.Abort()
		return nil
	}

	shouldNotExecute := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executed = append(executed, "should-not-execute")
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		AddHook(abortHook).
		AddHook(shouldNotExecute)

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(executed) != 1 {
		t.Errorf("Expected 1 hook executed, got %d", len(executed))
	}

	if executed[0] != "abort" {
		t.Errorf("Expected 'abort', got '%s'", executed[0])
	}
}

// TestDataSharing 测试数据共享
func TestDataSharing(t *testing.T) {
	setHook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		pipeCtx.Set("key1", "value1")
		pipeCtx.Set("key2", 123)
		return nil
	}

	getHook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		val1, ok1 := pipeCtx.Get("key1")
		val2, ok2 := pipeCtx.Get("key2")

		if !ok1 || val1 != "value1" {
			t.Errorf("Expected key1='value1', got ok=%v, val=%v", ok1, val1)
		}

		if !ok2 || val2 != 123 {
			t.Errorf("Expected key2=123, got ok=%v, val=%v", ok2, val2)
		}

		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		AddHook(setHook).
		AddHook(getHook)

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// TestMustGetPanic 测试 MustGet 在 key 不存在时 panic
func TestMustGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, but didn't get one")
		}
	}()

	hook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		_ = pipeCtx.MustGet("nonexistent")
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").AddHook(hook)
	_, _ = pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test-pipeline").
		AddHook(validateHook).
		AddHook(errorHook)

	payload := &TestPayload{UserID: 123}
	_, err := pipeline.Execute(newMockContext(), payload)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	pipeErr, ok := err.(*PipeError)
	if !ok {
		t.Fatalf("Expected *PipeError, got %T", err)
	}

	if pipeErr.PipelineName != "test-pipeline" {
		t.Errorf("Expected pipeline name 'test-pipeline', got '%s'", pipeErr.PipelineName)
	}

	if pipeErr.HookIndex != 1 {
		t.Errorf("Expected hook index 1, got %d", pipeErr.HookIndex)
	}
}

// TestSkipOnError 测试 SkipOnError 功能
func TestSkipOnError(t *testing.T) {
	var executed []string

	errorHook := NewHook(func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executed = append(executed, "error-hook")
		return errors.New("test error")
	}).SkipOnError().Build()

	successHook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		executed = append(executed, "success-hook")
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		AddHookWithOptions(errorHook).
		AddHook(successHook)

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(executed) != 2 {
		t.Errorf("Expected 2 hooks executed, got %d", len(executed))
	}
}

// TestLifecycleHooks 测试生命周期钩子
func TestLifecycleHooks(t *testing.T) {
	var lifecycle []string

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		OnBeforeExecute(func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) {
			lifecycle = append(lifecycle, "before")
		}).
		OnAfterExecute(func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult], err error) {
			lifecycle = append(lifecycle, "after")
		}).
		AddHook(func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
			lifecycle = append(lifecycle, "hook")
			return nil
		})

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(lifecycle) != 3 {
		t.Fatalf("Expected 3 lifecycle events, got %d", len(lifecycle))
	}

	if lifecycle[0] != "before" || lifecycle[1] != "hook" || lifecycle[2] != "after" {
		t.Errorf("Lifecycle hooks executed in wrong order: %v", lifecycle)
	}
}

// TestOnError 测试 OnError 钩子
func TestOnError(t *testing.T) {
	var errorCalled bool
	var capturedHookName string

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		OnError(func(ctx sylph.Context, hookName string, err error) {
			errorCalled = true
			capturedHookName = hookName
		}).
		AddNamedHook("failing-hook", errorHook)

	_, _ = pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})

	if !errorCalled {
		t.Error("OnError hook was not called")
	}

	if capturedHookName != "failing-hook" {
		t.Errorf("Expected hook name 'failing-hook', got '%s'", capturedHookName)
	}
}

// TestConcurrentSetGet 测试并发 Set/Get
func TestConcurrentSetGet(t *testing.T) {
	hook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(2)

			go func(n int) {
				defer wg.Done()
				key := string(rune('a' + n%26))
				pipeCtx.Set(key, n)
			}(i)

			go func(n int) {
				defer wg.Done()
				key := string(rune('a' + n%26))
				_, _ = pipeCtx.Get(key)
			}(i)
		}
		wg.Wait()
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").AddHook(hook)
	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})

	if err != nil {
		t.Fatalf("Concurrent access failed: %v", err)
	}
}

// TestConcurrentAbort 测试并发 Abort
func TestConcurrentAbort(t *testing.T) {
	hook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(2)

			go func() {
				defer wg.Done()
				pipeCtx.Abort()
			}()

			go func() {
				defer wg.Done()
				_ = pipeCtx.IsAborted()
			}()
		}
		wg.Wait()
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").AddHook(hook)
	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})

	if err != nil {
		t.Fatalf("Concurrent abort failed: %v", err)
	}
}

// TestExecutionStats 测试执行统计
func TestExecutionStats(t *testing.T) {
	var stats *ExecutionStats

	slowHook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult]("test").
		AddNamedHook("slow-hook", slowHook).
		OnAfterExecute(func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult], err error) {
			stats = pipeCtx.Stats()
		})

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if stats == nil {
		t.Fatal("Stats not captured")
	}

	if stats.PipelineName != "test" {
		t.Errorf("Expected pipeline name 'test', got '%s'", stats.PipelineName)
	}

	if len(stats.HookStats) != 1 {
		t.Errorf("Expected 1 hook stat, got %d", len(stats.HookStats))
	}

	if stats.HookStats[0].Name != "slow-hook" {
		t.Errorf("Expected hook name 'slow-hook', got '%s'", stats.HookStats[0].Name)
	}

	if stats.HookStats[0].Duration < 10*time.Millisecond {
		t.Errorf("Expected duration >= 10ms, got %v", stats.HookStats[0].Duration)
	}

	if !stats.Success {
		t.Error("Expected success=true")
	}
}

// TestOptionPointer 测试 Option 指针传递
func TestOptionPointer(t *testing.T) {
	hook := func(ctx sylph.Context, pipeCtx *PipeContext[TestOption, TestPayload, TestResult]) error {
		if pipeCtx.Option == nil {
			t.Error("Option should not be nil")
		}
		if !pipeCtx.Option.EnableCache {
			t.Error("Expected EnableCache=true")
		}
		if pipeCtx.Option.MaxRetries != 5 {
			t.Errorf("Expected MaxRetries=5, got %d", pipeCtx.Option.MaxRetries)
		}
		return nil
	}

	pipeline := NewPipeline[TestOption, TestPayload, TestResult](
		"test",
		func(opt *TestOption) {
			opt.EnableCache = true
			opt.MaxRetries = 5
		},
	).AddHook(hook)

	_, err := pipeline.Execute(newMockContext(), &TestPayload{UserID: 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
