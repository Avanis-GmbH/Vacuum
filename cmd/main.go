package main

import "fmt"

const VERSION = "0.0.1"

func main() {

	showStartupBanner()

	fmt.Println("Hewo UwU")
}

func showStartupBanner() {
	fmt.Println("=============================================================")
	fmt.Printf("Go Dust Vacuum v%v - Created by Simon Nils Rach \n", VERSION)
	fmt.Printf("=============================================================\n\n")
}
