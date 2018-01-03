package main

import (
  "path"
  "os"
  "log"
)

var (
  logger = log.New(os.Stdout, "", log.LstdFlags)
)

func inspectDepth(d *Directory, depth uint64, root string) {
  root = path.Join(root, d.Name)

  if len(d.Children) == 0 || depth == 0 {
    logger.Printf("Directory %s has size %d and %d files\n", root, d.Size, d.Files)
    return
  }

  depth = depth - 1
  for _, child := range d.Children {
    inspectDepth(child, depth, root)
  }
}

func InspectDepth(d *Directory, depth uint64) {
  inspectDepth(d, depth, "")
}
