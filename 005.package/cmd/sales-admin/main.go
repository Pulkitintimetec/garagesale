package main

import (
	"flag"
	"log"

	"garagesale/005.package/internal/schema"
)

func main() {

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		schema.OpenDb()
		log.Println("Migrations complete")
		return

	case "seed":
		schema.UsingDb()
		log.Println("Seed data complete")
		return
	}

}
