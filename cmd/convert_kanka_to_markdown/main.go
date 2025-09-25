package main

import (
	"fmt"
	"log"
	"os"

	"errenstar/internal/embeddings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/convert_kanka_to_markdown/main.go <directory_path>")
		fmt.Println("Example: go run cmd/convert_kanka_to_markdown/main.go ./raw_content/characters")
		os.Exit(1)
	}

	path := os.Args[1]
	
	fmt.Printf("Converting Kanka JSON files to Markdown in directory: %s\n", path)
	
	err := embeddings.ConvertDirectoriesToMarkdown(path)
	if err != nil {
		log.Fatalf("Error converting directories: %v", err)
	}
	
	fmt.Println("Conversion completed successfully!")
}
