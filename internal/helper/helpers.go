package helper

import (
	"log"
	"os"
)

func ShowHelp() {
	log.Printf("Usage: %s {params}", os.Args[0])
	log.Println("      -c {config file}")
}
