
在 Go 语言中，当我们需要初始化一个复杂结构体时，通常会使用构造函数（如 `New` 方法）来封装初始化逻辑。传统的构造函数传入固定参数，但随着结构体属性的增多或变化，这种方式可能会导致构造函数的参数列表过长、调用不灵活，甚至带来不必要的复杂性。

为了提高代码的灵活性和可读性，可以使用 **Option 模式**。Option 模式通过引入可选参数，将初始化逻辑模块化，使代码更易于扩展和维护。这种设计模式在许多知名 Go 库中（如 `gorm`、`go-redis`）被广泛使用。

``` go
package main

import (
	"fmt"
	"time"
)

// ConfigService 代表一个配置服务
type ConfigService struct {
	name        string        // 服务名称（必填字段）
	timeout     time.Duration // 超时时间（可选字段，默认值为 30 秒）
	enableCache bool          // 是否启用缓存（可选字段，默认值为 false）
	logger      interface{}   // 日志记录器（可选字段，默认值为 nil）
}

// ConfigOption 定义可选参数的函数类型
type ConfigOption func(*ConfigService)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(cs *ConfigService) {
		cs.timeout = timeout
	}
}

// WithCache 设置是否启用缓存
func WithCache(enable bool) ConfigOption {
	return func(cs *ConfigService) {
		cs.enableCache = enable
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger interface{}) ConfigOption {
	return func(cs *ConfigService) {
		cs.logger = logger
	}
}

// NewConfigService 初始化 ConfigService，支持必填和可选参数
func NewConfigService(name string, options ...ConfigOption) *ConfigService {
	// 初始化默认值
	cs := &ConfigService{
		name:        name,              // 必填字段
		timeout:     30 * time.Second,  // 默认超时时间
		enableCache: false,             // 默认不启用缓存
		logger:      nil,               // 默认无日志记录器
	}

	// 应用所有可选参数
	for _, option := range options {
		option(cs)
	}

	return cs
}

// 使用示例
func main() {
	// 只设置必填字段
	config1 := NewConfigService("Service1")
	fmt.Printf("Config1: %+v\n", config1)

	// 设置必填字段和部分可选字段
	config2 := NewConfigService("Service2", WithTimeout(10*time.Second), WithCache(true))
	fmt.Printf("Config2: %+v\n", config2)

	// 设置所有字段
	config3 := NewConfigService("Service3", WithTimeout(1*time.Minute), WithCache(true), WithLogger("CustomLogger"))
	fmt.Printf("Config3: %+v\n", config3)
}


```
说明：
    
`name` 是服务的必填参数，作为第一个参数直接传入 `NewConfigService`。

其他字段（如 `timeout`, `enableCache`, `logger`）被设计为可选参数，通过 `ConfigOption` 实现。

在初始化时（`NewConfigService`），为所有可选字段设置默认值，如超时时间默认为 30 秒，缓存默认关闭。

使用 `WithTimeout`, `WithCache`, `WithLogger` 函数为可选参数提供灵活的赋值方式，不需要所有调用方传递所有字段，简化了调用。

如果后续 `ConfigService` 增加了新的字段（如 `retryPolicy`），只需新增一个对应的 `WithRetryPolicy` 函数即可，调用方无需改动。

调用代码如 `NewConfigService("Service2", WithTimeout(10*time.Second), WithCache(true))` 清晰表达了 "服务名为 `Service2`，超时时间为 10 秒，启用了缓存"。

    不需要频繁修改 `NewConfigService` 的参数列表，新增功能只需添加新的 `WithXXX` 方法即可，降低代码耦合度。