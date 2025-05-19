package main

import (
	"dirsync/sync"
	"dirsync/validators"
	"flag"
	"log"
)

func main() {
	// Define flags for source, destination, and delete-missing option
	srcDir := flag.String("src", "", "Source folder path")
	dstDir := flag.String("dst", "", "Destination folder path")
	shouldDeleteMissing := flag.Bool(
		"delete-missing",
		false,
		"Delete files in destination that don't exist in source",
	)
	flag.Parse()

	if *srcDir == "" || *dstDir == "" {
		log.Fatal(
			"Source and destination folders are required. " +
				"Usage: program -src=<source_folder> -dst=<destination_folder> [-delete-missing]",
		)
	}
	if err := validators.ValidateSrcDir(*srcDir); err != nil {
		log.Fatalf("Error validating source directory: %v", err)
	}
	if err := sync.Dirs(*srcDir, *dstDir, *shouldDeleteMissing); err != nil {
		log.Fatalf("Error syncing directories: %v", err)
	}
	log.Println("Directory sync complete")
}
