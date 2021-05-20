package main

import (
  "fmt"
  "io"
  "os"

  "github.com/arglucas/hello-world/hello"
  "github.com/arglucas/hello-world/world"
)

func displayGreetings(w io.Writer) {
  fmt.Fprintln(w, hello.Greet())
  fmt.Fprintln(w, world.Greet())
}

func main() {
  //fmt.Println("This is the main package")
  //fmt.Println(hello.Greet())
  //fmt.Println(world.Greet())
  displayGreetings(os.Stdout)
}
