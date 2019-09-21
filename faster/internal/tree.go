package internal

// Tree -- assigns sequential indexes to a key hierarchy
type Tree struct {
	curIndex int
	root     *TreeNode
}

type TreeNode struct {
	Index    int
	Children map[string]*TreeNode
}

func (t *Tree) nextIndex() int {
	var rc = t.curIndex
	t.curIndex++
	return rc
}

// Clone -- create a deep copy of this Tree object
//
// TODO there is potential for optimizations here (e.g. `make([]treeNode, t.curIndex+1` and using that )
func (t *Tree) Clone() *TreeNode {
	var rc *TreeNode
	if t.root != nil {
		rc = t.root.cloneRec()
	}
	return rc
}

// Exists -- returns true if the path already exists
func (t *Tree) Exists(path ...string) bool {
	// maybe: if len(path) == 0 { return true }
	if t.root == nil {
		return false
	}

	return t.root.Exists(path...)
}

// GetIndex -- returns the sequential index assigned to the given path
// (will create nodes recursively)
func (t *Tree) GetIndex(path ...string) int {
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
func (t *Tree) Reset() {
	t.root = nil
	t.curIndex = 0
}

func (n *TreeNode) cloneRec() *TreeNode {
	var rc = TreeNode{
		Index:    n.Index,
		Children: make(map[string]*TreeNode, len(n.Children)),
	}

	for key, child := range n.Children {
		rc.Children[key] = child.cloneRec()
	}
	return &rc
}

// Exists -- returns true if the given tree node exists (synonymous to GetIndex(path) >= 0)
func (n *TreeNode) Exists(path ...string) bool {
	return n.GetIndex(path...) >= 0
}

// GetIndex -- returns the index of the given path (if found, -1 otherwise)
func (n *TreeNode) GetIndex(path ...string) int {
	if len(path) == 0 {
		return n.Index
	}
	if child, ok := n.Children[path[0]]; ok {
		return child.GetIndex(path[1:]...)
	}
	return -1
}

func (n *TreeNode) Keys() []string {
	var rc = make([]string, 0, len(n.Children))
	for k := range n.Children {
		rc = append(rc, k)
	}
	return rc
}
