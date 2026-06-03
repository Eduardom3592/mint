package models

import "time"

type NotificationType string

const (
	NotificationTypeUpdate       NotificationType = "update"
	NotificationTypeFollow       NotificationType = "follow"
	NotificationTypeStatusChange NotificationType = "status_change"
	NotificationTypeModeration   NotificationType = "moderation"
)

type Notification struct {
	ID      string               `json:"id"`
	UserID  string               `json:"user_id"`
	Type    NotificationType     `json:"type"`
	Title   string               `json:"title"`
	Text    string               `json:"text"`
	Link    string               `json:"link"`
	ActorID *string              `json:"actor_id"`
	Actor   *User                `json:"actor"`
	Read    bool                 `json:"read"`
	Created time.Time            `json:"created"`
	Actions []NotificationAction `json:"actions"`
}

type NotificationAction struct {
	ActionType string `json:"action_type"`
	Label      string `json:"label"`
	URL        string `json:"url"`
}

type NotificationCount struct {
	Total int `json:"total"`
}
