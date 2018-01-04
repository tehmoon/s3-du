package main

import (
  "fmt"
  "encoding/csv"
  "io"
  "encoding/json"
  "github.com/tehmoon/errors"
  "time"
  "strconv"
)

type Printer interface {
  Init() (error)
  Print(*Directory) (error)
  Name() (string)
  Close() (error)
}

func NewPrinter(writer io.Writer, output Output) (Printer, error) {
  var printer Printer

  switch output {
    case OUTPUT_LINE:
      printer = newLinePrinter(writer)
    case OUTPUT_JSON_LINE:
      printer = newJsonLinePrinter(writer)
    case OUTPUT_CSV:
      printer = newCsvPrinter(writer)
    default:
      return nil, errors.New("Unknown output type")
  }

  err := printer.Init()
  if err != nil {
    return nil, errors.Wrap(err, "Error setting up printer")
  }

  return printer, nil
}

func newLinePrinter(writer io.Writer) (Printer) {
  return &LinePrinter{
    Writer: writer,
    init: false,
    now: time.Now(),
  }
}

type LinePrinter struct {
  Writer io.Writer
  now time.Time
  init bool
}

func (p LinePrinter) Name() (string) {
  return "line"
}

func (p *LinePrinter) Init() (error) {
  p.init = true

  return nil
}

func (p LinePrinter) Print(d *Directory) (error) {
  if ! p.init {
    return errors.New("Init() has to be called first")
  }

  fmt.Fprintf(p.Writer, "Directory %s has size %d and %d files\n", d.Root, d.Size, d.Files)

  return nil
}

func (p LinePrinter) Close() (error) {
  return nil
}

func newJsonLinePrinter(writer io.Writer) (Printer) {
  return &JsonLinePrinter{
    Writer: writer,
    init: false,
    now: time.Now(),
  }
}

func (p *JsonLinePrinter) Init() (error) {
  p.init = true

  return nil
}

func (p JsonLinePrinter) Print(d *Directory) (error) {
  d.Now = p.now

  payload, err := json.Marshal(d)
  if err != nil {
    return errors.Wrap(err, "Error Marshaling to JSON_LINE")
  }

  payload = append(payload, '\n')
  _, err = p.Writer.Write(payload)
  if err != nil {
    return errors.Wrap(err, "Error writing JSON")
  }

  return nil
}

type JsonLinePrinter struct {
  now time.Time
  Writer io.Writer
  init bool
}

func (p JsonLinePrinter) Name() (string) {
  return "json"
}

func (p JsonLinePrinter) Close() (error) {
  return nil
}

func newCsvPrinter(writer io.Writer) (Printer) {
  return &CsvPrinter{
    Writer: csv.NewWriter(writer),
    init: false,
    now: strconv.FormatInt(time.Now().Unix(), 10),
  }
}

type CsvPrinter struct {
  Writer *csv.Writer
  init bool
  now string
}

func (p CsvPrinter) Name() (string) {
  return "csv"
}

func (p *CsvPrinter) Init() (error) {
  header := []string{
    "now",
    "path",
    "regular_files",
    "byte_size",
  }

  err := p.Writer.Write(header)
  if err != nil {
    return errors.Wrap(err, "Error writing CSV headers")
  }

  p.init = true

  return nil
}

func (p CsvPrinter) Print(d *Directory) (error) {
  data := []string{
    p.now,
    d.Root,
    strconv.FormatInt(d.Files, 10),
    strconv.FormatInt(d.Size, 10),
  }

  err := p.Writer.Write(data)
  if err != nil {
    return errors.Wrap(err, "Error writing data")
  }

  return nil
}

func (p CsvPrinter) Close() (error) {
  p.Writer.Flush()
  return nil
}
