package web

import (
	"sort"
	"strings"

	"github.com/mreithub/go-faster/faster"
)

type entry struct {
	Name string
	Path []string
	Data faster.Data
}

func flattenSnapshot(snap faster.Snapshot) []entry {
	return recFlattenSnapshot(nil, snap, nil)
}

func recFlattenSnapshot(list []entry, snap faster.Snapshot, pathPrefix []string) []entry {
	for k, v := range snap.Data {
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
		return pathLessThan(data[i].Path, data[j].Path)
	})
}

func pathLessThan(a, b []string) bool {
	// TODO avoid joining here (for added performance)
	return strings.Join(a, ".") < strings.Join(b, ".")
}
