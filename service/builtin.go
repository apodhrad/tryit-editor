package service

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var SERVICE_CAT Service = Service{
	Exec:     "cat",
	Examples: []string{"examples/default.txt"},
}

var BUILTIN_SERVICE_HTML Service = Service{
	Exec:     "html [built-in]",
	Examples: []string{"service/html [built-in]/example/items.html"},
	run: func(inputFile string) ([]byte, error) {
		return os.ReadFile(inputFile)
	},
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

var BUILTIN_SERVICE_MARKDOWN Service = Service{
	Exec:     "markdown [built-in]",
	Examples: []string{"service/markdown [built-in]/example/readme.md"},
	run: func(inputFile string) ([]byte, error) {
		input, err := os.ReadFile(inputFile)
		if err != nil {
			return input, err
		}
		output := mdToHTML(input)
		return output, nil
	},
}
