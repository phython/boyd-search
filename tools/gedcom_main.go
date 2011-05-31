
package main

import (
  "bytes"
  "flag"
  "./gedcom"
  "log"
  "os"
)

func main() {
  flag.Parse()
  for i := 0; i < flag.NArg(); i++ {
    buffer := new(bytes.Buffer)
    file, err := os.Open(flag.Arg(i))
    if err != nil {
      os.Exit(1)
    }
    bytes, read_err := buffer.ReadFrom(file)
    if bytes == 0 || read_err != nil {
      os.Exit(1)
    }
    log.Print("Parsing ", bytes, " bytes from ", flag.Arg(i))
    var data search.RawGedCom
    data.Parse(buffer)
  }
}
