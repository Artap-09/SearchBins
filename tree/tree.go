package tree

import (
	"searchbin/interfaces"
)

type Tree struct {
	root *node
	size uint64
}

func NewTree() *Tree {
	return &Tree{}
}

func (tree *Tree) Insert(value interfaces.Range) {
	if tree.root == nil {
		tree.root = newNode(value, nil)
		tree.size++
		return
	}

	node := tree.root.insert(value)
	if node != nil {
		root := node.balance()

		tree.root = root
		tree.size++
	}

}

func (tree *Tree) Find(bin uint64) interfaces.Range {
	return tree.root.find(bin)
}

func (tree *Tree) Len() uint64 {
	return tree.size
}

func (tree *Tree) Hight() int {
	return tree.root.getHeight()
}

func (tree *Tree) Root() *node {
	return tree.root
}

func (tree *Tree) String() string {
	return tree.root.prints()
}
