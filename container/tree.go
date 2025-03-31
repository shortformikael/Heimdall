package container

import (
	"errors"
	"fmt"
)

type TreeGraph struct {
	Head *TreeNode
}

func (g *TreeGraph) Search(s any) (*TreeNode, error) {
	if g.Head == nil {
		return nil, errors.New("EMPTY TREE")
	} else if g.Head.Children == nil {
		if g.Head.Value == s {
			return g.Head, nil
		} else {
			return nil, errors.New("VALUE DOES NOT EXIST IN TREE")
		}
	}

	current := &g.Head.Children[0]
	result := g.dfsSearch(*current, &s)

	if result != nil {
		return result, nil
	} else {
		return nil, errors.New("VALUE DOES NOT EXIST IN TREE")
	}

}

func (g *TreeGraph) dfsSearch(current *TreeNode, s *any) *TreeNode {
	for _, child := range current.Children {
		if child.Value == s {
			return child
		} else if child.Children != nil {
			return g.dfsSearch(child, s)
		}
	}
	return nil
}

type TreeNode struct {
	Value    any
	Children []*TreeNode
}

func (g *TreeNode) DfsPrintTraversal() {
	fmt.Println("Node Value: ", g.Value)
	for i := 0; i < len(g.Children); i++ {
		g.Children[i].DfsPrintTraversal()
	}
}
func (n *TreeNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
}

func (n *TreeNode) Get(i int) (*TreeNode, error) {
	if i < 0 || len(n.Children) < i {
		return nil, errors.New("OUT OF BOUNDS")
	}
	return n.Children[i], nil
}

func (n *TreeNode) SearchChildren(child any) (*TreeNode, error) {
	for i := 0; i < len(n.Children); i++ {
		if child == n.Children[i].Value {
			return n.Children[i], nil
		}
	}
	return nil, errors.New("DOES NOT EXIST IN ARRAY")
}

func (n *TreeNode) String() string {
	switch s := n.Value.(type) {
	case string:
		return s
	case MenuItem:
		return s.Name
	default:
		return "NO STRING AVAILABLE"
	}
}
