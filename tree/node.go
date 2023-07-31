package tree

import (
	"math"
	"searchbin/interfaces"
)

type node struct {
	value  interfaces.Range
	left   *node
	right  *node
	parent *node
	height int
}

func newNode(value interfaces.Range, parent *node) *node {
	return &node{value: value, height: 1, parent: parent}
}

func (n *node) insert(value interfaces.Range) *node {
	defer n.updateHigh()

	if n.value.Equal(value) {
		return nil
	}

	if n.value.RangeBelow(value) {

		if n.left == nil {
			n.left = newNode(value, n)
			return n.left
		}

		return n.left.insert(value)
	}

	if n.value.RangeHigher(value) {

		if n.right == nil {
			n.right = newNode(value, n)
			return n.right
		}

		return n.right.insert(value)
	}

	return nil
}

func (n *node) find(bin uint64) interfaces.Range {
	if result := n.value.Contains(bin); result != nil {
		return result
	}

	if n.value.BinBelow(bin) {
		if n.left == nil {
			return nil
		}

		return n.left.find(bin)
	}

	if n.value.BinHigher(bin) {
		if n.right == nil {
			return nil
		}

		return n.right.find(bin)
	}

	return nil
}

func (n *node) getHeight() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *node) balance() *node {
	nextBalance := n
	if math.Abs(float64(n.left.getHeight())-float64(n.right.getHeight())) > 1 {
		if n.right.getHeight() > n.left.getHeight() {
			if n.right.right.getHeight() >= n.right.left.getHeight() {
				nextBalance = n.smallLeftRotation()
			} else {
				nextBalance = n.largeLeftRotation()
			}
		} else {
			if n.left.left.getHeight() >= n.left.right.getHeight() {
				nextBalance = n.smallRightRotation()
			} else {
				nextBalance = n.largeRightRotation()
			}
		}
	}

	if nextBalance.parent == nil {
		return nextBalance
	}

	return nextBalance.parent.balance()
}

func (n *node) updateHigh() {
	if n.left.getHeight() >= n.right.getHeight() {
		n.height = n.left.getHeight() + 1
	} else {
		n.height = n.right.getHeight() + 1
	}
}

func (n *node) smallLeftRotation() *node {
	a := n
	b := a.right
	c := b.left

	if a.parent != nil {
		if a.parent.right == a {
			a.parent.right = b
		} else {
			a.parent.left = b
		}
	}

	a.parent, b.parent = b, a.parent
	b.left, a.right = a, c
	if c != nil {
		c.parent = a
	}

	a.updateHigh()
	b.updateHigh()

	return b
}

func (n *node) largeLeftRotation() *node {
	a := n
	b := n.right
	c := n.right.left

	if a.parent != nil {
		if a.parent.right == a {
			a.parent.right = c
		} else {
			a.parent.left = c
		}
	}

	a.parent, c.parent = c, a.parent
	b.parent, c.right, b.left = c, b, c.right
	a.right, c.left = c.left, a

	if a.right != nil {
		a.right.parent = a
	}

	if b.left != nil {
		b.left.parent = b
	}

	a.updateHigh()
	b.updateHigh()
	c.updateHigh()

	return c
}

func (n *node) smallRightRotation() *node {
	a := n
	b := a.left
	c := b.right

	if a.parent != nil {
		if a.parent.right == a {
			a.parent.right = b
		} else {
			a.parent.left = b
		}
	}

	a.parent, b.parent = b, a.parent
	b.right, a.left = a, c
	if c != nil {
		c.parent = a
	}

	a.updateHigh()
	b.updateHigh()

	return b
}

func (n *node) largeRightRotation() *node {
	a := n
	b := n.left
	c := n.left.right

	if a.parent != nil {
		if a.parent.right == a {
			a.parent.right = c
		} else {
			a.parent.left = c
		}
	}

	a.parent, c.parent = c, a.parent
	b.parent, c.left, b.right = c, b, c.left
	a.left, c.right = c.right, a

	if a.left != nil {
		a.left.parent = a
	}

	if b.right != nil {
		b.right.parent = b
	}

	a.updateHigh()
	b.updateHigh()
	c.updateHigh()

	return c
}

func (n *node) prints() string {
	if n == nil {
		return ""
	}

	result := ""

	result += n.left.prints()
	result += n.value.String()
	result += n.right.prints()

	return result
}
