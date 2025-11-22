package demo

const Key = "demo"

// Service Demo服务接口
type Service interface {
	GetFoo() Foo
}

// Foo Demo服务接口定义的一个数据结构
type Foo struct {
	Name string
}
