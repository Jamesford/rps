package main

// Team struct
type Team struct {
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
	Scope       string `json:"scope"`
}

// Player struct
type Player struct {
	Done bool         `json:"done"`
	Name string       `json:"name"`
	Move string       `json:"move"`
	TS string         `json:"ts"`
	Channel string     `json:"channel"`
}

// Game struct
type Game struct {
	ID         string `json:"id"`
	Done       bool   `json:"done"`
	Challenger Player `json:"challenger"`
	Challengee Player `json:"challengee"`
}

// Message struct
type Message struct {
	Token       string `form:"token"`
	TeamID      string `form:"team_id"`
	TeamDomain  string `form:"team_domain"`
	ChannelID   string `form:"channel_id"`
	ChannelName string `form:"channel_name"`
	UserID      string `form:"user_id"`
	UserName    string `form:"user_name"`
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
}

// InteractionMessage struct
type InteractionMessage struct {
	Actions []struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
	}                   `json:"actions"`
	CallbackID   string `json:"callback_id"`
	Team struct {
		ID       string `json:"id"`
		Domain   string `json:"domain"`
	}                   `json:"team"`
	Channel struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
	}                   `json:"channel"`
	User struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
	}                   `json:"user"`
	ActionTs     string `json:"action_ts"`
	MessageTs    string `json:"message_ts"`
	AttachmentID string `json:"attachment_id"`
	Token        string `json:"token"`
	ResponseURL  string `json:"response_url"`
}
