package cli

import (
	"bytes"
	"regexp"
	"strings"
)

const (
	esc      = "\033["
	black    = esc + "30m"
	red      = esc + "31m"
	green    = esc + "32m"
	yellow   = esc + "33m"
	blue     = esc + "34m"
	magenta  = esc + "35m"
	cyan     = esc + "36m"
	lgray    = esc + "37m"
	dgray    = esc + "90m"
	lred     = esc + "91m"
	lgreen   = esc + "92m"
	lyellow  = esc + "93m"
	lblue    = esc + "94m"
	lmagenta = esc + "95m"
	lcyan    = esc + "96m"
	white    = esc + "97m"
	bold     = esc + "1m"
	reset    = esc + "0m"
	clear    = esc + "2J" + esc + "f"

	h1     = lyellow
	h2     = yellow
	code   = lblue
	result = green
)

// This is *not* a full-featured markdown formatter, or even a compliant
// one for that matter. It only supports the productions that are
// currently used by documentation strings and will likely not evolve
// much beyond that

type formatter func(string) string

var (
	indent = regexp.MustCompile("^##? |^\\s\\s+")
	hashes = regexp.MustCompile("^##? ")
	ticks  = regexp.MustCompile("`[^`]*`")
	unders = regexp.MustCompile("_[^_]*_")
	stars  = regexp.MustCompile("[*][^*]*[*]")

	lineFormatters = map[*regexp.Regexp]formatter{
		regexp.MustCompile("^#\\s.*$"):   formatHeader1,
		regexp.MustCompile("^##\\s.*$"):  formatHeader2,
		regexp.MustCompile("^\\s\\s.*$"): formatIndent,
	}

	docFormatters = map[*regexp.Regexp]formatter{
		ticks:  trimmedFormatter("`", code),
		unders: trimmedFormatter("_", result),
		stars:  trimmedFormatter("*", bold),
	}
)

func formatMarkdown(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var out []string
	for _, l := range wrapLines(lines) {
		for r, f := range lineFormatters {
			l = r.ReplaceAllStringFunc(l, f)
		}
		out = append(out, l)
	}

	d := strings.Join(out, "\n")
	for r, f := range docFormatters {
		d = r.ReplaceAllStringFunc(d, f)
	}
	return d
}

func getFormatWidth() (width int) {
	return getScreenWidth() - 4
}

func wrapLines(s []string) []string {
	w := getFormatWidth()
	r := []string{}
	for _, e := range s {
		r = append(r, wrapLine(e, w)...)
	}
	return r
}

func wrapLine(s string, w int) []string {
	r := []string{}
	i, s := lineIndent(s)
	il := strippedLen(i)

	var b bytes.Buffer
	b.WriteString(i)
	for _, e := range strings.Split(s, " ") {
		bl := strippedLen(b.String())
		if bl > il {
			if bl+strippedLen(e) >= w {
				r = append(r, b.String())
				b.Reset()
				b.WriteString(i)
			} else if !isEmptyString(b.String()) {
				b.WriteString(" ")
			}
		}
		b.WriteString(e)
	}
	return append(r, b.String())
}

func lineIndent(s string) (string, string) {
	l := indent.FindString(s)
	return l, s[len(l):]
}

func stripDelimited(s string) string {
	s = ticks.ReplaceAllStringFunc(s, trimDelimited)
	s = stars.ReplaceAllStringFunc(s, trimDelimited)
	return hashes.ReplaceAllString(s, "")
}

func strippedLen(s string) int {
	return len(stripDelimited(s))
}

func trimDelimited(s string) string {
	if len(s) > 2 {
		return s[1 : len(s)-1]
	}
	return s[:1]
}

func formatHeader1(s string) string {
	return h1 + s[2:] + reset
}

func formatHeader2(s string) string {
	return h2 + s[3:] + reset
}

func formatIndent(s string) string {
	return code + s + reset
}

func trimmedFormatter(delim string, prefix string) formatter {
	return func(s string) string {
		t := trimDelimited(s)
		if t == delim {
			return t
		}
		return prefix + t + reset
	}
}
