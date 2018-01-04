package main

import (
  "github.com/aws/aws-sdk-go/aws/session"
  "fmt"
  "github.com/tehmoon/errors"
  "os"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
)

const (
  PREFIXES = "/0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!-_.*'()"
  OUTPUT_LINE Output = iota
  OUTPUT_CSV
  OUTPUT_JSON_LINE
)

type Output uint8

func main() {
  options, err := ParseFlags()
  if err != nil {
    fmt.Fprintln(os.Stderr, errors.Wrap(err, "Error parsing flags").Error())
    os.Exit(2)
  }

  sess := session.Must(session.NewSessionWithOptions(session.Options{
      SharedConfigState: session.SharedConfigEnable,
  }))

  svc := s3.New(sess)
  sync := make(chan *Directory)

  for _, prefix := range PREFIXES {
    p := fmt.Sprintf("%s%s", options.Prefix, string(prefix))
    go fetchPrefix(svc, options.Bucket, p, sync)
  }

  root := NewRootDirectory()

  for range PREFIXES {
    MergeDirectories(root, <- sync)
  }

  err = OutputTree(root, options.Depth, options.OutputFormat)
  if err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
  }
}

func fetchPrefix(svc *s3.S3, bucket, prefix string, syncMaster chan *Directory) {
  var (
    pages uint64 = 0
    sync = make(chan *Directory)
    objectVersionsInput = &s3.ListObjectVersionsInput{
      Bucket: aws.String(bucket),
      Prefix: aws.String(prefix),
    }
  )

  err := svc.ListObjectVersionsPages(objectVersionsInput, func (page *s3.ListObjectVersionsOutput, lastPage bool) (bool) {
    pages++

    go buildDirectory(page.Versions, sync)

    return ! lastPage
  })
  if err != nil {
    fmt.Printf("FATAL Error: %s\n", err.Error())
    os.Exit(1)
  }

  root := NewRootDirectory()

  for i := uint64(0); i < pages; i++ {
    MergeDirectories(root, <- sync)
  }

  syncMaster <- root
}

func buildDirectory(versions []*s3.ObjectVersion, sync chan *Directory) {
  root := NewRootDirectory()

  for _, version := range versions {
    key := *version.Key

    if *version.IsLatest && key[len(key) - 1] != '/' {
      root.CreateFullPathFile(*version.Key, *version.Size)
    }
  }

  sync <- root
}
