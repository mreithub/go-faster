package dashboard

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/mreithub/go-faster/faster"
)

type flatEntry struct {
	Name string
	Path []string
	Data faster.DataPoint
}

// Key -- returns Path + Name
func (e *flatEntry) Key() []string {
	return append(e.Path, e.Name)
}

func (e *flatEntry) JSONPath() string {
	var rc, _ = json.Marshal(append(e.Path, e.Name))
	return string(rc)
}

// formats a time.Duration as string (in msec) that can be easily parsed by the human eye when aligned right
func (e *flatEntry) toMsec(value time.Duration) string {
	if value == 0 {
		return ""
	}

	var mantissa = (value / time.Microsecond / 10) % 100 // two digits after the decimal point

	var unformatted = []rune(strconv.FormatInt(int64(value/time.Millisecond), 10))
	var formatted = make([]rune, 0, 3+int(float32(len(unformatted))*1.3))
	formatted = append(formatted, unformatted[:len(unformatted)%3]...)
	for i := len(unformatted) % 3; i < len(unformatted); i += 3 {
		if len(formatted) > 0 {
			formatted = append(formatted, ' ')
		}
		formatted = append(formatted, unformatted[i:i+3]...)
	}

	return fmt.Sprintf("%s.%02d", string(formatted), mantissa)
}

// PrettyAverage -- returns the average in msec (with space as thousands-separator)
func (e *flatEntry) PrettyAverage() string {
	return e.toMsec(e.Data.Average())
}

func (e *flatEntry) PrettyTotal() string {
	return e.toMsec(e.Data.TotalTime())
}

// flattenSnapshot -- takes the hierarchical data stored in a faster.Snapshot and puts it into a (sorted) slice
func flattenSnapshot(snap *faster.Snapshot) []flatEntry {
	return recFlattenSnapshot(nil, snap, nil)
}

func recFlattenSnapshot(rc []flatEntry, snap *faster.Snapshot, pathPrefix []string) []flatEntry {
	for _, k := range snap.Children(pathPrefix...) {
		if d := snap.Get(append(pathPrefix, k)...); d != nil {
			rc = append(rc, flatEntry{
				Name: k,
				Path: pathPrefix,
				Data: d,
			})
		}
	}

	for _, name := range snap.Children(pathPrefix...) {
		rc = recFlattenSnapshot(rc, snap, append(pathPrefix, name))
	}

	return rc
}

func sortByPath(data []flatEntry) {
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
