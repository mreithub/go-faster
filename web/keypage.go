package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/mreithub/go-faster/faster"
)

// KeyPage -- implements /key/*
type KeyPage struct {
	faster    *faster.Faster
	templates map[string]*template.Template
}

// InfoJSON -- implements GET key/info.json
func (p *KeyPage) InfoJSON(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}

	var key = r.URL.Query()["k"]
	if len(key) == 0 {
		http.Error(w, "missing key parameter(s): 'k'", http.StatusBadRequest)
		return
	}

	ref := p.faster.Track("_faster", "key", "info.json")
	defer ref.Done()

	var tickers = p.faster.ListTickers()
	var sortedTickers = p.sortHistoryByInterval(tickers)
	var selectedTicker *faster.History

	type RequestInfo struct {
		TS      []int64 `json:"ts"`
		Counts  []int64 `json:"counts"`
		AvgMsec []int64 `json:"avgMsec"`
	}
	type Response struct {
		Requests  RequestInfo              `json:"requests"`
		Tickers   []map[string]interface{} `json:"tickers"`
		Histogram []map[string]interface{} `json:"histogram,omitempty"`

		Active int32 `json:"active"`
		Total  int64 `json:"total"`
		AvgMS  int64 `json:"avgMS"`
	}
	var info Response

	if len(sortedTickers) > 0 {
		var req = &info.Requests
		selectedTicker = p.getTicker(r, tickers, sortedTickers[0])
		var timeseries = selectedTicker.GetData(key...).Relative()
		for i, snap := range timeseries.Data {
			req.TS = append(req.TS, timeseries.GetTimestamp(i).UnixNano()/int64(time.Millisecond))
			req.Counts = append(req.Counts, snap.Count())
			req.AvgMsec = append(req.AvgMsec, int64(snap.Average()/time.Millisecond))
		}

		for _, h := range sortedTickers {
			info.Tickers = append(info.Tickers, map[string]interface{}{
				"name":       h.Name,
				"range":      h.Duration().String(),
				"intervalNS": h.Interval(),
				"interval":   h.Interval().String(),
				"capacity":   h.Capacity,
			})
		}
	}

	if snap := p.faster.TakeSnapshot().Get(key...); snap != nil {
		info.Active = snap.Active()
		info.AvgMS = int64(snap.Average() / time.Millisecond)
		info.Total = snap.Count()
		/*		if h := snap.Histogram; h != nil {
				var durations, counts = h.GetValues()
				if len(durations) == len(counts) {
					for i, duration := range durations {
						var count = counts[i]

						var data = map[string]interface{}{
							"duration": duration.String(),
							"ns":       duration / time.Nanosecond,
							"count":    count,
						}
						info.Histogram = append(info.Histogram, data)
					}
				}
			}*/
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Print("Error: failed to encode info.json: ", err)
	}
}

// ServeHTTP -- implements GET key
func (p *KeyPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkMethod(w, r, "GET") {
		return
	}

	var key = r.URL.Query()["k"]
	if len(key) == 0 {
		http.Error(w, "missing key parameter(s): 'k'", http.StatusBadRequest)
		return
	}

	ref := p.faster.Track("_faster", "key")
	defer ref.Done()

	var tickers = p.faster.ListTickers()
	var sortedTickers = p.sortHistoryByInterval(tickers)
	var selectedTicker *faster.History
	if len(sortedTickers) > 0 {
		selectedTicker = p.getTicker(r, tickers, sortedTickers[0])
	}

	var tpl = p.templates["key.html"]
	var err = tpl.Execute(w, map[string]interface{}{
		"keyPath":       key[:len(key)-1],
		"keyName":       key[len(key)-1],
		"sortedTickers": sortedTickers,
		"ticker":        selectedTicker,
		"url":           &linkBuilder{*r.URL},
	})

	if err != nil {
		log.Print("Error: failed to render go-faster key.html template: ", err.Error())
	}

}

// returns the History object requested by the user (or 'default' if not specified/found)
func (p *KeyPage) getTicker(r *http.Request, tickers map[string]*faster.History, defaultValue *faster.History) *faster.History {
	var name = r.URL.Query().Get("ticker")
	if name != "" {
		if h, ok := tickers[name]; ok {
			return h
		}
	}
	return defaultValue
}

// sorts the given History tickers by their interval - lowest first
func (p *KeyPage) sortHistoryByInterval(data map[string]*faster.History) []*faster.History {
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
