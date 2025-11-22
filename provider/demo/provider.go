package demo

import (
	"fmt"

	"github.com/RZXBxie/web_server/framework"
)

type DemoServiceProvider struct{}

// Name 将服务对应的字符串凭证返回
func (sp *DemoServiceProvider) Name() string {
	return Key
}

// Register 方法是注册初始化服务实例的方法，我们这里先暂定为NewDemoService
func (sp *DemoServiceProvider) Register(container framework.Container) framework.NewInstance {
	return NewDemoService
}

// Boot 这里只打印日志信息
func (sp *DemoServiceProvider) Boot(c framework.Container) error {
	fmt.Println("demo service boot")
	return nil
}

// Params 这里只返回一个container参数
func (sp *DemoServiceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

// IsDefer 这里选择延迟实例化
func (sp *DemoServiceProvider) IsDefer() bool {
	return true
}
