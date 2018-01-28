package main

import (
  "fmt"
  "flag"
  "github.com/tehmoon/errors"
  "text/template"
)

func ParseFlags() (*Options, error) {
  flags := &Flags{}

  flag.StringVar(&flags.Prefix, "p", "", "Prefix for s3 object keys")
  flag.StringVar(&flags.Bucket, "b", "", "Bucket to fetch keys from")
  flag.Uint64Var(&flags.Depth, "d", 0, "Calculate directory sizes with specified depth")
  flag.StringVar(&flags.Template, "template", `directory {{ .Root }} has size {{ .Size }} and {{ .Files }} files.`, "Go text/template to use when output. Use json or json_indent functions if you want")

  flag.Parse()

  if flags.Bucket == "" {
    return nil, errors.New("Option -b is mandatory")
  }

  if flags.Template == "" {
    return nil, errors.New("Option -template cannot be empty")
  }

  flags.Template = fmt.Sprintf("%s\n", flags.Template)

  tmpl, err := template.New("root").Funcs(functionTemplates).Parse(flags.Template)
  if err != nil {
    return nil, errors.Wrap(err, "Error parsing template")
  }

  options := &Options{
    Depth: flags.Depth,
    Prefix: flags.Prefix,
    Bucket: flags.Bucket,
    Human: flags.Human,
    Template: tmpl,
  }

  return options, nil
}

type Options struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
  Template *template.Template
}

type Flags struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
  Template string
}
