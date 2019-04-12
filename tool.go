package main

import (
	"html/template"
	"strings"

	"github.com/pkg/errors"
)

type toolSpec struct {
	Url      string
	Nickname string
	BasedOn  string
	Format   string
	Defaults map[string]interface{}
}

func (t toolSpec) MergeUrl(url flatURL, config config) (cmdline string, err error) {
	template, err := template.New("format").Parse(t.Format)
	if err != nil {
		return
	}

	b := &strings.Builder{}

	err = template.Execute(b, url)
	if err != nil {
		return
	}

	cmdline = b.String()
	return
}

func (uq *Uq) selectTool(nickname string, urlSpec urlSpec) (tool toolSpec, err error) {
	for _, tool := range uq.config.Tools {
		if tool.Url == urlSpec.name && tool.Nickname == nickname {
			return tool, nil
		}
	}

	err = errors.Errorf("no tool found nickname %s url %s", nickname, urlSpec.name)
	return
}
