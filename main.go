package main

import (
	"aery-graphql/db"
	"aery-graphql/server"

	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fatih/color"
)

func main() {
	/*
	 * Start the program
	 */
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	environment := flag.String("env", "development", "")

	fmt.Println("")
	color.New(color.FgWhite).Println("     ¸__¸  ___ ¸___     / /   ¸__¸ __¸¸ ¸__¸")
	color.New(color.FgWhite).Println("    /__¸//_¸  /¸__//___/ /   /__ //¸__//__¸ ")
	color.New(color.FgWhite).Println("   /   //____/   \\  /   /___/   //___/¸___/ ")
	color.New(color.FgBlue).Println("  ----------------------------------------------")
	fmt.Println()
	fmt.Print("⇨ environment is ")
	color.New(color.FgBlue).Print(*environment)
	fmt.Println()

	/**
	 * Connect to db
	 */
	db.Init()

	/*
	 * Start http server
	 */
	server.New()

	/*
	 * Stop the program gracefully
	 */
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill)
	<-channel
	fmt.Println("Stopping the program")
}
