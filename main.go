package main

import (
  "bufio"
  "fmt"
  "os"
  "github.com/VanshikaaGuptaa/loghub/internal/aggregator"
)

func main() {
  f, _ := os.Open("sample.log")
  defer f.Close()
  sc := bufio.NewScanner(f)
  for sc.Scan() {
    e, _ := aggregator.ParseLine(sc.Text())
    fmt.Printf("%+v\n", e) // prints struct with all fields
  }
}
