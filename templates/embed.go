package templates

import (
	_ "embed"
	"github.com/go-ee/filegen/gen"
)

//go:embed macros/gitscript.gtpl
var macrosGitScript string

//go:embed clone.sh.gtpl
var clone string

//go:embed pull.sh.gtpl
var pull string

//go:embed status.sh.gtpl
var status string

func MacrosTemplates() (ret []*gen.TemplateSource) {
	ret = []*gen.TemplateSource{
		{
			Text:   macrosGitScript,
			Source: "github.com/go-ee/gitlab/templates/macros/gitscript.gtpl",
		},
	}
	return
}

func Templates() (ret []*gen.TemplateSource) {
	ret = []*gen.TemplateSource{
		{
			Text:   clone,
			Source: "github.com/go-ee/gitlab/templates/clone.sh.gtpl",
		}, {
			Text:   pull,
			Source: "github.com/go-ee/gitlab/templates/pull.sh.gtpl",
		}, {
			Text:   status,
			Source: "github.com/go-ee/gitlab/templates/status.sh.gtpl",
		},
	}
	return
}
