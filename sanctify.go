package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/as/edit"
	"github.com/as/text"
	"github.com/golang/lint"
)

var (
	omit       = flag.Bool("o", false, "add omitempty flag to json tag")
	typ        = flag.String("t", "X", "the root element's type name")
	pkg        = flag.String("p", "main", "the name of the package")
	omitstring = ""
	level      int
	pass       int
	b          = new(bytes.Buffer)
	rules      []*edit.Command
	pred       string
)

func Name(s string) string {
	b := []byte(s)
	if rules != nil {
		for _, v := range rules {
			b = transform(v, b)
		}
		s = string(b)

	}
	if pass > 0 && len(s) > 0 {
		b[0] = b[0] &^ 0x20
		s = string(b)
		if strings.HasPrefix(s, pred) && len(s) != len(pred) {
			// Ble, ble, bleh, bleh, ble, that's all folks!
			s = s[len(pred):]
		}
	}
	Printf(s)
	return s
}
func Printf(fm string, i ...interface{}) {
	for i := 0; i < level; i++ {
		fmt.Fprintf(b, "\t")
	}
	fmt.Fprintf(b, fm, i...)
}
func Parse(j interface{}) {
	fmt.Fprintf(b, "package %s\ntype %s ", *pkg, *typ)
	parse(j)
}

func addtag(s string) string {
	return fmt.Sprintf("`json:\"%s%s\"`\n", s, omitstring)
}

func parse(j interface{}) {
	v := reflect.ValueOf(j).Interface()
	switch reflect.TypeOf(j).Kind() {
	case reflect.Ptr:
		Printf(" peter")
	case reflect.Map:
		Printf(" struct{\n")
		level++
		for k, v := range v.(map[string]interface{}) {
			last := pred
			pred = Name(k)
			parse(v)
			pred = last
			Printf(addtag(k))
		}
		level--
		Printf(" }")
	case reflect.Slice:
		fmt.Fprint(b, " []")
		for _, v := range v.([]interface{}) {
			parse(v)
		}
	case reflect.Int, reflect.Float64:
		fmt.Fprint(b, " int")
	case reflect.String:
		fmt.Fprint(b, " string")
	case reflect.Bool:
		fmt.Fprint(b, " bool")
	}

}

func ck(s string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", s, err)
	}
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	ck("stdin", err)

	var mp interface{}
	err = json.Unmarshal(data, &mp)
	Parse(mp)
	// if golint complains about names, change them to the suggested name
	files := map[string][]byte{
		"main.go": b.Bytes(),
	}
	linter := &lint.Linter{}
	prob, err := linter.LintFiles(files)
	ck("lint", err)

	for _, v := range prob {
		if v.Category != "naming" {
			continue
		}
		if err != nil {
			log.Fatalln(err)
		}
		rules = append(rules, edit.MustCompile(mkX([]byte(v.Text))))
		ck("rules", err)
	}

	b.Truncate(0)
	pass++
	Parse(mp)

	// step 3
	// gofmt the resulting go source
	cmd := exec.Command("gofmt", "-s")
	cmd.Stdin = b
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	ck("gofmt", err)
}

func transform(c *edit.Command, s []byte) []byte {
	ed, _ := text.Open(text.BufferFrom(s))
	err := c.Run(ed)
	if err != nil {
		panic(err)
	}
	return ed.Bytes()
}

// mkX creates an edit program from the result of a go-lint name suggestion. the output is
// compiled to an edit program and executed on the source code
func mkX(data []byte) string {
	buf := text.BufferFrom(data)
	ed, _ := text.Open(buf)
	cc := edit.MustCompile
	for _, v := range []*edit.Command{cc(`,x, should be ,c,@,`), cc(`,y,[ \t],v,@,d`), cc(`,x,[ \t\n\r@]+,c, ,`)} {
		v.Run(ed)
	}
	xy := strings.Fields(string(buf.Bytes()))
	if len(xy) < 2 {
		panic("cant find replacement expression for  name")
	}
	return fmt.Sprintf(",x,%s,c,%s,", xy[0], xy[1])
}

func init() {
	flag.Usage = func() {
		fmt.Sprint(`
NAME
	sanctify - convert json to an idiomatic go struct

SYNOPSIS
	sanctify < data.json

DESCRIPTION
	Sanctify reads stdin and converts JSON to an idiomatic Go struct.
	It applies the following rules:
	
	1.) Marshal JSON into a Go struct
	2.) Remove underscores in variable names
	3.) Capitalize letter occupying position of deleted underscores
	4.) Run golint to enumerate improperly-punctuated acronyms
	5.) Run gofmt -s to simplify code
	6.) Remove stuttercase naming from nested structures (see example)
	
	There are a few options:
	
	-p    Package name to generate (default: main)
	-t    JSON root element type (default: X)
	-o    Add omitempty to all fields

EXAMPLES
    echo {"msg":{"msg_string":"hi","msg_num": 3}} | sanctify

        package main

        type X struct {
            Msg struct {
            String string ` + "`" + `json:"msg_string"` + "`" + `
            Num    int    ` + "`" + `json:"msg_num"` + "`" + `
            } ` + "`" + `json:"msg"` + "`" + `
        }

	`)
	}
	flag.Parse()
	if *omit {
		omitstring = ",omitempty"
	}
}
