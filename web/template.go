package web

import (
	"html/template"
	"net/url"
	"path"

	"github.com/mreithub/go-faster/web/internal"
)

func parseTemplates(prefix string) (map[string]*template.Template, error) {
	var tpls = map[string]string{
		"index.html": internal.IndexHTML,
		"key.html":   internal.KeyHTML,
	}
	var rc = map[string]*template.Template{}
	var err error

	var funcs = map[string]interface{}{
		"keyLink": func(key []string) string {
			var query = url.Values{
				"k": key,
			}
			var rc = url.URL{
				Path:     path.Join(prefix, "key"),
				RawQuery: query.Encode(),
			}
			return rc.String()
		},
	}

	for name, html := range tpls {
		var tpl *template.Template
		if tpl, err = template.New(name).Funcs(funcs).Parse(html); err != nil {
			return nil, err
		}
		rc[name] = tpl
	}

	return rc, nil
}
