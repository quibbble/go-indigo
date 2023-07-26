# Go-indigo

Go-indigo is a [Go](https://golang.org) implementation of the board game [Indigo](https://en.wikipedia.org/wiki/Indigo_(board_game)). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

## Usage

To play a game create a new Indigo instance:
```go
builder := Builder{}
game, err := builder.Create(&bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"}, // must contain at least 2 and at most 4 teams
    MoreOptions: TsuroMoreOptions{
        Seed: 123, // OPTIONAL - seed used to generate deterministic randomness which defaults to 0
        Variant: "Classic", // OPTIONAL - variants that change the game rules i.e. Classic (default), LargeHands
        RoundsUntilEnd: 10 - // OPTIONAL - the number of rounds played before the game ends
    }
})
```

To rotate a tile in your hand do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "RotateTile",
    MoreDetails: RotateTileActionDetails{
        Tile: "ABCDEF"
    },
})
```

To place a tile on the board do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "PlaceTile",
    MoreDetails: PlaceTileActionDetails{
        Row: 0,
        Column: 1,
        Tile: "ABCDEF"
    },
})
```

To get the current state of the game call the following:
```go
snapshot, err := game.GetSnapshot("TeamA")
```
