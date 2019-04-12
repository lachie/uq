package main

import (
	"errors"
	"net/url"
	"strings"
)

type urlSpec struct {
	name     string
	Scheme   string
	Contains string
	Defaults flatURL
}

type flatURL struct {
	Scheme      string
	Host        string
	Port        string
	Username    string
	Password    string
	HasPassword bool
	Path        string

	Opaque   string
	Fragment string

	// TODO query defaults
}

func FromURL(u *url.URL) (fullURL flatURL) {
	fullURL.Scheme = u.Scheme
	fullURL.Host = u.Hostname()
	fullURL.Port = u.Port()

	fullURL.Username = u.User.Username()
	fullURL.Password, fullURL.HasPassword = u.User.Password()

	fullURL.Path = u.Path

	fullURL.Opaque = u.Opaque
	fullURL.Fragment = u.Fragment

	return
}

func (ud flatURL) Merge(mu flatURL) (fullURL flatURL) {
	fullURL = mu

	setIfNotBlank(&fullURL.Scheme, ud.Scheme)
	setIfNotBlank(&fullURL.Host, ud.Host)
	setIfNotBlank(&fullURL.Port, ud.Port)
	setIfNotBlank(&fullURL.Path, ud.Path)

	setIfNotBlank(&fullURL.Opaque, ud.Opaque)
	setIfNotBlank(&fullURL.Fragment, ud.Fragment)

	setIfNotBlank(&fullURL.Username, ud.Username)
	setIfNotBlank(&fullURL.Password, ud.Password)

	// TODO HasPassword
	// TODO query

	return
}

func (ud flatURL) CleanPath() string {
	if ud.Path[0] == '/' {
		return ud.Path[1:]
	} else {
		return ud.Path
	}
}

type matcher func(*flatURL) bool
type matchers []matcher

func (m matchers) all(u *flatURL) bool {
	for _, mt := range m {
		if !mt(u) {
			return false
		}
	}
	return true
}

func scheme(value string) matcher {
	return func(u *flatURL) bool {
		return u.Scheme == value
	}
}

func contains(value string) matcher {
	return func(u *flatURL) bool {
		return strings.Contains(u.Host, value)
	}
}

func (us urlSpec) MatchURL(url *flatURL) (ok bool) {
	matchers := matchers([]matcher{})

	if us.Scheme != "" {
		matchers = append(matchers, scheme(us.Scheme))
	}
	if us.Contains != "" {
		matchers = append(matchers, contains(us.Contains))
	}

	return matchers.all(url)
}

func (us urlSpec) Merge(orig flatURL) flatURL {
	return us.Defaults.Merge(orig)
}

func setIfNotBlank(dest *string, source string) {
	if *dest == "" {
		*dest = source
	}
}

func (uq *Uq) selectUrl(url *flatURL) (u urlSpec, err error) {
	for name, spec := range uq.config.Urls {
		if spec.MatchURL(url) {
			spec.name = name
			return spec, nil
		}
	}

	err = errors.New("unable to match url")
	return
}
