package framework

import (
	"errors"
	"strings"
)

type Trie struct {
	root *node
}

type node struct {
	// isLast 代表这个节点是否可以成为最终的路由规则。该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	isLast bool

	// segment uri中的字符串，代表这个节点表示的是路由中某个段的字符串
	segment string

	// handler 代表这个节点对应的handler，便于后续加载调用
	handler ControllerHandler

	// children 子节点
	children []*node
}

func newNode() *node {
	return &node{
		isLast:   false,
		segment:  "",
		handler:  nil,
		children: []*node{},
	}
}

func NewTrie() *Trie {
	root := newNode()
	return &Trie{root: root}
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) filterChildNode(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}

	// 如果segment是通配符，则所有子节点都满足要求
	if isWildSegment(segment) {
		return n.children
	}

	childNodes := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		// 如果子节点有通配符，则满足要求
		if isWildSegment(child.segment) {
			childNodes = append(childNodes, child)
		} else if child.segment == segment {
			// 如果没有通配符，但是完全匹配，则也符合要求
			childNodes = append(childNodes, child)
		}
	}

	return childNodes
}

// matchNode 从root开始查找，找到了就返回isLast=true的节点的指针，找不到就返回nil
// 树结构：以 /user/name 和 /user/:id/name 为例
//
//		root
//		└── "", isLast=false
//		    └── user, isLast=false
//	         ├── name, isLast=true
//		        └── :id, isLast=false
//	             └── name, isLast=true
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]

	// 将路径全部转为大写格式，让用户使用时大小写不敏感
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	childNodes := n.filterChildNode(segment)
	// 没找到，则说明该路由一定不存在
	if childNodes == nil || len(childNodes) == 0 {
		return nil
	}

	// 如果只能分成1个段，则判断子节点中是否有结尾的
	if len(segments) == 1 {
		for _, child := range childNodes {
			if child.isLast {
				return child
			}
		}
		return nil

	}

	// 如果有2个segment，递归调用matchNode
	for _, child := range childNodes {
		childMatch := child.matchNode(segments[1])
		if childMatch != nil {
			return childMatch
		}
	}

	return nil
}

// AddRoute 添加路由
func (trie *Trie) AddRoute(uri string, handler ControllerHandler) error {
	root := trie.root
	// 去掉开头和结尾多余的"/" /user/name -> user/name
	uri = strings.Trim(uri, "/")
	// 如果路由存在，则返回错误
	if root.matchNode(uri) != nil {
		return errors.New("Route exist: " + uri)
	}

	// 分段添加路由
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		// 非通配符路由则转为大写
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1

		var objNode *node

		// 查找是否存在现存的子节点
		childNodes := root.filterChildNode(segment)
		// 如果存在已有的子节点
		if len(childNodes) > 0 {
			for _, child := range childNodes {
				if child.segment == segment {
					objNode = child
					break
				}
			}
		}
		// 没找到已有的子节点，则新建一个
		if objNode == nil {
			node := newNode()
			node.segment = segment
			if isLast {
				node.isLast = true
				node.handler = handler
			}
			root.children = append(root.children, node)
			objNode = node
		}
		root = objNode
	}

	return nil
}

// FindHandler 根据路由查找handler
func (trie *Trie) FindHandler(uri string) ControllerHandler {
	uri = strings.Trim(uri, "/")
	node := trie.root.matchNode(uri)
	if node == nil {
		return nil
	}
	return node.handler
}
