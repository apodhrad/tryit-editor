package service

import "os"

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
