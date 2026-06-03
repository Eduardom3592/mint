package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/programmersd21/mint/internal/models"
)

func (c *Client) GetProject(idOrSlug string) (*models.Project, error) {
	path := fmt.Sprintf("/project/%s", url.PathEscape(idOrSlug))
	var result models.Project
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	return &result, nil
}

func (c *Client) GetProjects(ids []string) ([]models.Project, error) {
	query := url.Values{}
	query.Set("ids", fmt.Sprintf("[%s]", strings.Join(quoteStrings(ids), ",")))

	var result []models.Project
	if err := c.get("/projects", query, &result); err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}
	return result, nil
}

func (c *Client) GetProjectVersions(projectID string, loaders []string, gameVersions []string) ([]models.Version, error) {
	path := fmt.Sprintf("/project/%s/version", url.PathEscape(projectID))
	query := url.Values{}

	if len(loaders) > 0 {
		query.Set("loaders", fmt.Sprintf("[%s]", strings.Join(quoteStrings(loaders), ",")))
	}
	if len(gameVersions) > 0 {
		query.Set("game_versions", fmt.Sprintf("[%s]", strings.Join(quoteStrings(gameVersions), ",")))
	}

	var result []models.Version
	if err := c.get(path, query, &result); err != nil {
		return nil, fmt.Errorf("get project versions: %w", err)
	}
	return result, nil
}

func (c *Client) GetProjectTeam(projectID string) ([]models.TeamMember, error) {
	path := fmt.Sprintf("/project/%s/members", url.PathEscape(projectID))
	var result []models.TeamMember
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get project team: %w", err)
	}
	return result, nil
}

func quoteStrings(strs []string) []string {
	out := make([]string, len(strs))
	for i, s := range strs {
		out[i] = quoteString(s)
	}
	return out
}

func quoteString(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
