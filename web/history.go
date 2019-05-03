package web

import (
	"sort"

	"github.com/mreithub/go-faster/faster"
)

func diff(data []*faster.Snapshot) []*faster.Snapshot {
	if len(data) == 0 {
		return nil
	}
	var rc = make([]*faster.Snapshot, 0, len(data)-1)

	var last *faster.Snapshot
	for _, snap := range data {
		if last != nil {
			rc = append(rc, last.Sub(snap))
		}
		last = snap
	}

	return rc
}

// sorts the given History tickers by their interval - lowest first
func sortHistoryByInterval(data map[string]*faster.History) []*faster.History {
	if data == nil || len(data) == 0 {
		return nil
	}
	var rc = make([]*faster.History, 0, len(data))

	for _, h := range data {
		rc = append(rc, h)
	}

	sort.Slice(rc, func(i int, j int) bool {
		return rc[i].Interval() < rc[j].Interval()
	})
	return rc
}
