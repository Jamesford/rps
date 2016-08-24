package main

import (
    "github.com/jamesford/slack"
    "github.com/ventu-io/go-shortid"
)

// Worker/Seed to ensure unique shortid's
var sid, err = shortid.New(1, shortid.DefaultABC, 4242)

// ID Helper
func newID() string {
    id, _ := sid.Generate()
    return id
}

// return true 1st if winner, 2nd if loser, 3rd if tie
func resultForCurrent(c string, o string) (bool, bool, bool) {
    // Tie
    if c == o { return false, false, true }

    // c wins
    if (c == "rock" && o == "scissors") || (c == "scissors" && o == "paper") || (c == "paper" && o == "rock") {
        return true, false, false
    }

    // c must have lost
    return false, true, false
}

// Converts "rock" to "Rock"
func choiceToText(choice string) string {
    var choiceText string;

    switch choice {
    case "rock":
        choiceText = "Rock"
    case "paper":
        choiceText = "Paper"
    case "scissors":
        choiceText = "Scissors"
    }

    return choiceText
}

// Converts "rock" to ":gem:"
func choiceToEmoji(choice string) string {
    var choiceEmoji string;

    switch choice {
    case "rock":
        choiceEmoji = ":gem:"
    case "paper":
        choiceEmoji = ":page_facing_up:"
    case "scissors":
        choiceEmoji = ":scissors:"
    }

    return choiceEmoji
}

// Challenger or Challengee's Message Text
func challengeMessageText(isChallenger bool, opponentName string) string {
    if isChallenger {
        return "Your challenge against " + opponentName + " has begun!"
    }
    return opponentName + " has challenged you!"
}

// Rock button
var rockAction = slack.AttachmentAction{
    Name: "rock",
    Text: ":gem: Rock",
    Type: "button",
    Value: "rock"}

// Paper button
var paperAction = slack.AttachmentAction{
    Name: "paper",
    Text: ":page_facing_up: Paper",
    Type: "button",
    Value: "paper"}

// Scissors button
var scissorsAction = slack.AttachmentAction{
    Name: "scissors",
    Text: ":scissors: Scissors",
    Type: "button",
    Value: "scissors"}

// Create attachment with callback id
func newAttachment(id string) slack.Attachment {
    attachment := slack.Attachment{
        CallbackID: id,
        Title: "Choose your weapon",
        Fallback: "Oh no! Something went wrong, battle not started :(",
        Color: "#3AA3E3",
        Actions: []slack.AttachmentAction{rockAction, paperAction, scissorsAction}}
    return attachment
}

// Create button message with callback id and given text
func newInteractiveMessage(id string, text string) slack.PostMessageParameters {
    attachment := newAttachment(id)

    message := slack.PostMessageParameters{
        Text: text,
        Attachments: []slack.Attachment{attachment}}

    return message
}

// Create waiting message with choice information
func newWaitingMessage(choice string, text string) slack.PostMessageParameters {
    attachment := slack.Attachment{
        Fallback: "Waiting for opponent...",
        Color: "#3AA3E3",
        Title: "Waiting for opponent...",
        Text: "You chose " + choiceToEmoji(choice) + " " + choiceToText(choice) + ", Good luck!"}

    message := slack.PostMessageParameters{
        Text: text,
        Attachments: []slack.Attachment{attachment}}

    return message
}

// Create the winner's message
func newWinnerMessage(playerMove string, opponentMove string, text string) slack.PostMessageParameters {
    attachment := slack.Attachment{
        Fallback: "You won the Rock Paper Scissors challenge!",
        Color: "#3AA3E3",
        Title: "You won! :smile:",
        Text: choiceToEmoji(playerMove) + " " + choiceToText(playerMove) + " vs " + choiceToText(opponentMove) + " " + choiceToEmoji(opponentMove)}

    message := slack.PostMessageParameters{
        Text: text,
        Attachments: []slack.Attachment{attachment}}

    return message
}

// Create the loser's message
func newLoserMessage(playerMove string, opponentMove string, text string) slack.PostMessageParameters {
    attachment := slack.Attachment{
        Fallback: "You lost the Rock Paper Scissors challenge!",
        Color: "#3AA3E3",
        Title: "You lost! :frowning:",
        Text: choiceToEmoji(playerMove) + " " + choiceToText(playerMove) + " vs " + choiceToText(opponentMove) + " " + choiceToEmoji(opponentMove)}

    message := slack.PostMessageParameters{
        Text: text,
        Attachments: []slack.Attachment{attachment}}

    return message
}

// Create tie game message
func newTieMessage(playerMove string, opponentMove string, text string) slack.PostMessageParameters {
    attachment := slack.Attachment{
        Fallback: "You tied with each other in the Rock Paper Scissors challenge!",
        Color: "#3AA3E3",
        Title: "You tied! :expressionless:",
        Text: choiceToEmoji(playerMove) + " " + choiceToText(playerMove) + " vs " + choiceToText(opponentMove) + " " + choiceToEmoji(opponentMove)}

    message := slack.PostMessageParameters{
        Text: text,
        Attachments: []slack.Attachment{attachment}}

    return message
}