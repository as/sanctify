## Sanctify

Convert JSON data into a Go struct with idiomatic field names
  
## Synopsis

`go get -u github.com/as/sanctify/...`

`sanctify < data.json`
  
## Description

Much of the work in creating a marshaller is in the tweaks that occur after the specification is
copied into Go source code. Because Go allows structs to have tags, you can keep the bad names in
an already-established API specification out of the struct field names. This package is an executable
that reads JSON data and converts it into Go, but also takes the labor and guess work out of such decisions. 

It uses the following process:
	
- Marshal JSON into a Go interface{}
- Recursively descend into arrays, amalagating fields of underlying JSON objects into a set
- Parse the amalagate tree, generating basic Go source in a main package
- Vet the package on the fly with golint, capturing naming suggestions in a buffer
- Remove underscores in variable names
- Capitalize letter occupying position of deleted underscores
- Compile rules to correct improperly-punctuated acronyms in struct field names
- Reparse the amalagate, applying corrections during the recursive descent step
- Prefix compare child fields to parent fields, remove stuttercase naming in Go field names
- Run gofmt -s to simplify code

Then you:

- Check the auto-generated code for consistency
- Paste it into your project

## Options
	
	-p    Package name to generate (default: main)
	-t    JSON root element type (default: X)
	-o    Add omitempty to all fields
  

## Why

- Generates Go without adding third party dependencies.
- Save time correcting huge datasets by hand

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

## Example 2

This hideous thing is straight from the json.org website. 

```
{"menu": {
    "header": "SVG Viewer",
    "items": [
        {"id": "Open"},
        {"id": "OpenNew", "label": "Open New"},
        null,
        {"id": "ZoomIn", "label": "Zoom In"},
        {"id": "ZoomOut", "label": "Zoom Out"},
        {"id": "OriginalView", "label": "Original View"},
        null,
        {"id": "Quality"},
        {"id": "Pause"},
        {"id": "Mute"},
        null,
        {"id": "Find", "label": "Find..."},
        {"id": "FindAgain", "label": "Find Again"},
        {"id": "Copy"},
        {"id": "CopyAgain", "label": "Copy Again"},
        {"id": "CopySVG", "label": "Copy SVG"},
        {"id": "ViewSVG", "label": "View SVG"},
        {"id": "ViewSource", "label": "View Source"},
        {"id": "SaveAs", "label": "Save As"},
        null,
        {"id": "Help"},
        {"id": "About", "label": "About Adobe CVG Viewer..."}
    ]
}}
```

``` cat horridmonkey.json | sanctify```

```
package main

type X struct {
        Menu struct {
                Header string `json:"header"`
                Items  []struct {
                        Id    string `json:"id"`
                        Label string `json:"label"`
                } `json:"items"`
        } `json:"menu"`
}
```


