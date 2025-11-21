package pipe

import (
	"time"

	"github.com/sylphbyte/sylph"
)

// Mock sylph.Context - 简化实现，仅用于测试
// 实现所有必需的 sylph.Context 接口方法
type mockContext struct{}

// Context methods
func (m *mockContext) Clone() sylph.Context                  { return m }
func (m *mockContext) Set(key string, value interface{})     {}
func (m *mockContext) MarkSet(key string, value interface{}) {}
func (m *mockContext) Get(key string) (interface{}, bool)    { return nil, false }
func (m *mockContext) GetBool(key string) (bool, bool)       { return false, false }
func (m *mockContext) GetInt(key string) (int, bool)         { return 0, false }
func (m *mockContext) GetString(key string) (string, bool)   { return "", false }
func (m *mockContext) TakeHeader() sylph.IHeader             { return nil }
func (m *mockContext) StoreHeader(header sylph.IHeader)      {}
func (m *mockContext) JwtClaim() sylph.IJwtClaim             { return nil }
func (m *mockContext) StoreJwtClaim(claim sylph.IJwtClaim)   {}
func (m *mockContext) TakeLogger() sylph.ILogger             { return nil }

// Logging methods
func (m *mockContext) Warn(pkg, action string, data any)             {}
func (m *mockContext) Info(pkg, action string, data any)             {}
func (m *mockContext) Error(pkg, action string, err error, data any) {}
func (m *mockContext) Debug(pkg, action string, data any)            {}
func (m *mockContext) Fatal(pkg, action string, data any)            {}
func (m *mockContext) Panic(pkg, action string, data any)            {}

// Control methods
func (m *mockContext) IsAbort() bool { return false }
func (m *mockContext) SetAbort()     {}
func (m *mockContext) Abort()        {}

// context.Context interface methods
func (m *mockContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (m *mockContext) Done() <-chan struct{}                   { return nil }
func (m *mockContext) Err() error                              { return nil }
func (m *mockContext) Value(key interface{}) interface{}       { return nil }

func newMockContext() sylph.Context {
	return &mockContext{}
}
