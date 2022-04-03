// Copyright Contributors to the Open Cluster Management project
package utils

import (
	"fmt"
	"io"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	corev1 "k8s.io/api/core/v1"
)

var suffixColor = color.New(color.Bold, color.FgGreen)

// PrefixWriter can write text at various indentation levels.
type PrefixWriter interface {
	// Write writes text with the specified indentation level.
	Write(level int, format string, a ...interface{})
	// WriteLine writes an entire line with no indentation level.
	WriteLine(a ...interface{})
	// Flush forces indentation to be reset.
	Flush()
}

// Each level has 2 spaces for PrefixWriter
const (
	LEVEL_0 = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
	LEVEL_4
)

// prefixWriter implements PrefixWriter
type prefixWriter struct {
	out io.Writer
}

var _ PrefixWriter = &prefixWriter{}

// NewPrefixWriter creates a new PrefixWriter.
func NewPrefixWriter(out io.Writer) PrefixWriter {
	return &prefixWriter{out: out}
}

func (pw *prefixWriter) Write(level int, format string, a ...interface{}) {
	levelSpace := "  "
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += levelSpace
	}
	fmt.Fprintf(pw.out, prefix+format, a...)
}

func (pw *prefixWriter) WriteLine(a ...interface{}) {
	fmt.Fprintln(pw.out, a...)
}

func (pw *prefixWriter) Flush() {
	if f, ok := pw.out.(flusher); ok {
		f.Flush()
	}
}

type flusher interface {
	Flush()
}

func NewSpinner(suffix string, interval time.Duration) *spinner.Spinner {
	return spinner.New(
		spinner.CharSets[14],
		interval,
		spinner.WithColor("green"),
		spinner.WithHiddenCursor(true),
		spinner.WithSuffix(suffixColor.Sprintf(" %s", suffix)))
}

func NewSpinnerWithStatus(suffix string, interval time.Duration, final string, statusFunc func() string) *spinner.Spinner {
	s := NewSpinner(suffix, interval)
	s.FinalMSG = final
	s.PreUpdate = func(s *spinner.Spinner) {
		status := statusFunc()
		if len(status) > 0 {
			s.Suffix = suffixColor.Sprintf(" %s (%s)", suffix, status)
		} else {
			s.Suffix = suffixColor.Sprintf(" %s", suffix)
		}
	}
	return s
}

func GetSpinnerPodStatus(pod *corev1.Pod) string {
	reason := string(pod.Status.Phase)
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			reason = containerStatus.State.Waiting.Reason
		}
	}
	return reason
}
