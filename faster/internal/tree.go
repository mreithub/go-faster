package internal

// Tree -- read only struct that maps hierarchical paths to integer indexes
//
// indexes are positive numbers (with 0 being the index of the root node)
type Tree struct {
	Index    int
	Children map[string]*Tree
}

func (n *Tree) cloneRec() *Tree {
	var rc = Tree{
		Index:    n.Index,
		Children: make(map[string]*Tree, len(n.Children)),
	}

	for key, child := range n.Children {
		rc.Children[key] = child.cloneRec()
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
		return n.Index
	}
	if child, ok := n.Children[path[0]]; ok {
		return child.GetIndex(path[1:]...)
	}
	return -1
}

func (n *Tree) Keys() []string {
	var rc = make([]string, 0, len(n.Children))
	for k := range n.Children {
		rc = append(rc, k)
	}
	return rc
}
