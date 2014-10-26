package gof

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// router handler
type RouterHandler func(*Context)

// router params
type RouterParams map[string]string

// router interface
type RouterInterface interface {
	Get(pattern string, fn ...RouterHandler)
	Post(pattern string, fn ...RouterHandler)
	Put(pattern string, fn ...RouterHandler)
	Delete(pattern string, fn ...RouterHandler)
	Route(pattern string, method string, fn ...RouterHandler)
	Group(pattern string, fn func(RouterInterface), fn2 ...RouterHandler)
	Find(pattern, method string) (RouterParams, []RouterHandler)
}

// router rule object
type routerRule struct {
	Url      string
	Pattern  string
	Regexp   *regexp.Regexp
	Method   string
	Ext      string
	Params   []string
	Handlers []RouterHandler
}

type Router struct {
	rules        []*routerRule
	group        []string
	groupHandler map[int][]RouterHandler
}

func NewRouter() *Router {
	r := new(Router)
	r.rules = []*routerRule{}
	r.group = []string{}
	r.groupHandler = make(map[int][]RouterHandler)
	return r
}

func (r *Router) buildRule(pattern string, method string, fn ...RouterHandler) *routerRule {
	rule := new(routerRule)
	rule.Url = pattern
	rule.Method = strings.ToUpper(method)
	rule.Handlers = fn
	rule.Params = []string{}

	// trim ext
	ext := filepath.Ext(pattern)
	if ext != "" {
		rule.Ext = ext
		pattern = strings.TrimSuffix(pattern, ext)
	}

	// parse to pattern or regexp
	patternSlice := strings.Split(pattern, "/")
	for i, s := range patternSlice {
		// ignore empty string
		if s == "" {
			continue
		}
		// param string
		if strings.HasPrefix(s, ":") {
			rule.Params = append(rule.Params, s[1:])
			patternSlice[i] = `([\w-%]+)`
		}
	}

	// add ending slash
	if patternSlice[len(patternSlice)-1] != "" {
		patternSlice = append(patternSlice, "")
	}

	pattern = strings.Join(patternSlice, "/")

	if len(rule.Params) > 0 {
		rule.Regexp, _ = regexp.Compile(pattern)
	} else {
		rule.Pattern = pattern
	}

	fmt.Println(patternSlice, rule)
	return rule
}

func (r *Router) Route(pattern string, method string, fn ...RouterHandler) {
	// add prefix with group
	prefix := ""
	prefixHandler := []RouterHandler{}
	fmt.Println(r.group)
	if len(r.group) > 0 {
		for _, g := range r.group {
			prefix = filepath.Join(prefix, g)
		}
		prefixHandler = r.groupHandler[len(r.group)]
	}
	pattern = filepath.Join(prefix, pattern)
	fn = append(prefixHandler, fn...)

	// build and append rule
	methods := strings.Split(method, ",")
	for _, m := range methods {
		if m == "" {
			continue
		}
		rule := r.buildRule(pattern, m, fn...)
		r.rules = append(r.rules, rule)
	}
}

func (r *Router) Get(pattern string, fn ...RouterHandler) {
	r.Route(pattern, "GET", fn...)
}

func (r *Router) Post(pattern string, fn ...RouterHandler) {
	r.Route(pattern, "POST", fn...)
}

func (r *Router) Put(pattern string, fn ...RouterHandler) {
	r.Route(pattern, "PUT", fn...)
}

func (r *Router) Delete(pattern string, fn ...RouterHandler) {
	r.Route(pattern, "DELETE", fn...)
}

func (r *Router) Group(pattern string, rFn func(RouterInterface), fn ...RouterHandler) {
	// append group , so that Route can use it
	r.group = append(r.group, pattern)
	l := len(r.group)
	r.groupHandler[l] = fn
	// call router to apply group
	rFn(r)
	// delete group, make sure it can't affect other routing
	r.group = r.group[:len(r.group)-1]
	delete(r.groupHandler, l)
}

func (r *Router) Find(pattern, method string) (RouterParams, []RouterHandler) {
	// parse ext
	ext := filepath.Ext(pattern)
	if ext != "" {
		pattern = strings.TrimSuffix(pattern, ext) + "/"
	}
	// add ending slash
	if pattern[len(pattern)-1] != 47 {
		pattern += "/"
	}

	// find route
	for _, rt := range r.rules {
		if rt.Method != method {
			continue
		}
		if ext != "" && rt.Ext != ext {
			continue
		}
		if rt.Pattern != "" && rt.Pattern == pattern {
			return nil, rt.Handlers
		}
		if rt.Regexp == nil {
			continue
		}
		if !rt.Regexp.MatchString(pattern) {
			continue
		}
		params := rt.Regexp.FindStringSubmatch(pattern)
		if len(params) != len(rt.Params)+1 {
			continue
		}
		resultParams := make(RouterParams)
		for i, n := range rt.Params {
			resultParams[n] = params[i+1]
		}
		return resultParams, rt.Handlers

	}
	return nil, nil
}
