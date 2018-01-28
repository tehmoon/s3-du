package main

import (
  "path"
  "os"
  "github.com/tehmoon/errors"
  "text/template"
  "encoding/json"
)

func inspectDepth(d *Directory, depth uint64, root string, send chan *DirectoryAttr) {
  root = path.Join(root, d.Name)

  if len(d.Children) == 0 || depth == 0 {
    send <- &d.Attr
    return
  }

  depth = depth - 1
  for _, child := range d.Children {
    inspectDepth(child, depth, root, send)
  }
}

func OutputTree(d *Directory, depth uint64, tmpl *template.Template) (error) {
  send := make(chan *DirectoryAttr)
  stop := make(chan struct{})
  syncBack := make(chan error)

  go func () {
    var err error

    LOOP: for {
      select {
        case d := <- send:
          err = tmpl.Execute(os.Stdout, d)
          if err != nil {
            err = errors.Wrap(err, "Error templating the output")
            break LOOP
          }
        case <- stop:
          break LOOP
      }
    }

    syncBack <- err
  }()

  go func() {
    inspectDepth(d, depth, "", send)
    stop <- struct{}{}
  }()

  return <- syncBack
}

var (
  functionTemplates = template.FuncMap{
    "json": func(d interface{}) (string) {
      payload, err := json.Marshal(d)
      if err != nil {
        return ""
      }

      return string(payload[:])
    },
    "json_indent": func(d interface{}) (string) {
      payload, err := json.MarshalIndent(d, "", "  ")
      if err != nil {
        return ""
      }

      return string(payload[:])
    },
  }
)
