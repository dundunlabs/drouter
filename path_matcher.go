package prenn

import (
	"fmt"
	"regexp"
)

func newPathMatcher(path string) PathMacher {
	p := generatePattern(path)
	if p == path {
		return PathMacherString(p)
	}
	return PathMacherRegexp{
		regex: regexp.MustCompile("^" + p + "$"),
	}
}

type Params map[string]string

type PathMacher interface {
	Match(path string) (Params, bool)
}

type PathMacherString string

func (p PathMacherString) Match(path string) (Params, bool) {
	return nil, path == string(p)
}

type PathMacherRegexp struct {
	regex *regexp.Regexp
}

func (p PathMacherRegexp) Match(path string) (Params, bool) {
	match := p.regex.FindStringSubmatch(path)
	if len(match) == 0 {
		return nil, false
	}
	params := make(Params)
	for i, name := range p.regex.SubexpNames() {
		if i != 0 {
			params[name] = match[i]
		}
	}
	return params, true
}

var dynamicRegexp = regexp.MustCompile(`[:\*][^/]+`)

var patternGenerators = map[string]func(name string) string{
	":": func(name string) string {
		return fmt.Sprintf("(?P<%s>[^/]+)", name)
	},
	"*": func(name string) string {
		return fmt.Sprintf("(?P<%s>.+)", name)
	},
}

func generatePattern(path string) string {
	return dynamicRegexp.ReplaceAllStringFunc(path, func(s string) string {
		c, name := s[:1], s[1:]
		return patternGenerators[c](name)
	})
}
