package main

import (
	"context"
	"fmt"

	db "errenstar/internal/embeddings/db"
)

func main() {
	fmt.Println("Seeding embeddings database...")

	// Initialize the database
	embeddingsDB := db.InitializeDB()

	// Create context
	ctx := context.Background()

	// Seed the database
	embeddingsDB.SeedDB(ctx)

	fmt.Println("Database seeded successfully!")
}
