package demo

import (
	"fmt"

	"github.com/RZXBxie/web_server/framework"
)

type DemoService struct {
	Service

	// c 参数
	c framework.Container
}

// NewDemoService 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	fmt.Println("new demo service")
	return &DemoService{c: c}, nil
}

func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}
