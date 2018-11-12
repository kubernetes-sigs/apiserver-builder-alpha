/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators/api"
)

type Manifest struct {
	Docs            []Doc   `json:"docs,omitempty"`
	Title           string  `json:"title,omitempty"`
	Copyright       string  `json:"copyright,omitempty"`
}

type MarkdownWriter struct {
	Config *api.Config
	Manifest Manifest
}

func NewMarkdownWriter(config *api.Config, copyright, title string) DocWriter {
	writer := MarkdownWriter{
		Config: config,
		Manifest: Manifest{
			Copyright: copyright,
			Title: title,
		},
	}
	return &writer
}

func (m *MarkdownWriter) Extension() string {
	return ".md"
}

func (m *MarkdownWriter) DefaultStaticContent(title string) string {
	return fmt.Sprintf("# <strong>%s</strong>\n\n----------\n\n", title)
}

func (m *MarkdownWriter) WriteOverview() {
	writeStaticFile("Overview", "_overview.md", m.DefaultStaticContent("Overview"))
	if *api.BuildOps {
		m.Manifest.Docs = append(m.Manifest.Docs, Doc{"_overview.md"})
	}
}

func (m *MarkdownWriter) WriteResourceCategory(name, file string) {
	writeStaticFile(name, file + ".md", m.DefaultStaticContent(name))
	m.Manifest.Docs = append(m.Manifest.Docs, Doc{file + ".md"})
}

func (m *MarkdownWriter) writeFields(w io.Writer, d *api.Definition) {
	fmt.Fprintf(w, "Field        | Description\n------------ | -----------\n")
	for _, field := range d.Fields {
		fmt.Fprintf(w, "`%s`", field.Name)
		if field.Link() != "" {
			fmt.Fprintf(w, "<br /> *%s*", field.Link())
		}
		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "<br /> **patch strategy**: *%s*", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "<br /> **patch merge key**: *%s*", field.PatchMergeKey)
		}
		fmt.Fprintf(w, " | %s\n", field.DescriptionWithEntities)
	}
}

func (m *MarkdownWriter) writeOtherVersions(w io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() != 0 {
		fmt.Fprintf(w, "<aside class=\"notice\">Other API versions of this object exist:\n")
		for _, v := range d.OtherVersions {
			fmt.Fprintf(w, "%s\n", v.VersionLink())
		}
		fmt.Fprintf(w, "</aside>\n\n")
	}
	fmt.Fprintf(w, "%s\n\n", d.DescriptionWithEntities)
}

func (m *MarkdownWriter) writeAppearsIn(w io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		fmt.Fprintf(w, "<aside class=\"notice\">\nAppears In:\n\n<ul>\n")
		for _, a := range d.AppearsIn {
			fmt.Fprintf(w, "<li>%s</li>\n", a.FullHrefLink())
		}
		fmt.Fprintf(w, "</ul></aside>\n\n")
	}
}

func (m *MarkdownWriter) WriteDefinition(d *api.Definition) {
	fn := "_" + definitionFileName(d) + ".md"
	path := *api.ConfigDir + "/includes/" + fn
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	fmt.Fprintf(f, "## %s %s %s\n\n", d.Name, d.Version, d.Group)
	fmt.Fprintf(f, "Group        | Version    | Kind\n------------ | ---------- | -----------\n")
	fmt.Fprintf(f, "`%s` | `%s` | `%s`\n", d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "\n")

	m.writeOtherVersions(f, d)
	m.writeAppearsIn(f, d)
	m.writeFields(f, d)
	m.Manifest.Docs = append(m.Manifest.Docs, Doc{fn})
}

func (m *MarkdownWriter) WriteResource(r *api.Resource) {
	fn := "_" + conceptFileName(r.Definition) + ".md"
	path := *api.ConfigDir + "/includes/" + fn
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(f, "-----------\n\n")
	fmt.Fprintf(f, "# %s %s", r.Name, r.Definition.Version)
	if r.Definition.ShowGroup {
		fmt.Fprintf(f, " %s\n", r.Definition.GroupDisplayName())
	} else {
		fmt.Fprintf(f, "\n")
	}

	if r.Definition.Sample.Sample != "" {
		note := r.Definition.Sample.Note
		for _, t := range r.Definition.GetSamples() {
			// 't' is a ExampleText
			fmt.Fprintf(f, ">%s %s\n\n", t.Tab, note)
			fmt.Fprintf(f, "```%s\n%s\n```\n\n", t.Type, t.Text)
		}
	}

	// GVK
	fmt.Fprintf(f, "Group        | Version    | Kind\n------------ | ---------- | -----------\n")
	fmt.Fprintf(f, "`%s` | `%s` | `%s`\n\n", r.Definition.GroupDisplayName(), r.Definition.Version, r.Name)

	if r.DescriptionWarning != "" {
		fmt.Fprintf(f, "<aside class=\"warning\">%s</aside>\n\n", r.DescriptionWarning)
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(f, "<aside class=\"notice\">%s</aside>\n\n", r.DescriptionNote)
	}

	m.writeOtherVersions(f, r.Definition)
	m.writeAppearsIn(f, r.Definition)
	m.writeFields(f, r.Definition)

	fmt.Fprintf(f, "\n")
	if r.Definition.Inline.Len() > 0 {
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(f, "### %s %s %s\n", d.Name, d.Version, d.Group)
			m.writeAppearsIn(f, d)
			m.writeFields(f, d)
		}
	}

	if len(r.Definition.OperationCategories) > 0 {
		for _, c := range r.Definition.OperationCategories {
			if len(c.Operations) > 0 {
				fmt.Fprintf(f, "## <strong>%s</strong>\n", c.Name)
				for _, o := range c.Operations {
					fmt.Fprintf(f, "\n## %s\n", o.Type.Name)

					// Example requests
					requests := o.GetExampleRequests()
					if len(requests) > 0 {
						for _, r := range requests {
							fmt.Fprintf(f, ">%s %s\n\n", r.Tab, r.Msg)
							fmt.Fprintf(f, "```%s\n%s```\n\n", r.Type, r.Text)
						}
					}
					// Example responses
					responses := o.GetExampleResponses()
					if len(responses) > 0 {
						for _, r := range responses {
							fmt.Fprintf(f, ">%s %s\n\n", r.Tab, r.Msg)
							fmt.Fprintf(f, "```%s\n%s```\n\n", r.Type, r.Text)
						}
					}

					fmt.Fprintf(f, "%s\n", o.Description())
					fmt.Fprintf(f, "\n### HTTP Request\n\n`%s`\n", o.GetDisplayHttp())

					// Operation path params
					if o.PathParams.Len() > 0 {
						fmt.Fprintf(f, "\n### Path Parameters\n\n")
						fmt.Fprintf(f, "Parameter    | Description\n------------ | -----------\n")
						for _, p := range o.PathParams {
							fmt.Fprintf(f, "`%s`", p.Name)
							if p.Link() != "" {
								fmt.Fprintf(f, "<br /> *%s*", p.Link())
							}
							fmt.Fprintf(f, " | %s\n", p.Description)
						}
					}

					// operation query params
					if o.QueryParams.Len() > 0 {
						fmt.Fprintf(f, "\n### Query Parameters\n\n")
						fmt.Fprintf(f, "Parameter    | Description\n------------ | -----------\n")
						for _, p := range o.QueryParams {
							fmt.Fprintf(f, "`%s`", p.Name)
							if p.Link() != "" {
								fmt.Fprintf(f, "<br /> *%s*", p.Link())
							}
							fmt.Fprintf(f, " | %s\n", p.Description)
						}
					}
					// operation body params
					if o.BodyParams.Len() > 0 {
						fmt.Fprintf(f, "\n### Body Parameters\n\n")
						fmt.Fprintf(f, "Parameter    | Description\n------------ | -----------\n")
						for _, p := range o.BodyParams {
							fmt.Fprintf(f, "`%s`", p.Name)
							if p.Link() != "" {
								fmt.Fprintf(f, "<br /> *%s*", p.Link())
							}
							fmt.Fprintf(f, " | %s\n", p.Description)
						}
					}

					// operation response body
					if o.HttpResponses.Len() > 0 {
						fmt.Fprintf(f, "\n### Response\n\n")
						fmt.Fprintf(f, "Code         | Description\n------------ | -----------\n")
						for _, r := range o.HttpResponses {
							fmt.Fprintf(f, "%s ", r.Code)
							if r.Field.Link() != "" {
								fmt.Fprintf(f, "<br /> *%s*", r.Field.Link())
							}
							fmt.Fprintf(f, " | %s\n", r.Field.Description)
						}
					}
				}
			}
		}
	}

	m.Manifest.Docs = append(m.Manifest.Docs, Doc{fn})
}

func (m *MarkdownWriter) WriteDefinitionsOverview() {
	writeStaticFile("Definitions", "_definitions.md", m.DefaultStaticContent("Definitions"))
	m.Manifest.Docs = append(m.Manifest.Docs, Doc{"_definitions.md"})
}

func (m *MarkdownWriter) WriteOldVersionsOverview() {
	writeStaticFile("Old Versions", "_oldversions.md", m.DefaultStaticContent("Old Versions"))
	m.Manifest.Docs = append(m.Manifest.Docs, Doc{"_oldversions.md"})
}

func (m *MarkdownWriter) Finalize() {
	// Write out the json manifest
	jsonbytes, err := json.MarshalIndent(m.Manifest, "", "  ")
	if err != nil {
		fmt.Printf("Could not Marshal manfiest %+v due to error: %v.\n", m.Manifest, err)
	} else {
		jsonfile, err := os.Create(*api.ConfigDir + "/" + JsonOutputFile)
		if err != nil {
			fmt.Printf("Could not create file %s due to error: %v.\n", JsonOutputFile, err)
		} else {
			defer jsonfile.Close()
			_, err := jsonfile.Write(jsonbytes)
			if err != nil {
				fmt.Printf("Failed to write bytes %s to file %s: %v.\n", jsonbytes, JsonOutputFile, err)
			}
		}
	}
}
