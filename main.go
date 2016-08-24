package main

import (
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    "github.com/jamesford/slack"
)

// Load config file
var config = ReadConfig()

// Respond to a ping request
func pingHandler(c echo.Context) error {
    return c.String(http.StatusOK, "pong")
}

// Redirect requester to Slack add app page
func addHandler(c echo.Context) error {
    return c.Redirect(http.StatusFound,
        "https://slack.com/oauth/authorize?scope=" + config.Scope + "&client_id=" + config.ClientID)
}

// Get Oauth token for team adding app
func authHandler(c echo.Context) error {
    code := c.QueryParam("code")

    // Retreive Slack AccessToken for Team
    r, err := slack.GetOAuthResponse(config.ClientID,
        config.ClientSecret,
        code,
        config.RedirectURI,
        false)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Construct team to save
    team := Team{ID: r.TeamID, AccessToken: r.AccessToken, Scope: r.Scope}

    // Save Team to redis
    err = saveTeam(team)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Oauth success, redirect user to slack app management
    return c.Redirect(http.StatusFound, config.SuccessRedirect)
}

// Handle a new /rps command from slack
func challengeHandler(c echo.Context) error {
    // Bind message form to Message struct
    m := new(Message)
    if err := c.Bind(m); err != nil {
        return c.NoContent(http.StatusBadRequest)
    }

    // The token from Slack should equal the verificationToken
    // If not equal request could come from outside of slack
    if m.Token != config.VerifyToken {
        return c.NoContent(http.StatusBadRequest)
    }

    // Token Verified, send OK response
    c.NoContent(http.StatusOK)

    // Format Challengee string
    s := strings.TrimSpace(m.Text)
    s = strings.TrimPrefix(s, "@")

    // Construct game to save
    game := Game{ID: newID(), Done: false, Challenger: Player{Name: m.UserName, Done: false}, Challengee: Player{Name: s, Done: false}}

    // Get Challenger's Team from Redis
    team, err := getTeam(m.TeamID)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Init slack client for team
    client := slack.New(team.AccessToken)

    // Send Challengee message
    challengeeMessage := newInteractiveMessage("gee:"+game.ID, game.Challenger.Name+" has challenged you!")
    _, _, err = client.PostMessage("@"+game.Challengee.Name, challengeeMessage.Text, challengeeMessage)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }

    // Send Challenger message
    challengerMessage := newInteractiveMessage("ger:"+game.ID, "Your challenge against "+ game.Challengee.Name +" has begun!")
    _, _, err = client.PostMessage("@"+game.Challenger.Name, challengerMessage.Text, challengerMessage)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }

    // Save Game to redis
    err = saveGame(game)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }
    return nil
}

// Handle a button interaction from slack
func interactionHandler(c echo.Context) error {

    // Transform json payload into Go struct (InteractionMessage)
    payload := c.FormValue("payload")
    interactionMessage := new(InteractionMessage)
    json.Unmarshal([]byte(payload), &interactionMessage)

    // The token from Slack should equal the verificationToken
    // If not equal request could come from outside of slack
    if interactionMessage.Token != config.VerifyToken {
        return c.NoContent(http.StatusBadRequest)
    }

    // Determinte Game ID and Challenger or Challengee's response
    split := strings.Split(interactionMessage.CallbackID, ":")
    playerType := split[0]
    gameID := split[1]

    // Load game from redis
    game, err := getGame(gameID)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Update the Challeger or Challengee fields
    var currentIsChallenger bool
    var currentPlayer *Player
    var opponentPlayer *Player

    if playerType == "ger" {
        currentIsChallenger = true
        currentPlayer = &game.Challenger
        opponentPlayer = &game.Challengee
    }

    if playerType == "gee" {
        currentIsChallenger = false
        currentPlayer = &game.Challengee
        opponentPlayer = &game.Challenger
    }

    currentPlayer.Move = interactionMessage.Actions[0].Value
    currentPlayer.TS = interactionMessage.MessageTs
    currentPlayer.Channel = interactionMessage.Channel.ID
    currentPlayer.Done = true

    // If both players done, mark game done
    if game.Challenger.Done && game.Challengee.Done {
        game.Done = true
    }

    // Save the updated game & players to redis
    err = saveGame(game)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Get team from redis
    team, err := getTeam(interactionMessage.Team.ID)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }

    // Init slack client for team
    client := slack.New(team.AccessToken)

    // Check if game is done
    if game.Done {
        // Update both messages with a result

        // Determine winner
        currentIsWinner, currentIsLoser, tieGame := resultForCurrent(currentPlayer.Move, opponentPlayer.Move)

        // Prep Both Message Texts
        opponentMessageText := challengeMessageText(false, currentPlayer.Name)
        currentMessageText := challengeMessageText(true, opponentPlayer.Name)

        // Handle tie game
        if tieGame {
            opponentMessage := newTieMessage(opponentPlayer.Move, currentPlayer.Move, opponentMessageText)
            _, _, _, err = client.UpdateMessage(opponentPlayer.Channel, opponentPlayer.TS, opponentMessageText, slack.UpdateMessageParameters{Attachments: opponentMessage.Attachments})
            if err != nil {
                fmt.Println(err.Error())
                return c.String(http.StatusInternalServerError, err.Error())
            }

            currentMessage := newTieMessage(currentPlayer.Move, opponentPlayer.Move, currentMessageText)
            cM, _ := json.Marshal(currentMessage)
            return c.JSONBlob(http.StatusOK, []byte(cM))
        }

        // Handle currentPlayer wins
        if currentIsWinner {
            opponentMessage := newLoserMessage(opponentPlayer.Move, currentPlayer.Move, opponentMessageText)
            _, _, _, err = client.UpdateMessage(opponentPlayer.Channel, opponentPlayer.TS, opponentMessageText, slack.UpdateMessageParameters{Attachments: opponentMessage.Attachments})
            if err != nil {
                fmt.Println(err.Error())
                return c.String(http.StatusInternalServerError, err.Error())
            }

            currentMessage := newWinnerMessage(currentPlayer.Move, opponentPlayer.Move, currentMessageText)
            cM, _ := json.Marshal(currentMessage)
            return c.JSONBlob(http.StatusOK, []byte(cM))
        }

        // Handle opponentPlayer wins
        if currentIsLoser {
            opponentMessage := newWinnerMessage(opponentPlayer.Move, currentPlayer.Move, opponentMessageText)
            _, _, _, err = client.UpdateMessage(opponentPlayer.Channel, opponentPlayer.TS, opponentMessageText, slack.UpdateMessageParameters{Attachments: opponentMessage.Attachments})
            if err != nil {
                fmt.Println(err.Error())
                return c.String(http.StatusInternalServerError, err.Error())
            }

            currentMessage := newLoserMessage(currentPlayer.Move, opponentPlayer.Move, currentMessageText)
            cM, _ := json.Marshal(currentMessage)
            return c.JSONBlob(http.StatusOK, []byte(cM))
        }
    }

    // Send 200 OK
    c.NoContent(http.StatusOK)

    // Game is not done, send waiting message to first interactor
    waitingMessageText := challengeMessageText(currentIsChallenger, opponentPlayer.Name)
    waitingMessage := newWaitingMessage(currentPlayer.Move, waitingMessageText)

    // Send waiting message via slack api
    channel, timestamp, _, err := client.UpdateMessage(currentPlayer.Channel, currentPlayer.TS, waitingMessageText, slack.UpdateMessageParameters{Attachments: waitingMessage.Attachments})
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }

    currentPlayer.Channel = channel
    currentPlayer.TS = timestamp

    // Save the updated game & players to redis
    err = saveGame(game)
    if err != nil {
        fmt.Println(err.Error())
        return nil
    }
    return nil
}

func main() {
    // Create an Echo
    e := echo.New()

    // Register Echo Routes
    e.GET("/ping", pingHandler)

    e.GET("/add", addHandler)

    e.GET("/auth", authHandler)

    e.POST("/challenge", challengeHandler)

    e.POST("/interaction", interactionHandler)

    // Run the Echo
    e.Run(standard.New(config.Port))
}