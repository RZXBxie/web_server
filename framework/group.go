package framework

// IGroup 路由组接口，新建路由组的时候返回IGroup类型，便于后续扩展。
type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// Group 实现嵌套路由组的功能
	Group(prefix string) IGroup
	Use(middlewares ...ControllerHandler)
}

type Group struct {
	Core   *Core
	Prefix string
	Parent *Group

	Middlewares []ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		Core:   core,
		Prefix: prefix,
		Parent: nil,
	}
}

func (group *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	allHandlers := append(group.getMiddlewares(), handlers...)
	group.Core.Get(uri, allHandlers...)
}

func (group *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	allHandlers := append(group.getMiddlewares(), handlers...)
	group.Core.Post(uri, allHandlers...)
}

func (group *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	allHandlers := append(group.getMiddlewares(), handlers...)
	group.Core.Put(uri, allHandlers...)
}

func (group *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = group.getAbsolutePrefix() + uri
	allHandlers := append(group.getMiddlewares(), handlers...)
	group.Core.Delete(uri, allHandlers...)
}

func (group *Group) Group(uri string) IGroup {
	cGroup := NewGroup(group.Core, uri)
	cGroup.Parent = group
	return cGroup
}

func (group *Group) Use(middlewares ...ControllerHandler) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}

func (group *Group) getAbsolutePrefix() string {
	if group.Parent == nil {
		return group.Prefix
	}
	return group.Parent.getAbsolutePrefix() + group.Prefix
}

// getMiddlewares 获取某个group的中间件
// 这里就是获取除了Get/Post/Put/Delete之外设置的middleware
func (group *Group) getMiddlewares() []ControllerHandler {
	if group.Parent == nil {
		return group.Middlewares
	}
	return append(group.Parent.getMiddlewares(), group.Middlewares...)
}
