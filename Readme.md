# Quake log parser

Parse logs of the quake 3 arena server and group them by each match.

## Thought Process

To group information for each match, I need to find all the properties of the match, and then choose what properties are worth displaying/saving.

At first need to cleanup the log file, at line 97 and at 3948. Removed some matches that were abruptly ended.

We are looking for a properties of a match and users. 

Commands: `InitGame, ShutdownGame, ClientConnect, ClientBegin, ClientUserinfoChanged, ClientDisconnect, Exit, Kill, Item`

InitGame initializes the match, ShutdownGame shuts down the match, while Exit seems to indicate that the game was successfully ended and comes before ShutdownGame. Notice that, game can be shut down without Exiting or finishing.

Exit indicates that the match is concluded and shows final score, including ping, client info. Or red vs blue scores depending on the match type. The match type and other properties of the match are decided at the InitGame.  

All the players have ID, including \<world>. The logs with commands Kill, Item, ClientConnect, ClientBegin, ClientUserInfoChanged, ClientDisconnect display player ID after a command. 
I need to collect all the information about the match that can be gathered or aggregated from this logs. For example death_count, total_kills, player_score, match_type and so on.

## Structure 

The game report will be a json list of separate matches.

Match structure: 

```JSON
{
  "Settings": {
    "bot_minplayers": "0",
    "capturelimit": "8",
    "dmflags": "0"
  },
  "Players": {
    "Assasinu Credi": {
      "Name": "Assasinu Credi",
      "DeathCount": 30,
      "KillCount": 27
    },
    "Dono da Bola": {
      "Name": "Dono da Bola",
      "DeathCount": 19,
      "KillCount": 17
    }
  },
  "Kills": 131,
  "Mod": {
    "MOD_FALLING": 3,
    "MOD_MACHINEGUN": 4,
    "MOD_RAILGUN": 9
  },
  "Finished": true,
  "ExitReason": "Fraglimit hit."
}
```

`Finished` property indicates if a match successfully exited (`EXIT` command). `ExitReason` displays the conclusion of the match.

## Running the program

You can run this program just by using this command. 

```go run .```