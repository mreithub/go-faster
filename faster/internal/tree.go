package internal

// Tree -- read only struct that maps hierarchical paths to integer indexes
//
// indexes are positive numbers (with 0 being the index of the root node)
type Tree struct {
	index    int
	children map[string]*Tree
}

func (n *Tree) cloneRec() *Tree {
	var rc = Tree{
		index:    n.index,
		children: make(map[string]*Tree, len(n.children)),
	}

	for key, child := range n.children {
		rc.children[key] = child.cloneRec()
	}
	return &rc
}

// getNode -- recursively traverses the tree to find the node with the given path (returns nil if not found)
func (n *Tree) getNode(path ...string) *Tree {
	if len(path) == 0 {
		return n
	}

	if child, ok := n.children[path[0]]; ok {
		return child.getNode(path[1:]...)
	}
	return nil

}

// Exists -- returns true if the given tree node exists
func (n *Tree) Exists(path ...string) bool {
	return n.getNode(path...) != nil
}

// GetIndex -- returns the index of the given path (if found, -1 otherwise)
func (n *Tree) GetIndex(path ...string) int {
	var node = n.getNode(path...)
	if node != nil {
		return node.index
	}
	return -1
}

func (n *Tree) Keys() []string {
	var rc = make([]string, 0, len(n.children))
	for k := range n.children {
		rc = append(rc, k)
	}
	return rc
}
