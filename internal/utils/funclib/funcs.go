package funclib

import (
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func WhiteSpace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func NonWhiteSpace(s string) bool {
	return !WhiteSpace(s)
}

type number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Gte[N number](a, b N) bool {
	return a >= b
}

func GteCUR[N number](b N) func(N) bool {
	return func(a N) bool {
		return Gte(a, b)
	}
}

func Lte[N number](a, b N) bool {
	return a <= b
}

func LteCUR[N number](b N) func(N) bool {
	return func(a N) bool {
		return Lte(a, b)
	}
}

func Slugify(s string) string {
	pattern := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	s = pattern.ReplaceAllString(s, "")
	return strings.ToLower(strings.Join(strings.Split(s, " "), "-"))
}

func Map[T any, U any](ts []T, fn func(T) U) []U {
	u := make([]U, len(ts))

	for i := range len(ts) {
		u[i] = fn(ts[i])
	}

	return u
}

func ParseToHTML(b []byte) string {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(b)

	flags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: flags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func SanitizeHTML(s string) string {
	return bluemonday.UGCPolicy().Sanitize(s)
}
