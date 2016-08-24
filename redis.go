package main

import (
    "encoding/json"
    "gopkg.in/redis.v4"
)

// Connect to redis
var rediscli = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

// Save team to redis with "team:" prefix
func saveTeam(team Team) error {
    t, _ := json.Marshal(team)
    err := rediscli.Set("team:" + team.ID, t, 0).Err()
    if err != nil {
        return err
    }
    return nil
}

// Save game to redis with "game:" prefix
func saveGame(game Game) error {
    g, _ := json.Marshal(game)
    err := rediscli.Set("game:" + game.ID, g, 0).Err()
    if err != nil {
        return err
    }
    return nil
}

// Get team from redis as Team struct
func getTeam(id string) (Team, error) {
    team := new(Team)
    t, err := rediscli.Get("team:" + id).Result()
    if err != nil {
        return *team, err
    }
    json.Unmarshal([]byte(t), &team)
    return *team, nil
}

// Get game from redis as Game struct
func getGame(id string) (Game, error) {
    game := new(Game)
    g, err := rediscli.Get("game:" + id).Result()
    if err != nil {
        return *game, err
    }
    json.Unmarshal([]byte(g), &game)
    return *game, nil
}