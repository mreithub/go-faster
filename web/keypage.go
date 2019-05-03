package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mreithub/go-faster/faster"
)

// KeyPage -- impl
type KeyPage struct {
	faster    *faster.Faster
	templates map[string]*template.Template
}

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
	var sortedTickers = sortHistoryByInterval(tickers)
	var selectedTicker *faster.History
	var data []*faster.Snapshot

	if len(sortedTickers) > 0 {
		selectedTicker = tickers[sortedTickers[0].Name] // TODO allow the user to pick another ticker
		data = diff(selectedTicker.ForKey(key...))
	}

	var tpl = p.templates["key.html"]
	var err = tpl.Execute(w, map[string]interface{}{
		"key":     key,
		"keyName": key[len(key)-1],
		"tickers": sortedTickers,
		"data":    data,
	})

	if err != nil {
		log.Print("Error: failed to render go-faster key.html template: ", err.Error())
	}

}
