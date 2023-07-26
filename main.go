package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Player struct {
	Name       string
	DeathCount int
	KillCount  int
}

type Game struct {
	Settings   map[string]string
	Players    map[string]Player
	Kills      int
	Mod        map[string]int
	Finished   bool
	ExitReason string
}

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

func main() {
	filePath := "qgames.log"
	isConsistent, err := checkConsistency(filePath)

	if err != nil || !isConsistent {
		log.Fatal(err)
	}

	file, _ := os.Open(filePath)

	defer file.Close()

	var games []Game
	var game *Game
	var players map[string]*Player

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		command := parts[1][:len(parts[1])-1] // Remove the ":" at the end

		if command == "InitGame" {
			// If there is a current game, add it to the list of games
			if game != nil {
				games = append(games, *game)
			}
			// Start a new game and reset the players
			game = &Game{Settings: make(map[string]string), Players: make(map[string]Player), Mod: make(map[string]int)}
			// Keep track of players with their ID
			players = make(map[string]*Player)
			settings := strings.Split(line[strings.Index(line, "\\")+1:], "\\")
			for i := 0; i < len(settings); i += 2 {
				game.Settings[settings[i]] = settings[i+1]
			}

		} else if game != nil && command == "ClientConnect" {
			playerID := parts[2]
			players[playerID] = &Player{Name: "", DeathCount: 0, KillCount: 0}

		} else if game != nil && command == "ClientUserinfoChanged" {
			playerID := parts[2]
			info := strings.Join(parts[3:], " ")
			re := regexp.MustCompile(`n\\(.+?)\\t`)
			match := re.FindStringSubmatch(info)
			if len(match) > 1 {
				name := match[1]
				players[playerID].Name = name
			}

		} else if game != nil && command == "Kill" {
			killerID, killedID, modID := parts[2], parts[3], parts[4]
			game.Kills++
			game.Mod[modID]++

			// Increase killcount for killer if it's not <world>
			if killerID != "1022" {
				players[killedID].KillCount++
			}

			// Increase deathcount for killed player
			players[killedID].DeathCount++
		} else if game != nil && command == "ClientDisconnect" {
			playerID := parts[2]
			
			// Remove disconnected player from players but add his info to the game report
			game.Players[players[playerID].Name] = *players[playerID]
			delete(players, playerID)

		} else if game != nil && command == "ShutdownGame" {
			for _, player := range players {
				game.Players[player.Name] = *player
			}
			games = append(games, *game)
			game = nil
			players = nil

		} else if game != nil && command == "Exit" {
			game.Finished = true
			game.ExitReason = strings.Join(parts[2:], " ")
		}
	}

	// If there's an ongoing game at the end of the logs, add it to the games
	if game != nil {
		games = append(games, *game)
	}

	//Print the games
	for _, game := range games {
		fmt.Println(game)
	}
}

// 1. Read the log file
// 2. Group the game data of each match
// 3. Collect kill data
