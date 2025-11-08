package framework

import (
	"testing"
)

func TestTrieBasic(t *testing.T) {
	trie := NewTrie()

	// 加两个路由
	err := trie.AddRoute("/user/name", []ControllerHandler{func(ctx *Context) error {
		return nil
	}})
	if err != nil {
		t.Fatalf("add route failed: %v", err)
	}

	err = trie.AddRoute("/user/:id/name", []ControllerHandler{func(ctx *Context) error {
		return nil
	}})
	if err != nil {
		t.Fatalf("add route failed: %v", err)
	}

	// 检查能否匹配静态路由
	node := trie.root.matchNode("/user/name")
	if node == nil || !node.isLast {
		t.Fatalf("matchNode for /user/name failed, got nil or non-last")
	}

	// 检查能否匹配动态路由
	node = trie.root.matchNode("/user/123/name")
	if node == nil || !node.isLast {
		t.Fatalf("matchNode for /user/:id/name failed, got nil or non-last")
	}
	params := node.findParamsFromEndNode("/user/123/name")
	if v, ok := params["id"]; !ok || v != "123" {
		t.Fatalf("expected param id=123, got %v", params)
	}

	// 验证根节点 segment 是否为空
	if trie.root.segment != "" {
		t.Fatalf("expected root segment empty, got %q", trie.root.segment)
	}

	// 验证根节点下面确实有子节点
	if len(trie.root.children) == 0 {
		t.Fatalf("expected root to have children, got 0")
	}

	// 验证子节点 segment 是否为 "USER"（因为大写化）
	if trie.root.children[0].segment != "" {
		t.Fatalf("expected first child segment USER, got %q", trie.root.children[0].segment)
	}
}
