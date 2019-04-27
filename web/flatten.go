package web

import (
	"sort"

	"github.com/mreithub/go-faster/faster"
)

type entry struct {
	Name string
	Path []string
	Data *faster.Snapshot
}

func flattenSnapshot(snap *faster.Snapshot) []entry {
	return recFlattenSnapshot(nil, snap, nil)
}

func recFlattenSnapshot(list []entry, snap *faster.Snapshot, pathPrefix []string) []entry {
	for k, v := range snap.Children {
		list = append(list, entry{
			Name: k,
			Path: pathPrefix,
			Data: v,
		})
	}

	for name, child := range snap.Children {
		list = recFlattenSnapshot(list, child, append(pathPrefix, name))
	}

	return list
}

func sortByPath(data []entry) {
	sort.Slice(data, func(i, j int) bool {
		return pathLessThan(
			append(data[i].Path, data[i].Name),
			append(data[j].Path, data[j].Name))
	})
}

func pathLessThan(a, b []string) bool {
	if len(b) == 0 {
		return false
	} else if len(a) == 0 {
		return true
	}

	var headA, tailA, headB, tailB = a[0], a[1:], b[0], b[1:]

	if headA < headB {
		return true
	} else if headA > headB {
		return false
	}

	return pathLessThan(tailA, tailB)
}
