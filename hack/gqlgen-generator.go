package main

import (
	"fmt"

	"github.com/99designs/gqlgen/cmd"
)

func main() {
	fmt.Println("\033[0;36m Generating Graphql assets... \033[0m")
	cmd.Execute()
	fmt.Println("\033[1;32m Done! \033[1;30m(copy or update all of the assets from /tmp folder to resolver package!)\033[0m")
}
