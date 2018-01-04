package main

import (
  "flag"
  "strings"
  "github.com/tehmoon/errors"
)

func ParseFlags() (*Options, error) {
  flags := &Flags{}

  flag.StringVar(&flags.Prefix, "p", "", "Prefix for s3 object keys")
  flag.StringVar(&flags.Bucket, "b", "", "Bucket to fetch keys from")
  flag.Uint64Var(&flags.Depth, "d", 0, "Calculate directory sizes with specified depth")
  flag.StringVar(&flags.OutputFormat, "f", "line", "Output format to use. One of: line, json_line or csv")

  flag.Parse()

  if flags.Bucket == "" {
    return nil, errors.New("Option -b is mandatory")
  }

  var outputFormat Output

  switch strings.ToLower(flags.OutputFormat) {
    case "line":
      outputFormat = OUTPUT_LINE
    case "json_line":
      outputFormat = OUTPUT_JSON_LINE
    case "csv":
      outputFormat = OUTPUT_CSV
    default:
      return nil, errors.Errorf("Unknown option -f format type: %s", flags.OutputFormat)
  }

  options := &Options{
    Depth: flags.Depth,
    Prefix: flags.Prefix,
    Bucket: flags.Bucket,
    Human: flags.Human,
    OutputFormat: outputFormat,
  }

  return options, nil
}

type Options struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
  OutputFormat Output
}

type Flags struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
  OutputFormat string
}
