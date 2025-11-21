package examples

import "testing"

func TestExamples(t *testing.T) {
	ctx := newMockContext()

	t.Run("BasicExample", func(t *testing.T) {
		BasicExample(ctx)
	})

	t.Run("MiddlewareExample", func(t *testing.T) {
		MiddlewareExample(ctx)
	})

	t.Run("LifecycleExample", func(t *testing.T) {
		LifecycleExample(ctx)
	})

	t.Run("AdvancedHookExample", func(t *testing.T) {
		AdvancedHookExample(ctx)
	})
}
