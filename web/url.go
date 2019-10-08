package web

import "net/url"

// urlBuilder -- wraps an url.URL object and adds .With() functions (making it easier to manipulate URLs in templates)
type urlBuilder struct {
	URL url.URL
}

func (l urlBuilder) WithParam(key string, values ...string) urlBuilder {
	return l.WithParamArray(key, values)
}

// WithParamArray -- flavor of SetParam that accepts an array (for use in templates)
func (l urlBuilder) WithParamArray(key string, values []string) urlBuilder {
	var query = l.URL.Query()
	query[key] = values
	l.URL.RawQuery = query.Encode()
	return l
}

func (l urlBuilder) WithPath(path string) urlBuilder {
	l.URL.Path = path
	return l
}

func (l urlBuilder) String() string {
	return l.URL.String()
}
