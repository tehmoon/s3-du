package main

import (
  "path"
  "os"
  "github.com/tehmoon/errors"
)

func inspectDepth(d *Directory, depth uint64, root string, send chan *Directory) {
  root = path.Join(root, d.Name)

  if len(d.Children) == 0 || depth == 0 {
    send <- d
    return
  }

  depth = depth - 1
  for _, child := range d.Children {
    inspectDepth(child, depth, root, send)
  }
}

func OutputTree(d *Directory, depth uint64, ot Output) (error) {
  send := make(chan *Directory)
  stop := make(chan struct{})
  syncBack := make(chan error)

  printer, err := NewPrinter(os.Stdout, ot)
  if err != nil {
    return errors.Wrap(err, "Error calling NewPrinter()")
  }

  go func () {
    var err error

    LOOP: for {
      select {
        case d := <- send:
          err = printer.Print(d)
          if err != nil {
            err = errors.Wrap(err, "Error printing output")
            break LOOP
          }
        case <- stop:
          break LOOP
      }
    }

    printer.Close()
    syncBack <- err
  }()

  go func() {
    inspectDepth(d, depth, "", send)
    stop <- struct{}{}
  }()

  return <- syncBack
}
