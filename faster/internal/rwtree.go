package internal

// RWTree -- read/write wrapper around the read-only TreeNode struct
type RWTree struct {
	curIndex int
	root     *Tree

	// Limit -- after curIndex reached this Limit, a catch-all root node named '_overflow'
	// will be created.
	//
	// any attempt to create new nodes will return the index of that overflow node
	//
	// set to <= 0 to disable
	Limit int
}

// nextIndex -- increments .curIndex, returning the old value
func (t *RWTree) nextIndex(ignoreLimit bool) int {
	var rc = t.curIndex
	if ignoreLimit || t.Limit <= 0 || rc < t.Limit {
		t.curIndex++
	}
	return rc
}

// Clone -- returns a read-only deep copy of the internal Tree structure
func (t *RWTree) Clone() *Tree {
	var rc *Tree
	if t.root != nil {
		rc = t.root.cloneRec()
	}
	return rc
}

// Exists -- returns true if the path already exists
func (t *RWTree) Exists(path ...string) bool {
	// maybe: if len(path) == 0 { return true }
	if t.root == nil {
		return false
	}

	return t.root.Exists(path...)
}

// GetIndex -- returns the sequential index assigned to the given path
// (will create new tree nodes recursively)
func (t *RWTree) GetIndex(path ...string) int {
	if t.root == nil {
		t.root = &Tree{
			index: t.nextIndex(false), // 0
		}
	}

	var curNode = t.root
	for len(path) > 0 {
		if curNode.children == nil {
			curNode.children = make(map[string]*Tree)
		}

		var child *Tree
		var ok bool
		if child, ok = curNode.children[path[0]]; !ok {
			var nextIndex = t.nextIndex(false)
			if t.Limit <= 0 || nextIndex < t.Limit {
				child = &Tree{
					index: nextIndex,
				}
				curNode.children[path[0]] = child
			} else {
				return t.getOverflowIndex()
			}
		}

		curNode = child
		path = path[1:]
	}
	return curNode.index
}

// getOverflowIndex -- returns the index of the root tree entry '_overflow' (creates the node if neccessary)
func (t *RWTree) getOverflowIndex() int {
	if t.root == nil {
		t.root = &Tree{
			index: t.nextIndex(true), // 0
		}
	}

	var node *Tree
	if node = t.root.children["_overflow"]; node == nil {
		node = &Tree{
			index: t.nextIndex(true),
		}
		t.root.children["_overflow"] = node
	}

	return node.index
}

// Reset -- removes all nodes and resets the sequential index to 0
func (t *RWTree) Reset() {
	t.root = nil
	t.curIndex = 0
}
