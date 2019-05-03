package main

import (
  "fmt"

  "github.com/AlexisBRENON/fishbeats/engine"
  "github.com/AlexisBRENON/fishbeats/ui/gtk3"
)

func main() {
  fmt.Println("Hello fishbeats.")
  e := engine.NewEngine()
  gtk3.Main(e)
}
