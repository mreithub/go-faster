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

// Exists -- returns true if the given tree node exists (synonymous to GetIndex(path) >= 0)
func (n *Tree) Exists(path ...string) bool {
	return n.GetIndex(path...) >= 0
}

// GetIndex -- returns the index of the given path (if found, -1 otherwise)
func (n *Tree) GetIndex(path ...string) int {
	if len(path) == 0 {
		return n.index
	}
	if child, ok := n.children[path[0]]; ok {
		return child.GetIndex(path[1:]...)
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
