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
