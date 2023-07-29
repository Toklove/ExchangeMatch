package engine

import (
	"encoding/json"
	"fmt"
)

type NodeLink struct {
	Current *OrderNode
	Node    *OrderNode
}

func (nl *NodeLink) InitOrderLink() {
	nl.Node.IsFirst = true
	nl.Node.IsLast = true

	nl.SetFirstPointer(nl.Node.NodeName)
	nl.SetLastPointer(nl.Node.NodeName)
	nl.SetLinkNode(nl.Node, nl.Node.NodeName)
}

func (nl *NodeLink) GetLinkNode(nodeName string) *OrderNode {
	res := cache.HGet(ctx, nl.Node.NodeLink, nodeName)
	if res.Val() == "" {
		return &OrderNode{}
	}

	node := &OrderNode{}
	err := json.Unmarshal([]byte(res.Val()), &node)
	if err != nil {
		fmt.Println(err)
	}
	nl.Current = node // 获取某个节点时，把此节点设置为当前节点

	return node
}

func (nl *NodeLink) SetFirstPointer(nodeName string) {
	cache.HSet(ctx, nl.Node.NodeLink, "f", nodeName)
}

func (nl *NodeLink) GetFirstNode() *OrderNode {
	res := cache.HGet(ctx, nl.Node.NodeLink, "f")
	if res.Val() == "" {
		return &OrderNode{}
	}

	node := nl.GetLinkNode(res.Val())
	if node.Uuid != "" {
		return node
	}
	nl.Current = node

	return node
}

func (nl *NodeLink) SetLast() {
	nl.GetLast()
	nl.Current.IsLast = false
	nl.Current.NextNode = nl.Node.NodeName
	nl.SetLinkNode(nl.Current, nl.Current.NodeName) // 更新本身节点信息

	nl.Node.PrevNode = nl.Current.NodeName // 重置last节点信息
	nl.SetLastPointer(nl.Node.NodeName)

	nl.Node.IsLast = true
	nl.SetLinkNode(nl.Node, nl.Node.NodeName)
}

func (nl *NodeLink) SetLastPointer(nodeName string) {
	cache.HSet(ctx, nl.Node.NodeLink, "l", nodeName)
}

func (nl *NodeLink) GetLast() *OrderNode {
	res := cache.HGet(ctx, nl.Node.NodeLink, "l")
	if res.Val() == "" {
		return &OrderNode{}
	}

	node := nl.GetLinkNode(res.Val())
	if node.Uuid == "" {
		return node
	}
	nl.Current = node

	return node
}

func (nl *NodeLink) GetCurrent() *OrderNode {
	return nl.Current
}

func (nl *NodeLink) GetPrev() *OrderNode {
	current := nl.GetCurrent()
	prevName := current.PrevNode
	if prevName == "" {
		return &OrderNode{}
	}
	node := nl.GetLinkNode(prevName)
	if node.Oid == "" {
		return &OrderNode{}
	}
	//nl.Current = node //是否需要重置当前节点?

	return node
}

func (nl *NodeLink) GetNext() *OrderNode {
	current := nl.GetCurrent()
	nextName := current.NextNode
	if nextName == "" {
		return &OrderNode{}
	}
	node := nl.GetLinkNode(nextName)
	if node.Oid == "" {
		return &OrderNode{}
	}
	//nl.Current = node //是否需要重置当前节点?

	return node
}

func (nl *NodeLink) SetLinkNode(node *OrderNode, nodeName string) {
	nodeJson, _ := json.Marshal(node)
	cache.HSet(ctx, nl.Node.NodeLink, nodeName, nodeJson)
}

func (nl *NodeLink) DeleteLinkNode(node *OrderNode) {
	if node.IsFirst && node.IsLast {
		cache.HDel(ctx, node.NodeLink, "f")
		cache.HDel(ctx, node.NodeLink, "l")
		cache.HDel(ctx, node.NodeLink, node.NodeName)
	} else if node.IsFirst && !node.IsLast {
		next := nl.GetNext()
		if next.Oid == "" {
			panic("expects next node is not empty.")
		}
		cache.HDel(ctx, node.NodeLink, node.NodeName)
		next.IsFirst = true
		next.PrevNode = ""
		nl.SetFirstPointer(next.NodeName)
		nl.SetLinkNode(next, next.NodeName)
	} else if !node.IsFirst && node.IsLast {
		prev := nl.GetPrev()
		if prev.Oid == "" {
			panic("expects prev node is not empty.")
		}
		cache.HDel(ctx, node.NodeLink, node.NodeName)
		prev.IsLast = true
		prev.NextNode = ""
		nl.SetLastPointer(prev.NodeName)
		nl.SetLinkNode(prev, prev.NodeName)
	} else {
		prev := nl.GetPrev()
		current := nl.GetNext()
		next := nl.GetNext()

		if prev.Oid == "" && next.Oid == "" {
			panic("expects relation node is not empty.")
		}
		cache.HDel(ctx, current.NodeLink, current.NodeName)

		prev.NextNode = next.NodeName
		next.PrevNode = prev.NodeName
		nl.SetLinkNode(prev, prev.NodeName)
		nl.SetLinkNode(next, next.NodeName)
	}
}
