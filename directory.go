package main

import (
  "time"
  "strings"
  "path/filepath"
  "path"
)

type Directory struct {
  Parent *Directory `json:"-"`
  Children []*Directory `json:"-"`
  Root string `json:"path"`

  // Number of files incremented by CreateFullPathFile()
  Files int64 `json:"regular_files"`
  Size int64 `json:"byte_size"`
  Name string `json:"-"`
  Now time.Time `json:"now"`
}

// Recursively create directory for filepath. Filepath is obviously a file
// and not a directory.
func (d *Directory) CreateFullPathFile(path string, size int64) (*Directory) {
  directories := strings.Split(filepath.Dir(path), "/") /* separator is always / on s3 */
  cwd := d

  // If it is the root then update size and files
  if d.Parent == nil {
    d.Size += size
    d.Files++
  }

  for _, directory := range directories {
    var child *Directory

    for _, c := range cwd.Children {
      if directory == c.Name {
        c.Size += size
        c.Files++
        child = c
        break
      }
    }

    if child == nil {
      child = NewDirectory(cwd, directory)
      child.Size = size
      child.Files++
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
    Size: 0,
    Name: name,
    Parent: parent,
    Root: path.Join(parent.Root, name),
  }
}

func NewRootDirectory() (*Directory) {
  return &Directory{
    Children: make([]*Directory, 0),
    Size: 0,
    Name: "/",
    Root: "/",
  }
}

func MergeDirectories(to, from *Directory) {
  to.Size += from.Size
  to.Files += from.Files

  if len(from.Children) == 0 {
    return
  }

  for _, child := range from.Children {
    c := to.CreateCWDDirectory(child.Name)
    MergeDirectories(c, child)
  }
}
