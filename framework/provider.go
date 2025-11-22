package framework

// NewInstance 定义了如何创建一个新实例，这是所有服务容器的创建方法
type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 定义一个服务提供者需要实现的接口
type ServiceProvider interface {
	// Register 在服务容器中注册了一个实例化服务的方法，是否在注册时就实例化这个服务，需要由IsDefer()方法来决定
	// 参数Container是服务容器
	// 返回值 NewInstance 是实例化服务的方法，会在服务需要使用时调用这个方法
	Register(Container) NewInstance

	// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次make的时候进行实例化操作
	// false表示不需要延迟实例化，在注册的时候就实例化。true表示延迟实例化
	IsDefer() bool

	// Params 定义传递给NewInstance方法的参数，可以自定义多个，但建议将Container作为第一个参数
	Params(Container) []interface{}

	// Name 定义了这个服务提供者的名称
	Name() string

	// Boot 在调用服务实例化的时候会调用这个方法，可以把一些准备工作：基础配置、初始化参数的操作放在这里
	// 如果 Boot 返回error，整个服务实例化就会实例化失败
	Boot(Container) error
}
