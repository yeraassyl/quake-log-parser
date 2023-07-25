package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"log"
)

// I assume that all logs are formatted correctly

// All I care about is if all the matches are initialized and shut down correctly
// Because if not, it will mess with the consistency and correctness of match reports.

func checkConsistency(filePath string) (bool, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return false, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	initGame := false
	shutdownGame := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "InitGame") {
			initGame = true
			shutdownGame = false
		} else if strings.Contains(line, "ShutdownGame") {
			shutdownGame = true
			if !initGame {
				return false, nil
			} else {
				initGame = false
				shutdownGame = false
			}
		}
	}

	if initGame && !shutdownGame {
		return false, nil
	}

	return true, nil
}

func main(){
	filePath := "qgames.log"
	isConsistent, err := checkConsistency(filePath)
	
	if err != nil || !isConsistent {
		log.Fatal(err)
	}

	


	fmt.Println("Hello")
}

// 1. Read the log file
// 2. Group the game data of each match
// 3. Collect kill data
