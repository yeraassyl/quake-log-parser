package main

import (
	"bufio"
	"os"
	"strings"
)

// I assume that all logs are formatted correctly

// All I care about is if all the matches are initialized and shut down correctly
// Because if not, it will mess with the consistency and correctness of match reports.

func CheckConsistency(filePath string) (bool, int, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return false, 0, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	initGame := false
	exit := false

	cnt := 1
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "InitGame") {
			if initGame {
				return false, cnt, nil
			}
			initGame = true
		} else if strings.Contains(line, "Exit") {
			if !initGame {
				return false, cnt, nil
			} else {
				initGame = false
			}
			exit = true
		} else if strings.Contains(line, "ShutdownGame") {
			if exit {
				exit = false
				initGame = false
				cnt++
				continue
			}
			if !initGame {
				return false, cnt, nil
			} else {
				initGame = false
			}
		}
		cnt++
	}

	if initGame {
		return false, cnt, nil
	}

	return true, cnt, nil
}
