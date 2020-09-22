package main

import (
	"fmt"
	"text/template"
	"log"
	"os"
)

func main() {
	data := map[string]string{
		"Object":   "*Object",
		"Segment":  "Segment",
		"Storager": "Storager",
	}

	generateT(tmpl, "types/iterator.generated.go", data)
}

func generateT(tmpl *template.Template, filePath string, data interface{}) {
	errorMsg := fmt.Sprintf("generate template %s to %s", tmpl.Name(), filePath) + ": %v"

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
}

var tmpl = template.Must(template.New("iterator").Parse(`
package types

import (
	"errors"
	"fmt"
)

{{- range $k, $v := . }}
/*
NextObjectFunc is the func used in iterator.

Notes
- ErrDone should be return while there are no items any more.
- Input objects slice should be set every time.
*/
type Next{{$k}}Func func(*{{$k}}Page) error

type {{$k}}Page struct {
	Token string
	Data  []{{$v}}
}

type {{$k}}Iterator struct {
	next Next{{$k}}Func

	index int
	done  bool

	o {{$k}}Page
}

func (it *{{$k}}Iterator) Next() (object {{$v}}, err error) {
	if it.index < len(it.o.Data) {
		it.index++
		return it.o.Data[it.index-1], nil
	}
	if it.done {
		return nil, IterateDone
	}

	// Reset buf before call next.
	it.o.Data = it.o.Data[:0]

	err = it.next(&it.o)
	if err == nil {
		it.index = 1
		return it.o.Data[0], nil
	}
	if !errors.Is(err, IterateDone) {
		return nil, fmt.Errorf("iterator next failed: %w", err)
	}

	// Mark this iterator has been done, no more elem will be fetched.
	it.done = true
	if len(it.o.Data) == 0 {
		return nil, IterateDone
	}

	it.index = 1
	return it.o.Data[0], nil
}
{{- end }}
`))
