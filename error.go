package pipeline

import "fmt"

// PipeError 管道执行错误
type PipeError struct {
	PipelineName string // 管道名称
	HookName     string // Hook 名称
	HookIndex    int    // Hook 索引
	Err          error  // 原始错误
}

func (e *PipeError) Error() string {
	if e.HookName != "" {
		return fmt.Sprintf(
			"pipeline '%s' failed at hook '%s' (index %d): %v",
			e.PipelineName, e.HookName, e.HookIndex, e.Err,
		)
	}
	return fmt.Sprintf(
		"pipeline '%s' failed at hook index %d: %v",
		e.PipelineName, e.HookIndex, e.Err,
	)
}

func (e *PipeError) Unwrap() error {
	return e.Err
}

// newPipeError 创建管道错误
func newPipeError(pipelineName, hookName string, hookIndex int, err error) *PipeError {
	return &PipeError{
		PipelineName: pipelineName,
		HookName:     hookName,
		HookIndex:    hookIndex,
		Err:          err,
	}
}
