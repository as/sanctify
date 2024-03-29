package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"

	"github.com/as/edit"
	"github.com/as/text"
	"golang.org/x/lint"
)

const Magic = "EX" // thanks microsoft

var (
	omit = flag.Bool("o", false, "add omitempty flag to json tag")
	typ  = flag.String("t", "X", "the root element's type name")
	pkg  = flag.String("p", "main", "the name of the package")

	omitstring = ""
	level      int
	pass       int
	b          = new(bytes.Buffer)
	rules      map[string]*edit.Command
	pred       string
)

const forbidden = `!@#$%^*()-=+{}/.`

func nodash(s string) string {
	b := []byte(s)
	if n := bytes.IndexAny(b, forbidden); n != -1 {
		if copy(b[n:], b[n+1:]) > 0 {
			b[n] = b[n] &^ 0x20
		}
		s = string(b[:len(b)-1])
	}
	return s
}

func Name(s string) string {
	s += Magic
	if rules != nil {
		if r, ok := rules[s]; ok {
			s = string(transform(r, []byte(s)))
		}
	}

	if pass > 0 && len(s) > 0 {
		b := []byte(s)
		b[0] = b[0] &^ 0x20
		s = string(b[:len(b)-len(Magic)])
		if strings.HasPrefix(s, pred) && len(s) != len(pred) {
			// Ble, ble, bleh, bleh, ble, that's all folks!
			s = s[len(pred):]
		}
	}
	Printf("%s", nodash(s))
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

func Unite(j interface{}) {
	unite(j)
}
func unite2(a []interface{}) []interface{} {
	if len(a) == 0 {
		return a
	}
	m3 := make(map[string]interface{})
	for _, v := range a {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			for k, v := range reflect.ValueOf(v).Interface().(map[string]interface{}) {
				m3[k] = v
			}
		}
	}
	a[0] = m3
	return a
}
func unite(j interface{}) {
	if j == nil {
		return
	}
	v := reflect.ValueOf(j).Interface()
	switch reflect.TypeOf(j).Kind() {
	case reflect.Map:
		for _, v := range v.(map[string]interface{}) {
			defer unite(v)
		}
	case reflect.Slice:
		defer unite2(v.([]interface{}))
	}
}

type bySuffix []string

func (b bySuffix) Suffix(i int) string {
	s := b[i]
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	t := strings.LastIndexAny(s[:len(s)-1], "-_")
	u := strings.LastIndexAny(s[:len(s)-1], "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if t >= u {
		if t <= 0 {
			return s
		}
		return s[t+1:]
	}
	return s[u:]
}

func (b bySuffix) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b bySuffix) Less(i, j int) bool {
	if v := strings.Compare(b.Suffix(i), b.Suffix(j)); v != 0 {
		return v == -1
	}
	return strings.Compare(b[i], b[j]) == -1
}
func (b bySuffix) Len() int { return len(b) }

func arrange(m map[string]interface{}) (om map[string]interface{}, keys []string) {
	m2 := make(map[string]interface{}, len(m))
	for k, v := range m {
		m2[k] = v
	}

	common := []string{"id", "name", "label", "tags"}
	for _, v := range common {
		if _, ok := m[v]; ok {
			delete(m2, v)
			keys = append(keys, v)
		}
	}

	exist := []string{}
	for k := range m2 {
		exist = append(exist, k)
	}

	sort.Sort(bySuffix(exist))
	return m, append(keys, exist...)
}

func parse(j interface{}) {
	if j == nil {
		Printf(" interface{}")
		return
	}
	v := reflect.ValueOf(j).Interface()
	switch reflect.TypeOf(j).Kind() {
	case reflect.Ptr:
		Printf(" peter")
	case reflect.Map:
		Printf(" struct{\n")
		level++
		m, keys := arrange(v.(map[string]interface{}))
		for _, k := range keys {
			v := m[k]
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
		i := 0
		for _, v := range v.([]interface{}) {
			parse(v)
			i++
			break // we don't need to parse them all, it's an array
		}
		if i == 0 {
			fmt.Fprint(b, "interface{}")
		}
	case reflect.Float64:
		_, f := math.Modf(v.(float64))
		if f == 0 {
			fmt.Fprint(b, " int")
		} else {
			fmt.Fprint(b, " float64")
		}
	case reflect.Int:
		fmt.Fprint(b, " int")
	case reflect.String:
		fmt.Fprint(b, " string")
	case reflect.Bool:
		fmt.Fprint(b, " bool")
	default:
		fmt.Fprint(b, " interface{}")
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
	ck("json", err)

	Unite(mp)
	Parse(mp)
	// if golint complains about names, change them to the suggested name
	files := map[string][]byte{
		"main.go": b.Bytes(),
	}
	linter := &lint.Linter{}
	prob, err := linter.LintFiles(files)
	if err != nil {
		log.Printf("ERROR DUMP: %s\n", b.Bytes())
		ck("lint", err)
	}

	rules = make(map[string]*edit.Command)
	for _, v := range prob {
		if v.Category != "naming" {
			continue
		}
		if err != nil {
			log.Fatalln(err)
		}
		src, prog := mkX([]byte(v.Text))
		rules[src] = edit.MustCompile(prog)
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
func mkX(data []byte) (string, string) {
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
	return xy[0], fmt.Sprintf(",x,%s,c,%s,", xy[0], xy[1])
}

func init() {
	flag.Usage = func() {
		fmt.Sprint(`
NAME
	sanctify - convert json to an idiomatic go struct

SYNOPSIS
	sanctify < data.json

OPTIONS
	
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
