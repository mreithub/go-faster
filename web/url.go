package web

import "net/url"

type linkBuilder struct {
	URL url.URL
}

func (l *linkBuilder) WithParam(key string, values ...string) *linkBuilder {
	l.WithParamArray(key, values)
	return l
}

// WithParamArray -- flavor of SetParam that accepts an array (for use in templates)
func (l *linkBuilder) WithParamArray(key string, values []string) *linkBuilder {
	var query = l.URL.Query()
	query[key] = values
	l.URL.RawQuery = query.Encode()
	return l
}

func (l *linkBuilder) WithPath(path string) *linkBuilder {
	l.URL.Path = path
	return l
}

func (l *linkBuilder) String() string {
	return l.URL.String()
}
