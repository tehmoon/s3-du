package main

import (
  "flag"
  "github.com/tehmoon/errors"
)

func ParseFlags() (*Options, error) {
  flags := &Flags{}

  flag.StringVar(&flags.Prefix, "p", "", "Prefix for s3 object keys")
  flag.StringVar(&flags.Bucket, "b", "", "Bucket to fetch keys from")
  flag.Uint64Var(&flags.Depth, "d", 0, "Calculate directory sizes with specified depth")

  flag.Parse()

  if flags.Bucket == "" {
    return nil, errors.New("Option -b is mandatory")
  }

  options := &Options{
    Depth: flags.Depth,
    Prefix: flags.Prefix,
    Bucket: flags.Bucket,
    Human: flags.Human,
  }

  return options, nil
}

type Options struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
}

type Flags struct {
  Depth uint64
  Prefix string
  Bucket string
  Human bool
}
