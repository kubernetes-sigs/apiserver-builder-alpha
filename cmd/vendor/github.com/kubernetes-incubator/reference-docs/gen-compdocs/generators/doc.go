/*
Copyright 2018 The Kubernetes Authors.

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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	// "time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

func GenMarkdownTree(cmd *cobra.Command, dir string, with_title bool) error {
	identity := func(s string) string { return s }
	emptyStr := func(s string) string { return "" }
	return GenMarkdownTreeCustom(cmd, dir, emptyStr, identity, with_title)
}

func GenMarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string, with_title bool) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenMarkdownTreeCustom(c, dir, filePrepender, linkHandler, with_title); err != nil {
			return err
		}
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := GenMarkdownCustom(cmd, f, linkHandler, with_title); err != nil {
		return err
	}
	return nil
}

func GenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string, with_title bool) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	// buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	short := cmd.Short
	long := cmd.Long
	if len(long) == 0 {
		long = short
	}

    if with_title {
		if _, err := fmt.Fprintf(w, "---\ntitle: %s\nnotitle: true\n---\n", name); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "## %s\n\n", name); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "%s\n\n", short); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "### Synopsis\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\n%s\n\n", long); err != nil {
		return err
	}

	if cmd.Runnable() {
		if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.UseLine()); err != nil {
			return err
		}
	}

	if len(cmd.Example) > 0 {
		if _, err := fmt.Fprintf(w, "### Examples\n\n"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "```\n%s\n```\n\n", cmd.Example); err != nil {
			return err
		}
	}

	if err := printOptions(w, cmd, name); err != nil {
		return err
	}
	if hasSeeAlso(cmd) {
		if _, err := fmt.Fprintf(w, "### SEE ALSO\n"); err != nil {
			return err
		}
		if cmd.HasParent() {
			parent := cmd.Parent()
			pname := parent.CommandPath()
			link := pname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", pname, linkHandler(link), parent.Short); err != nil {
				return err
			}
			cmd.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Commands()
		sort.Sort(byName(children))

		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + ".md"
			link = strings.Replace(link, " ", "_", -1)
			if _, err := fmt.Fprintf(w, "* [%s](%s)\t - %s\n", cname, linkHandler(link), child.Short); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func printOptions(w io.Writer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(w)
	if flags.HasFlags() {
		if _, err := fmt.Fprintf(w, "### Options\n\n"); err != nil {
			return err
		}
		usages := flagUsages(flags)
		fmt.Fprintf(w, usages)
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(w)
	if parentFlags.HasFlags() {
		if _, err := fmt.Fprintf(w, "### Options inherited from parent commands\n\n"); err != nil {
			return err
		}
		usages := flagUsages(parentFlags)
		fmt.Fprintf(w, usages)

		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func hasSeeAlso(cmd *cobra.Command) bool {
	if cmd.HasParent() {
		return true
	}
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

func flagUsages(f *pflag.FlagSet) string {
	x := new(bytes.Buffer)

	lines := make([]string, 0)

	lines = append(lines, "<table style=\"width: 100%%; table-layout: fixed;\">\n  <colgroup>\n" +
		"    <col span=\"1\" style=\"width: 10px;\" />\n" +
		"    <col span=\"1\" />\n" +
		"  </colgroup>\n" +
        "  <tbody>\n")
	f.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 || flag.Hidden {
			return
		}

		line := "    <tr>\n      <td colspan=\"2\">"
		if len(flag.Shorthand) > 0 && len(flag.ShorthandDeprecated) == 0 {
			line += fmt.Sprintf("-%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line += fmt.Sprintf("--%s", flag.Name)
		}

		varname, usage := UnquoteUsage(flag)
		if len(varname) > 0 {
			line += " " + varname
		}
		if len(flag.NoOptDefVal) > 0 {
			switch flag.Value.Type() {
			case "string":
				line += fmt.Sprintf("[=\"%s\"]", flag.NoOptDefVal)
			case "bool":
				if flag.NoOptDefVal != "true" {
					line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
				}
			default:
				line += fmt.Sprintf("[=%s]", flag.NoOptDefVal)
			}
		}
		if !defaultIsZeroValue(flag) {
			if flag.Value.Type() == "string" {
				// There are cases where the string is very very long, split
				// it to mutiple lines manually
				defaultValue := flag.DefValue
				if len(defaultValue) > 40 {
					defaultValue = strings.Replace(defaultValue, ",", ",<br />", -1)
				}
				line += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Default: \"%s\"", defaultValue)
			} else {
				line += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Default: %s", flag.DefValue)
			}
		}
		line += "</td>\n    </tr>\n    <tr>\n      <td></td><td style=\"line-height: 130%%; word-wrap: break-word;\">"

		// escape '<' and '>', force wrap for "\n"
		usage = strings.Replace(usage, "<", "&lt;", -1)
		usage = strings.Replace(usage, ">", "&gt;", -1)
		usage = strings.Replace(usage, "\n", "<br/>", -1)
		line += usage + "</td>\n    </tr>\n"

		lines = append(lines, line)
	})
	lines = append(lines, "  </tbody>\n</table>\n\n")

	for _, line := range lines {
		// fmt.Fprintln(x, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
		fmt.Fprintln(x, line)
	}

	return x.String()
}

func defaultIsZeroValue(f *pflag.Flag) bool {
	switch f.Value.Type() {
	case "bool":
		return f.DefValue == "false"
	case "duration":
		return f.DefValue == "0" || f.DefValue == "0s"
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint16", "uint32", "count", "float32", "float64":
		return f.DefValue == "0"
	case "string":
		return f.DefValue == ""
	case "ip", "ipMask", "ipNet":
		return f.DefValue == "<nil>"
	case "intSlice", "stringSlice", "stringArray":
		return f.DefValue == "[]"
	default:
		switch f.Value.String() {
		case "false":
			return true
		case "<nil>":
			return true
		case "":
			return true
		case "0":
			return true
		}
		return false
	}
}

func UnquoteUsage(flag *pflag.Flag) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}

	name = flag.Value.Type()
	switch name {
	case "bool":
		name = ""
	case "float64":
		name = "float"
	case "int64":
		name = "int"
	case "uint64":
		name = "uint"
	}

	return
}
