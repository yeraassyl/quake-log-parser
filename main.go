package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "qgames.log"
	isConsistent, line, err := CheckConsistency(filePath)

	if err != nil {
		log.Fatal(err)
	}

	if !isConsistent {
		log.Fatalf("Log is not consistent at line: %d", line)
	}

	games, err := ParseLogs(filePath)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(
		"output.json",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(games)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Report generated")
}

// 1. Read the log file
// 2. Group the game data of each match
// 3. Collect kill data
