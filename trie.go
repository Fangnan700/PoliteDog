package PoliteDog

import "strings"

// Trie Router的本质是一棵前缀树，路径存储路由，末尾节点存储该路由对应的method、handler等信息
type Trie struct {
	next *TrieNode
}

type TrieNode struct {
	part     string
	path     string
	method   string
	children []*TrieNode
	key      string
	end      bool
}

// Insert 插入节点
func (tn *TrieNode) Insert(method string, pattern string, key string) {
	root := tn
	parts := strings.Split(pattern, "/")

	for i, part := range parts {
		if i == 0 {
			continue
		}

		matched := false
		for _, child := range tn.children {
			if part == child.part {
				tn = child
				matched = true
				break
			}
		}

		if !matched {
			child := &TrieNode{
				part:     part,
				children: make([]*TrieNode, 0),
			}
			tn.children = append(tn.children, child)
			tn = child
		}

		if i >= len(parts)-1 {
			tn.key = key
			tn.end = true
			tn.path = pattern
			tn.method = method
		}
	}

	tn = root
}

// Search 搜索路由
func (tn *TrieNode) Search(path string) *TrieNode {
	root := tn
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if i == 0 {
			continue
		}

		for _, child := range tn.children {
			if child.part == part || strings.Contains(child.part, ":") || strings.Contains(child.part, "*") {
				tn = child
				if tn.end {
					return child
				}
			}
		}
	}

	tn = root
	return nil
}
