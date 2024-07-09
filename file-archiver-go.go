package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	// Define command-line flags
	dbPath := flag.String("db", "file_search_results.db", "Path to the SQLite database")
	sourcePath := flag.String("source", "d:\\csdportal", "Source path prefix")
	archivePath := flag.String("archive", "d:\\archive", "Archive path")
	flag.Parse()

	// Open the SQLite database
	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query for orphaned files
	rows, err := db.Query("SELECT path FROM file_search_results WHERE is_orphaned = 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Process each orphaned file
	for rows.Next() {
		var filePath string
		err := rows.Scan(&filePath)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Check if the file path starts with the source path
		if strings.HasPrefix(filePath, *sourcePath) {
			// Calculate the new path in the archive
			relPath, err := filepath.Rel(*sourcePath, filePath)
			if err != nil {
				log.Println("Error getting relative path:", err)
				continue
			}
			newPath := filepath.Join(*archivePath, relPath)

			// Create the directory structure
			err = os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
			if err != nil {
				log.Println("Error creating directory:", err)
				continue
			}

			// Move the file
			err = os.Rename(filePath, newPath)
			if err != nil {
				log.Println("Error moving file:", err)
				continue
			}

			fmt.Printf("Moved %s to %s\n", filePath, newPath)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("File archiving process completed.")
}
