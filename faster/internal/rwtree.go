package internal

// RWTree -- read/write wrapper around the read-only TreeNode struct
type RWTree struct {
	curIndex int
	root     *TreeNode
}

// nextIndex -- increments .curIndex, returning the old value
func (t *RWTree) nextIndex() int {
	var rc = t.curIndex
	t.curIndex++
	return rc
}

// Clone -- create a deep copy of this Tree object
//
// TODO there is potential for optimizations here (e.g. `make([]treeNode, t.curIndex+1` and using that )
func (t *RWTree) Clone() *TreeNode {
	var rc *TreeNode
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
		t.root = &TreeNode{
			Index: t.nextIndex(), // 0
		}
	}

	var curNode = t.root
	for len(path) > 0 {
		if curNode.Children == nil {
			curNode.Children = make(map[string]*TreeNode)
		}

		var child *TreeNode
		var ok bool
		if child, ok = curNode.Children[path[0]]; !ok {
			child = &TreeNode{
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
