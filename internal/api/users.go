package api

import (
	"fmt"
	"net/url"

	"github.com/programmersd21/mint/internal/models"
)

func (c *Client) GetUser(idOrUsername string) (*models.User, error) {
	path := fmt.Sprintf("/user/%s", url.PathEscape(idOrUsername))
	var result models.User
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &result, nil
}

func (c *Client) GetUserProjects(userID string) ([]models.Project, error) {
	path := fmt.Sprintf("/user/%s/projects", url.PathEscape(userID))
	var result []models.Project
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get user projects: %w", err)
	}
	return result, nil
}

func (c *Client) GetUserNotifications(userID string) ([]models.Notification, error) {
	path := fmt.Sprintf("/user/%s/notifications", url.PathEscape(userID))
	var result []models.Notification
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get user notifications: %w", err)
	}
	return result, nil
}
