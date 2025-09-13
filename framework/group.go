package framework

// IGroup 路由组接口，新建路由组的时候返回IGroup类型，便于后续扩展。
type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
}

type Group struct {
	Core   *Core
	Prefix string
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		Core:   core,
		Prefix: prefix,
	}
}

func (group *Group) Get(uri string, handler ControllerHandler) {
	uri = group.Prefix + uri
	group.Core.Get(uri, handler)
}

func (group *Group) Post(uri string, handler ControllerHandler) {
	uri = group.Prefix + uri
	group.Core.Post(uri, handler)
}

func (group *Group) Put(uri string, handler ControllerHandler) {
	uri = group.Prefix + uri
	group.Core.Put(uri, handler)
}

func (group *Group) Delete(uri string, handler ControllerHandler) {
	uri = group.Prefix + uri
	group.Core.Delete(uri, handler)
}
