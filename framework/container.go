package framework

import (
	"errors"
	"sync"
)

// Container 服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，不返回error
	Bind(provider ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务的实例
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务的实例，如果实例不存在，则会panic
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务的实例，只是这个服务的实例是新的实例，不是单例
	// 它是根据服务提供者注册时的启动函数和传递的 params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params ...interface{}) (interface{}, error)
}

// MyContainer 是服务容器的具体实现
type MyContainer struct {
	// Container 表示MyContainer必须显示实现Container接口
	Container

	// providers 存储注册的服务提供者，key为服务名，value为服务提供者
	providers map[string]ServiceProvider

	// instances 存储实例化后的服务实例，key为服务名，value为服务实例
	instances map[string]interface{}

	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

func NewContainer() *MyContainer {
	return &MyContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
		lock:      sync.RWMutex{},
	}
}

func (c *MyContainer) Bind(provider ServiceProvider) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	key := provider.Name()
	c.providers[key] = provider

	if !provider.IsDefer() {
		if err := provider.Boot(c); err != nil {
			return err
		}
		params := provider.Params(c)
		method := provider.Register(c)
		instance, err := method(params...)
		if err != nil {
			return err
		}
		c.instances[key] = instance
	}

	return nil
}

func (c *MyContainer) IsBind(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.providers[key]
	return ok
}

// findServiceProvider 查找是否已经注册了这个服务提供者，如果没有注册，则返回nil
func (c *MyContainer) findServiceProvider(key string) ServiceProvider {
	if sp, ok := c.providers[key]; ok {
		return sp
	}
	return nil
}

func (c *MyContainer) Make(key string) (interface{}, error) {
	return c.make(key, nil, false)
}

func (c *MyContainer) MustMake(key string) interface{} {
	service, err := c.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return service
}

func (c *MyContainer) MakeNew(key string, params ...interface{}) (interface{}, error) {
	return c.make(key, params, false)
}

// make 真正的实例化一个服务
func (c *MyContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 查询是否已经注册了这个服务提供者，如果没有注册，则报错
	sp, ok := c.providers[key]
	if !ok {
		return nil, errors.New("contract " + key + " not found.")
	}
	if forceNew {
		return c.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := c.instances[key]; ok {
		return ins, nil
	}

	// 如果容器中不存在实例化服务，则进行实例化
	instance, err := c.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	c.instances[key] = instance
	return instance, nil
}

// newInstance 实例化一个服务
func (c *MyContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(c); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(c)
	}

	method := sp.Register(c)
	instance, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return instance, nil
}
