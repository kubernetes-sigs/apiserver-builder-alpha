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
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators/api"
)

type TOCItem struct {
	Level int
	Title string
	Link string
	File string
	SubSections []*TOCItem
}

type TOC struct {
	Title string
	Copyright string
	Sections []*TOCItem
}

type HTMLWriter struct {
	Config *api.Config
	TOC TOC
	CurrentSection *TOCItem
}

func NewHTMLWriter(config *api.Config, copyright, title string) DocWriter {
	writer := HTMLWriter{
		Config: config,
		TOC: TOC{
			Copyright: copyright,
			Title: title,
			Sections: []*TOCItem{},
		},
	}
	return &writer
}

func (h *HTMLWriter) Extension() string {
	return ".html"
}

func (h *HTMLWriter) WriteOverview() {
	fn := "_overview.html"
	writeStaticFile("Overview", fn, h.DefaultStaticContent("Overview"))
	item := TOCItem {
		Level: 1,
		Title: "Overview",
		Link: "-strong-api-overview-strong-",
		File: fn,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

func (h *HTMLWriter) WriteResourceCategory(name, file string) {
	writeStaticFile(name, file + ".html", h.DefaultStaticContent(name))
	link := strings.Replace(strings.ToLower(name), " ", "-", -1)
	item := TOCItem {
		Level: 1,
		Title: strings.ToUpper(name),
		Link: "-strong-" + link + "-strong-",
		File: file + ".html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

func (h *HTMLWriter) DefaultStaticContent(title string) string {
	titleID := strings.ToLower(strings.Replace(title, " ", "-", -1))
	return fmt.Sprintf("<H1 id=\"-strong-%s-strong-\"><STRONG>%s</STRONG></H1>\n", titleID, title)
}

func (h *HTMLWriter) writeOtherVersions(w io.Writer, d *api.Definition) {
	if d.OtherVersions.Len() == 0 {
		return
	}

	fmt.Fprint(w, "<DIV class=\"alert alert-success col-md-8\"><I class=\"fa fa-toggle-right\"></I> Other API versions of this object exist:\n")
	for _, v := range d.OtherVersions {
		fmt.Fprintf(w, "%s\n", v.VersionLink())
	}
	fmt.Fprintf(w, "</DIV>\n")
}

func (h *HTMLWriter) writeAppearsIn(w io.Writer, d *api.Definition) {
	if d.AppearsIn.Len() != 0 {
		fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-info-circle\"></I> Appears In:\n <UL>\n")
		for _, a := range d.AppearsIn {
			fmt.Fprintf(w, "  <LI>%s</LI>\n", a.FullHrefLink())
		}
		fmt.Fprintf(w, " </UL>\n</DIV>\n")
	}
}

func (h *HTMLWriter) writeFields(w io.Writer, d *api.Definition) {
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Field</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")

	for _, field := range d.Fields {
		fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE>", field.Name)
		if field.Link() != "" {
			fmt.Fprintf(w, "<BR /><I>%s</I>", field.FullLink())
		}
		if field.PatchStrategy != "" {
			fmt.Fprintf(w, "<BR /><B>patch strategy</B>: <I>%s</I>", field.PatchStrategy)
		}
		if field.PatchMergeKey != "" {
			fmt.Fprintf(w, "<BR /><B>patch merge key</B>: <I>%s</I>", field.PatchMergeKey)
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", field.DescriptionWithEntities)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *HTMLWriter) WriteDefinitionsOverview() {
	writeStaticFile("Definitions", "_definitions.html", h.DefaultStaticContent("Definitions"))
	item := TOCItem {
		Level: 1,
		Title: "DEFINITIONS",
		Link: "-strong-definitions-strong-",
		File: "_definitions.html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

func (h *HTMLWriter) WriteDefinition(d *api.Definition) {
	fn := "_" + definitionFileName(d) + ".html"
	path := *api.ConfigDir + "/includes/" + fn
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	nvg := fmt.Sprintf("%s %s %s", d.Name, d.Version, d.GroupDisplayName())
	linkID := getLink(nvg)
	fmt.Fprintf(f, "<H2 id=\"%s\">%s</H2>\n", linkID, nvg)
	fmt.Fprintf(f, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH><TH>Kind</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(f, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
	            d.GroupDisplayName(), d.Version, d.Name)
	fmt.Fprintf(f, "</TBODY>\n</TABLE>\n")

	fmt.Fprintf(f, "<P>%s</P>\n", d.DescriptionWithEntities)
	h.writeOtherVersions(f, d)
	h.writeAppearsIn(f, d)
	h.writeFields(f, d)

	item := TOCItem {
		Level: 2,
		Title: nvg,
		Link: linkID,
		File: fn,
	}
	h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &item)
}

func (h *HTMLWriter) writeSample(w io.Writer, d *api.Definition) {
	if d.Sample.Sample == "" {
		return
	}

	note := d.Sample.Note
	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", linkID, linkID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">%s example</BUTTON>\n", sType)
	}

	for _, s := range d.GetSamples() {
		sType := strings.Split(s.Tab, ":")[1]
		linkID := sType + "-" + d.LinkID()
		lType := strings.Split(s.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", linkID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", note)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", sType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", html.EscapeString(s.Text))
	}
}

func (h *HTMLWriter) writeOperationSample(w io.Writer, req bool, op string, examples []api.ExampleText) {
	// e.Tab bdocs-tab:kubectl  | bdocs-tab:curl
	// e.Msg `kubectl` Command  | Output | Response Body | `curl` Command (*requires `kubectl proxy` to be running*)
	// e.Type bdocs-tab:kubectl_shell
	// e.Text <actual command>

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		var btnText string
		if req {
			sampleID = "req-" + eType + "-" + op
			btnText = eType + " request"
		} else {
			sampleID = "res-" + eType + "-" + op
			btnText = eType + " response"
		}
		fmt.Fprintf(w, "<BUTTON class=\"btn btn-info\" type=\"button\" data-toggle=\"collapse\"\n")
		fmt.Fprintf(w, "  data-target=\"#%s\" aria-controls=\"%s\"\n", sampleID, sampleID)
		fmt.Fprintf(w, "  aria-expanded=\"false\">%s example</BUTTON>\n", btnText)
	}

	for _, e := range examples {
		eType := strings.Split(e.Tab, ":")[1]
		var sampleID string
		if req {
			sampleID = "req-" + eType + "-" + op
		} else {
			sampleID = "res-" + eType + "-" + op
		}
		msg := e.Msg
		if eType == "curl" && strings.Contains(msg, "proxy") {
			msg = "<CODE>curl</CODE> command (<I>requires <code>kubectl proxy</code> to be running</I>)"
		} else if eType == "kubectl" && strings.Contains(msg, "Command") { // `kubectl` command
			msg = "<CODE>kubectl</CODE> command"
		}
		lType := strings.Split(e.Type, ":")[1]
		lang := strings.Split(lType, "_")[1]
		fmt.Fprintf(w, "<DIV class=\"collapse\" id=\"%s\">\n", sampleID)
		fmt.Fprintf(w, "  <DIV class=\"panel panel-default\">\n<DIV class=\"panel-heading\">%s</DIV>\n", msg)
		fmt.Fprintf(w, "  <DIV class=\"panel-body\">\n<PRE class=\"%s\">", eType)
		fmt.Fprintf(w, "<CODE class=\"lang-%s\">\n", lang)
		// TODO: Add language highlight
		fmt.Fprintf(w, "%s\n</CODE></PRE></DIV></DIV></DIV>\n", e.Text)
	}
}

func (h *HTMLWriter) writeParams(w io.Writer, title string, params api.Fields) {
	fmt.Fprintf(w, "<H3>%s</H3>\n", title)
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Parameter</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range params {
		fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE>", p.Name)
		if p.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}

func (h *HTMLWriter) writeRequestParams(w io.Writer, o *api.Operation) {
	// Operation path params
	if o.PathParams.Len() > 0 {
		h.writeParams(w, "Path Parameters", o.PathParams)
	}

	// operation query params
	if o.QueryParams.Len() > 0 {
		h.writeParams(w, "Query Parameters", o.QueryParams)
	}

	// operation body params
	if o.BodyParams.Len() > 0 {
		h.writeParams(w, "Body Parameters", o.BodyParams)
	}
}

func (h *HTMLWriter) writeResponseParams(w io.Writer, o *api.Operation) {
	if o.HttpResponses.Len() == 0 {
		return
	}

	fmt.Fprintf(w, "<H3>Response</H3>\n")
	fmt.Fprintf(w, "<TABLE>\n<THEAD><TR><TH>Code</TH><TH>Description</TH></TR></THEAD>\n<TBODY>\n")
	for _, p := range o.HttpResponses {
		fmt.Fprintf(w, "<TR><TD>%s", p.Name)
		if p.Field.Link() != "" {
			fmt.Fprintf(w, "<br /><I>%s</I>", p.Field.FullLink())
		}
		fmt.Fprintf(w, "</TD><TD>%s</TD></TR>\n", p.Field.Description)
	}
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")
}


func (h *HTMLWriter) WriteResource(r *api.Resource) {
	fn := "_" + conceptFileName(r.Definition) + ".html"
	path := *api.ConfigDir + "/includes/" + fn
	w, err := os.Create(path)
	defer w.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	dvg := fmt.Sprintf("%s %s %s", r.Name, r.Definition.Version, r.Definition.GroupDisplayName())
	linkID := getLink(dvg)
	fmt.Fprintf(w, "<H1 id=\"%s\">%s</H1>\n", linkID, dvg)

	h.writeSample(w, r.Definition)

	// GVK
	fmt.Fprintf(w, "<TABLE class=\"col-md-8\">\n<THEAD><TR><TH>Group</TH><TH>Version</TH><TH>Kind</TH></TR></THEAD>\n<TBODY>\n")
	fmt.Fprintf(w, "<TR><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD><TD><CODE>%s</CODE></TD></TR>\n",
	            r.Definition.GroupDisplayName(), r.Definition.Version, r.Name)
	fmt.Fprintf(w, "</TBODY>\n</TABLE>\n")

	if r.DescriptionWarning != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-warning col-md-8\"><P><I class=\"fa fa-warning\"></I> <B>Warning:</B></P><P>%s</P></DIV>\n", r.DescriptionWarning)
	}
	if r.DescriptionNote != "" {
		fmt.Fprintf(w, "<DIV class=\"alert alert-info col-md-8\"><I class=\"fa fa-bullhorn\"></I> %s</DIV>\n", r.DescriptionNote)
	}

	h.writeOtherVersions(w, r.Definition)
	h.writeAppearsIn(w, r.Definition)
	h.writeFields(w, r.Definition)

	// Inline
	if r.Definition.Inline.Len() > 0 {
		for _, d := range r.Definition.Inline {
			fmt.Fprintf(w, "<H3 id=\"%s\">%s %s %s</H3>\n", d.LinkID(), d.Name, d.Version, d.Group)
			h.writeAppearsIn(w, d)
			h.writeFields(w, d)
		}
	}

	item := TOCItem {
		Level: 1,
		Title: dvg,
		Link: linkID,
		File: fn,
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item

	// Operations
	if len(r.Definition.OperationCategories) == 0 {
		return
	}

	for _, c := range r.Definition.OperationCategories {
		if len(c.Operations) == 0 {
			continue
		}
		catID := strings.Replace(strings.ToLower(c.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
		catID = "-strong-" + catID + "-strong-"
		fmt.Fprintf(w, "<H2 id=\"%s\"><STRONG>%s</STRONG></H2>\n", catID, c.Name)
		OCItem := TOCItem {
			Level: 2,
			Title: c.Name,
			Link: catID,
		}
		h.CurrentSection.SubSections = append(h.CurrentSection.SubSections, &OCItem)

		for _, o := range c.Operations {
			opID := strings.Replace(strings.ToLower(o.Type.Name), " ", "-", -1) + "-" + r.Definition.LinkID()
			fmt.Fprintf(w, "<H2 id=\"%s\">%s</H2>\n", opID, o.Type.Name)
			OPItem := TOCItem {
				Level: 2,
				Title: o.Type.Name,
				Link: opID,
			}
			OCItem.SubSections = append(OCItem.SubSections, &OPItem)

			// Example requests
			requests := o.GetExampleRequests()
			if len(requests) > 0 {
				h.writeOperationSample(w, true, opID, requests)
			}
			// Example responses
			responses := o.GetExampleResponses()
			if len(responses) > 0 {
				h.writeOperationSample(w, false, opID, responses)
			}

			fmt.Fprintf(w, "<P>%s</P>\n", o.Description())
			fmt.Fprintf(w, "<H3>HTTP Request</H3>\n")
			fmt.Fprintf(w, "<CODE>%s</CODE>\n", o.GetDisplayHttp())

			h.writeRequestParams(w, o)
			h.writeResponseParams(w, o)
		}
	}
}

func (h *HTMLWriter) WriteOldVersionsOverview() {
	writeStaticFile("Old Versions", "_oldversions.html", h.DefaultStaticContent("Old Versions"))
	item := TOCItem {
		Level: 1,
		Title: "OLD API VERSIONS",
		Link: "-strong-old-api-versions-strong-",
		File: "_oldversions.html",
	}
	h.TOC.Sections = append(h.TOC.Sections, &item)
	h.CurrentSection = &item
}

func (h *HTMLWriter) generateNavContent() string {
	nav := ""
	for _, sec := range h.TOC.Sections {
		// class for level-1 navigation item
		nav += "<UL>\n"
		if strings.Contains(sec.Link, "strong") {
			nav += fmt.Sprintf(" <LI class=\"nav-level-1 strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n",
							   sec.Link, sec.Title)
		} else {
			nav += fmt.Sprintf(" <LI class=\"nav-level-1\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
							   sec.Link, sec.Title)
		}

		// close H1 items which have no subsections or strong navs
		if len(sec.SubSections) == 0 || (sec.Level == 1 && strings.Contains(sec.Link, "strong")) {
			nav += "</UL>\n"
		}

		// short circuit to next if no sub-sections
		if len(sec.SubSections) == 0 {
			continue
		}

		// wrapper1
	    nav += fmt.Sprintf(" <UL id=\"%s-nav\" style=\"display: none;\">\n", sec.Link)
		for _, sub := range sec.SubSections {
			nav += "  <UL>\n"
			if strings.Contains(sub.Link, "strong") {
				nav += fmt.Sprintf("   <LI class=\"nav-level-%d strong-nav\"><A href=\"#%s\" class=\"nav-item\"><STRONG>%s</STRONG></A></LI>\n",
							        sub.Level, sub.Link, sub.Title)
			} else {
				nav += fmt.Sprintf("   <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
							        sub.Level, sub.Link, sub.Title)
			}
			// close this H1/H2 if possible
			if len(sub.SubSections) == 0 {
				nav += " </UL>\n"
				continue
			}

			// 3rd level
			// another wrapper
			nav += fmt.Sprintf("   <UL id=\"%s-nav\" style=\"display: none;\">\n", sub.Link)
			for _, subsub := range sub.SubSections {
				nav += fmt.Sprintf("    <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n", subsub.Level, subsub.Link, subsub.Title)
				if len(subsub.SubSections) == 0 {
					continue
				}

				fmt.Printf("*** found third level!\n")
				nav += fmt.Sprintf("   <UL id=\"%s-nav\" style=\"display: none;\">\n", subsub.Link)
				for _, subsubsub := range subsub.SubSections {
					nav += fmt.Sprintf("    <LI class=\"nav-level-%d\"><A href=\"#%s\" class=\"nav-item\">%s</A></LI>\n",
					                   subsubsub.Level, subsubsub.Link, subsubsub.Title)
				}
				nav += "   </UL>\n"
			}
			// end wrapper2
			nav += "   </UL>\n"
			nav += "  </UL>\n"
		}
		// end wrapper1
		nav += " </UL>\n"
		// end top UL
		nav += "</UL>\n"
	}

	return nav
}

func (h *HTMLWriter) generateNavJS() {
	// generate NavData
	var tmp string
	flatToc := []string{}
	os.MkdirAll(*api.ConfigDir + "/build", os.ModePerm)

	navjs, err := os.Create(*api.ConfigDir + "/build/navData.js")
	defer navjs.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	s1 := []string{}
	for _, sec := range h.TOC.Sections {
		flatToc = append([]string{"\"" + sec.Link + "\""}, flatToc...)
		s2 := []string{}
		for _, sub := range sec.SubSections {
			flatToc = append([]string{"\"" + sub.Link + "\""}, flatToc...)
			s3 := []string{}
			for _, subsub := range sub.SubSections {
				flatToc = append([]string{"\"" + subsub.Link + "\""}, flatToc...)
				s4 := []string{}
				for _, subsubsub := range subsub.SubSections {
					flatToc = append([]string{"\"" + subsubsub.Link + "\""}, flatToc...)
					tmp = "{\"section\":\"" + subsubsub.Link + "\"}"
					s4 = append([]string{tmp}, s4...)
				}
				tmp = "{\"section\":\"" + subsub.Link + "\",\"subsections\":[" + strings.Join(s4, ",") + "]}"
				s3 = append([]string{tmp}, s3...)
			}
			tmp = "{\"section\":\"" + sub.Link + "\",\"subsections\":[" + strings.Join(s3, ",") + "]}"
			s2 = append([]string{tmp}, s2...)
		}

		if strings.Contains(sec.Link, "strong") {
			tmp =  "{\"section\":\"" + sec.Link + "\",\"subsections\":[]}"
			s2res := strings.Join(s2, ",")
			if len(s2res) > 0 {
				tmp = s2res + "," + tmp
			}
		} else {
			tmp = "{\"section\":\"" + sec.Link + "\",\"subsections\":[" + strings.Join(s2, ",") + "]}"
		}
		s1 = append([]string{tmp}, s1...)
	}
	fmt.Fprintf(navjs, "(function(){navData={\"toc\":[" + strings.Join(s1, ",") + "],\"flatToc\":[" + strings.Join(flatToc, ",") + "]};})();")
}

func (h *HTMLWriter) generateHTML(navContent string) {
	html, err := os.Create(*api.ConfigDir + "/build/index.html")
	defer html.Close()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	fmt.Fprintf(html, "<!DOCTYPE html>\n<HTML>\n<HEAD>\n<META charset=\"UTF-8\">\n")
	fmt.Fprintf(html, "<TITLE>%s</TITLE>\n", h.TOC.Title)
	fmt.Fprintf(html, "<LINK rel=\"shortcut icon\" href=\"favicon.ico\" type=\"image/vnd.microsoft.icon\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/bootstrap.min.css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/font-awesome.min.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "<LINK rel=\"stylesheet\" href=\"css/stylesheet.css\" type=\"text/css\">\n")
	fmt.Fprintf(html, "</HEAD>\n<BODY>\n")
	fmt.Fprintf(html, "<DIV id=\"sidebar-wrapper\" class=\"side-nav side-bar-nav\">\n")

	// html buffer
	buf := ""
	for _, sec := range h.TOC.Sections {
		fmt.Printf("Collecting %s ... ", sec.File)
		content, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sec.File))
		if err == nil {
			buf += string(content)
		}
		fmt.Printf("OK\n")

		for _, sub := range sec.SubSections {
			if len(sub.File) > 0 {
				subdata, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", sub.File))
				fmt.Printf("Collecting %s ... ", sub.File)
				if err == nil {
					buf += string(subdata)
					fmt.Printf("OK\n")
				}
			}

			for _, subsub := range sub.SubSections {
				if len(subsub.File) > 0 {
					subsubdata, err := ioutil.ReadFile(filepath.Join(*api.ConfigDir, "includes", subsub.File))
					fmt.Printf("Collecting %s ... ", subsub.File)
					if err == nil {
						buf += string(subsubdata)
						fmt.Printf("OK\n")
					}
				}
			}
		}
	}

	fmt.Fprintf(html, "%s", navContent)
	fmt.Fprintf(html, "<BR/><DIV class=\"copyright\">%s</DIV></DIV>\n", h.TOC.Copyright)
	fmt.Fprintf(html, "<DIV id=\"wrapper\">\n")
	fmt.Fprintf(html, "<DIV id=\"page-content-wrapper\" class=\"body-content container-fluid\">\n")
	fmt.Fprintf(html, "%s", string(buf))
	fmt.Fprintf(html, "</DIV>\n</DIV>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/jquery-3.2.1.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"jquery.scrollTo.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"/js/bootstrap-3.3.7.min.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"navData.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "<SCRIPT src=\"scroll.js\"></SCRIPT>\n")
	fmt.Fprintf(html, "</BODY>\n</HTML>\n")
}

func (h *HTMLWriter) Finalize() {
	// generate NavData
	os.MkdirAll(*api.ConfigDir + "/build", os.ModePerm)

	h.generateNavJS()
	navContent := h.generateNavContent()
	h.generateHTML(navContent)
}
