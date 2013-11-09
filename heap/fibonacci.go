package heap

import "math"

/*
  Ported to Go from http://code.google.com/p/graphmaker/source/browse/core/src/com/bluemarsh/graphmaker/core/util/FibonacciHeap.java

  The contents of this file are subject to the terms of the Common Development
  and Distribution License (the License). You may not use this file except in
  compliance with the License.

  You can obtain a copy of the License at http://www.netbeans.org/cddl.html
  or http://www.netbeans.org/cddl.txt.

  When distributing Covered Code, include this CDDL Header Notice in each file
  and include the License file at http://www.netbeans.org/cddl.txt.
  If applicable, add the following below the CDDL Header, with the fields
  enclosed by brackets [] replaced by your own identifying information:
  "Portions Copyrighted [year] [name of copyright owner]"

  The Original Software is GraphMaker. The Initial Developer of the Original
  Software is Nathan L. Fiedler. Portions created by Nathan L. Fiedler
  are Copyright (C) 1999-2008. All Rights Reserved.

  Contributor(s): Nathan L. Fiedler. Remco Schoeman

  This class implements a Fibonacci heap data structure.

  The code in this class is based on the algorithms in Chapter 21 of the
  "Introduction to Algorithms" by Cormen, Leiserson, Rivest, and Stein.
  The amortized running time of most of these methods is O(1), making
  it a very fast data structure. Several have an actual running time
  of O(1). removeMin() and delete() have O(log n) amortized running
  times because they do the heap consolidation.

  Note that this implementation is not synchronized.
  If multiple threads access a set concurrently, and at least one of the
  threads modifies the set, it must be synchronized externally.
  This is typically accomplished by synchronizing on some object that
  naturally encapsulates the set.

  author:  Nathan Fiedler
*/
type FibHeap struct {
	/* Points to the minimum node in the heap. */
	min *fibNode
	/* Points to the head of the free nodes in the heap. */
	free *fibNode
	/* Number of nodes in the heap. If the type is ever widened,
	(e.g. changed to long) then recalcuate the maximum degree
	value used in the consolidate() method. */
	length int
}

func NewFibonacci() *FibHeap {
	return &FibHeap{nil, nil, 0}
}

/*
  Removes all elements from this heap.

  Running time: O(1)
*/
func (h *FibHeap) Clear() {
	h.min = nil
	h.free = nil
	h.length = 0
}

/*
  Consolidates the trees in the heap by joining trees of equal
  degree until there are no more trees of equal degree in the
  root list

  Running time: O(log n) amortized
*/
func (h *FibHeap) consolidate() {
	// The magic 45 comes from log base phi of Integer.MAX_VALUE,
	// which is the most elements we will ever hold, and log base
	// phi represents the largest degree of any root list node.
	A := make([]*fibNode, 45)

	// For each root list node look for others of the same degree.
	start := h.min
	w := h.min
	for {
		x := w
		// Because x might be moved, save its sibling now.
		nextW := w.right
		d := x.degree
		for A[d] != nil {
			// Make one of the nodes a child of the other.
			y := A[d]
			if x.key > y.key {
				x, y = y, x
			}
			if y == start {
				// Because removeMin() arbitrarily assigned the min
				// reference, we have to ensure we do not miss the
				// end of the root node list.
				start = start.right
			}
			if y == nextW {
				// If we wrapped around we need to check for this case.
				nextW = nextW.right
			}
			// Node y disappears from root list.
			y.link(x)
			// We've handled this degree, go to next one.
			A[d] = nil
			d++
		}
		// Save this node for later when we might encounter another
		// of the same degree.
		A[d] = x
		// Move forward through list.
		w = nextW

		if w == start {
			// we're done
			break
		}
	}

	// The node considered to be min may have been changed above.
	h.min = start

	// Find the minimum key again.
	for _, a := range A {
		if a != nil && a.key < h.min.key {
			h.min = a
		}
	}
}

/*
  Deletes a node from the heap given the data handle.
  The trees in the heap will be consolidated, if necessary.

  Running time: O(log n)

  data:  data handle to remove from heap.
*/
func (h *FibHeap) Delete(data int) {
	n, ok := find(h.min, data)
	if ok {
		// make n the minimum node.
		h.decreaseKey(n, -math.MaxFloat64, true)
		// remove it
		h.DeleteMin()
	}
}

/*
  Finds the first node with the given data handle
*/
func find(n *fibNode, data int) (*fibNode, bool) {
	if n == nil {
		return nil, false
	}

	if n.data == data {
		return n, true
	}

	for sibling := n.right; sibling != n; sibling = sibling.right {
		if sibling.data == data {
			return sibling, true
		}
		if child, ok := find(sibling.child, data); ok {
			return child, ok
		}
	}
	return nil, false
}

/*
  Decreases the key value for a heap node, given the new value
  to take on. The structure of the heap may be changed, but will
  not be consolidated.

  Running time: O(log N) amortized

  x:  node to decrease the key of.

  k:  new key value for node x.

*/
func (h *FibHeap) DecreaseKey(data int, k float64) bool {
	x, ok := find(h.min, data)
	if ok {
		return h.decreaseKey(x, k, false)
	}
	return false
}

/*
  Decrease the key value of a node, or simply bubble it up to the
  top of the heap in preparation for a delete operation.

  Running time: O(1) amortized

  x:       node to decrease the key of.
  k:       new key value for node x.
  del:     true if deleting node (in which case, k is ignored).
*/
func (h *FibHeap) decreaseKey(x *fibNode, k float64, del bool) bool {
	if !del && k > x.key {
		return false
	}
	x.key = k
	y := x.parent
	if y != nil && (del || k < y.key) {
		y.cut(x, h.min)
		y.cascadingCut(h.min)
	}
	if del || k < h.min.key {
		h.min = x
	}
	return true
}

/*
  Tests if the Fibonacci heap is empty or not. Returns true if
  the heap is empty, false otherwise.

  Running time: O(1)

  Returns true if the heap is empty, false otherwise.
*/
func (h *FibHeap) IsEmpty() bool {
	return h.min == nil
}

/*
  Inserts a new data element into the heap. No heap consolidation
  is performed at this time, the new node is simply inserted into
  the root list of this heap.

  Running time: O(1)

  data:    data handle to insert into heap.

  key:    key value associated with data handle.

  Returns the newly created heap node.
*/
func (h *FibHeap) Insert(data int, key float64) {
	if data < 0 {
		panic("Negative data handles are not supported")
	}

	node := h.newNode(data, key)
	// concatenate node into min list
	if h.min != nil {
		node.right = h.min
		node.left = h.min.left
		h.min.left = node
		node.left.right = node
		if key < h.min.key {
			h.min = node
		}
	} else {
		h.min = node
	}
	h.length++
}

func (h *FibHeap) Contains(data int) bool {
	_, ok := find(h.min, data)
	return ok
}

/*
  Returns the smallest element in the heap. This smallest element
  is the one with the minimum key value.

  Running time: O(1)

  Returns the data handle with the smallest key and true, or 0 and false if empty.
*/
func (h *FibHeap) Min() (int, bool) {
	if h.min == nil {
		return 0, false
	}
	return h.min.data, true
}

/*
  Removes the smallest element from the heap. This will cause
  the trees in the heap to be consolidated, if necessary.

  Running time: O(log n) amortized

  Returns the  data handle with the smallest key and true, or 0 and false if the heap is empty.
*/
func (h *FibHeap) DeleteMin() (int, bool) {
	z := h.min
	if z == nil {
		return 0, false
	}
	data := z.data
	if z.child != nil {
		z.child.parent = nil
		// for each child of z do...
		for x := z.child.right; x != z.child; x = x.right {
			// set parent[x] to nil
			x.parent = nil
		}
		// merge the children into root list
		minleft := h.min.left
		zchildleft := z.child.left
		h.min.left = zchildleft
		zchildleft.right = h.min
		z.child.left = minleft
		minleft.right = z.child
	}
	// remove z from root list of heap
	z.left.right = z.right
	z.right.left = z.left
	if z == z.right {
		h.min = nil
	} else {
		h.min = z.right
		h.consolidate()
	}
	// decrement size of heap
	h.length--
	h.addFree(z)
	return data, true
}

/*
  Returns the size of the heap which is measured in the
  number of elements contained in the heap.

  Running time: O(1)

  Returns the number of elements in the heap.
*/
func (h *FibHeap) Len() int {
	return h.length
}

/*
  Moves the items of h2 to this one. No heap consolidation is
  performed at this time. The two root lists are simply joined together.
  h2 wil be empty after the Join

  Running time: O(1)

  h2  the heap to join
*/
func (h1 *FibHeap) Join(h2 *FibHeap) {
	h := FibHeap{} // temp storage
	if h1 != nil && h2 != nil {
		h.min = h1.min
		if h.min != nil {
			if h2.min != nil {
				h.min.right.left = h2.min.left
				h2.min.left.right = h.min.right
				h.min.right = h2.min
				h2.min.left = h.min
				if h2.min.key < h1.min.key {
					h.min = h2.min
				}
			}
		} else {
			h.min = h2.min
		}
		h.length = h1.length + h2.length
	}
	h2.Clear()
	h1.min = h.min
	h1.length = h.length
}

/*
  Implements a node of the Fibonacci heap. It holds the information
  necessary for maintaining the structure of the heap. It acts as
  an opaque handle for the data element, and serves as the key to
  retrieving the data from the heap.
*/
type fibNode struct {
	/* Key value for this node. */
	key float64
	/* Parent node. */
	parent *fibNode
	/* First child node. */
	child *fibNode
	/* Right sibling node. */
	right *fibNode
	/* Left sibling node. */
	left *fibNode
	/* handle to the data */
	data int
	/* Number of children of this node. */
	degree int
	/* True if this node has had a child removed since this node was
	   added to its parent. */
	mark bool
}

/*
  Gets a free node or creates a new node for the fibonacci heap,

  data:  data handle to associate with the new node

  key:   key value for the data handle
*/
func (h *FibHeap) newNode(data int, key float64) *fibNode {
	n := h.free
	if n == nil {
		n = &fibNode{}
	} else {
		h.free = n.child
		n.child = nil
	}
	n.data = data
	n.key = key
	return n
}

/*
  Adds a node to the free list for the fibonacci heap,
*/
func (h *FibHeap) addFree(n *fibNode) {
	if n == nil {
		return
	}
	n.reset(0, 0)
	n.child = h.free
	h.free = n
}

// resets the FibNode to it's initial state
func (n *fibNode) reset(data int, key float64) {
	n.data = data
	n.key = key
	n.left = n
	n.right = n
	n.degree = 0
	n.mark = false
	n.parent = nil
	n.child = nil
}

/*
  Performs a cascading cut operation. Cuts this from its parent
  and then does the same for its parent, and so on up the tree.

  Running time: O(log n)

  min:  the minimum heap node, to which nodes will be added.
*/
func (n *fibNode) cascadingCut(min *fibNode) {
	p := n.parent
	// if there's a parent...
	if p != nil {
		if n.mark {
			// it's marked, cut it from parent
			p.cut(n, min)
			// cut its parent as well
			p.cascadingCut(min)
		} else {
			// if y is unmarked, set it marked
			n.mark = true
		}
	}
}

/*
  The reverse of the link operation: removes x from the child
  list of this node.

  Running time: O(1)

  x:    child to be removed from this node's child list

  min:  the minimum heap node, to which x is added.
*/
func (n *fibNode) cut(x, min *fibNode) {
	// remove x from childlist and decrement degree
	x.left.right = x.right
	x.right.left = x.left
	n.degree--
	// reset child if necessary
	if n.degree == 0 {
		n.child = nil
	} else if n.child == x {
		n.child = x.right
	}
	// add x to root list of heap
	x.right = min
	x.left = min.left
	min.left = x
	x.left.right = x
	// set parent[x] to nil
	x.parent = nil
	// set mark[x] to false
	x.mark = false
}

/*
  Make this node a child of the given parent node. All linkages
  are updated, the degree of the parent is incremented, and
  mark is set to false.

  parent:  the new parent node.
*/
func (n *fibNode) link(parent *fibNode) {
	// Note: putting this code here in Node makes it faster
	// because it doesn't have to use generated accessor methods,
	// which add a lot of time when called millions of times.

	// remove n from its circular list
	n.left.right = n.right
	n.right.left = n.left
	// make this a child of x
	n.parent = parent
	if parent.child == nil {
		parent.child = n
		n.right = n
		n.left = n
	} else {
		n.left = parent.child
		n.right = parent.child.right
		parent.child.right = n
		n.right.left = n
	}
	// increase degree[x]
	parent.degree++
	// set mark false
	n.mark = false
}
