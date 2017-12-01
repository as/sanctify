## Sanctify

Convert JSON to an idiomatic go struct
  
## Synopsis

`go get -u github.com/as/sanctify/...`

`sanctify < data.json`
  
## Description

Sanctify reads stdin and converts JSON to an idiomatic Go struct.
It applies the following rules:
	
- Marshal JSON into a Go struct
- Remove underscores in variable names
- Capitalize letter occupying position of deleted underscores
- Run golint to enumerate improperly-punctuated acronyms
- Run gofmt -s to simplify code
- Remove stuttercase naming from nested structures (see example)
	
## Options
	
	-p    Package name to generate (default: main)
	-t    JSON root element type (default: X)
	-o    Add omitempty to all fields
  
 ## Example
 
`echo {"msg":{"msg_string":"hi","msg_num": 3}} | sanctify`
   
```
package main
type X struct {
  Msg struct {
    String string `json:"msg_string"`
    Num    int    `json:"msg_num"`
  } `json:"msg"`
}
```
