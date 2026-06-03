package models

import "time"

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Email       *string   `json:"email"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	Joined      time.Time `json:"joined"`
	Role        string    `json:"role"`
	Badges      int       `json:"badges"`
	GitHubID    *int      `json:"github_id"`
}

type TeamMember struct {
	TeamID      string    `json:"team_id"`
	User        User      `json:"user"`
	Role        string    `json:"role"`
	Permissions int       `json:"permissions"`
	Accepted    bool      `json:"accepted"`
	Created     time.Time `json:"created"`
}

type Team struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Color       *int         `json:"color"`
	Members     []TeamMember `json:"members"`
}
