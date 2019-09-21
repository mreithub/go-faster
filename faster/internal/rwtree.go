package internal

// RWTree -- read/write wrapper around the read-only TreeNode struct
type RWTree struct {
	curIndex int
	root     *Tree
}

// nextIndex -- increments .curIndex, returning the old value
func (t *RWTree) nextIndex() int {
	var rc = t.curIndex
	t.curIndex++
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
			Index: t.nextIndex(), // 0
		}
	}

	var curNode = t.root
	for len(path) > 0 {
		if curNode.Children == nil {
			curNode.Children = make(map[string]*Tree)
		}

		var child *Tree
		var ok bool
		if child, ok = curNode.Children[path[0]]; !ok {
			child = &Tree{
				Index: t.nextIndex(),
			}
			curNode.Children[path[0]] = child
		}

		curNode = child
		path = path[1:]
	}
	return curNode.Index
}

// Reset -- removes all nodes and resets the sequential index to 0
func (t *RWTree) Reset() {
	t.root = nil
	t.curIndex = 0
}
