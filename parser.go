package main

import (
	"bufio"
	"fmt"
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

var MOD map[string]string = map[string]string{
	"0":  "MOD_UNKNOWN",
	"1":  "MOD_SHOTGUN",
	"2":  "MOD_GAUNTLET",
	"3":  "MOD_MACHINEGUN",
	"4":  "MOD_GRENADE",
	"5":  "MOD_GRENADE_SPLASH",
	"6":  "MOD_ROCKET",
	"7":  "MOD_ROCKET_SPLASH",
	"8":  "MOD_PLASMA",
	"9":  "MOD_PLASMA_SPLASH",
	"10": "MOD_RAILGUN",
	"11": "MOD_LIGHTNING",
	"12": "MOD_BFG",
	"13": "MOD_BFG_SPLASH",
	"14": "MOD_WATER",
	"15": "MOD_SLIME",
	"16": "MOD_LAVA",
	"17": "MOD_CRUSH",
	"18": "MOD_TELEFRAG",
	"19": "MOD_FALLING",
	"20": "MOD_SUICIDE",
	"21": "MOD_TARGET_LASER",
	"22": "MOD_TRIGGER_HURT",
	"23": "MOD_GRAPPLE",
}

func ParseLogs(filePath string) ([]Game, error) {
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
			if game != nil {
				return nil, fmt.Errorf("Logfile is not consistent")
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
			game.Mod[MOD[modID]]++

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
			for _, player := range players {
				game.Players[player.Name] = *player
			}
			game.Finished = true
			game.ExitReason = strings.Join(parts[2:], " ")
			games = append(games, *game)
			game = nil
			players = nil
		}
	}

	return games, nil
}
