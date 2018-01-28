package main

import (
  "strings"
  "path/filepath"
  "path"
)

type Directory struct {
  Parent *Directory
  Children []*Directory
  Name string
  Attr DirectoryAttr
}

type DirectoryAttr struct {
  Root string `json:"path"`

  // Number of files incremented by CreateFullPathFile()
  Files int64 `json:"regular_files"`
  Size int64 `json:"byte_size"`
}

// Recursively create directory for filepath. Filepath is obviously a file
// and not a directory.
func (d *Directory) CreateFullPathFile(path string, size int64) (*Directory) {
  directories := strings.Split(filepath.Dir(path), "/") /* separator is always / on s3 */
  cwd := d

  // If it is the root then update size and files
  if d.Parent == nil {
    d.Attr.Size += size
    d.Attr.Files++
  }

  for _, directory := range directories {
    var child *Directory

    for _, c := range cwd.Children {
      if directory == c.Name {
        c.Attr.Size += size
        c.Attr.Files++
        child = c
        break
      }
    }

    if child == nil {
      child = NewDirectory(cwd, directory)
      child.Attr.Size = size
      child.Attr.Files++
      cwd.Children = append(cwd.Children, child)
    }

    cwd = child
  }

  return cwd
}

func (d *Directory) CreateCWDDirectory(name string) (*Directory) {
  for _, child := range d.Children {
    if child.Name == name {
      return child
    }
  }

  dir := NewDirectory(d, name)

  d.Children = append(d.Children, dir)

  return dir
}

func NewDirectory(parent *Directory, name string) (*Directory) {
  return &Directory{
    Children: make([]*Directory, 0),
    Attr: DirectoryAttr{
      Root: path.Join(parent.Attr.Root, name),
      Size: 0,
    },
    Name: name,
    Parent: parent,
  }
}

func NewRootDirectory() (*Directory) {
  return &Directory{
    Children: make([]*Directory, 0),
    Attr: DirectoryAttr{
      Root: "/",
      Size: 0,
    },
    Name: "/",
  }
}

func MergeDirectories(to, from *Directory) {
  to.Attr.Size += from.Attr.Size
  to.Attr.Files += from.Attr.Files

  if len(from.Children) == 0 {
    return
  }

  for _, child := range from.Children {
    c := to.CreateCWDDirectory(child.Name)
    MergeDirectories(c, child)
  }
}
