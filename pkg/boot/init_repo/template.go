package init_repo

import _ "embed"

var (
	//go:embed templates/go.mod.tpl
	goModTemplate string
	//go:embed templates/apiserver.tpl
	apiserverTemplate string
	//go:embed templates/api.doc.tpl
	apisDocTemplate string
	//go:embed templates/package.doc.tpl
	packageDocTemplate string
)
