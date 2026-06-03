package api

import (
	"fmt"
	"net/url"

	"github.com/programmersd21/mint/internal/models"
)

func (c *Client) GetVersion(versionID string) (*models.Version, error) {
	path := fmt.Sprintf("/version/%s", url.PathEscape(versionID))
	var result models.Version
	if err := c.get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("get version: %w", err)
	}
	return &result, nil
}

func (c *Client) GetVersions(versionIDs []string) ([]models.Version, error) {
	query := url.Values{}
	ids := "["
	for i, id := range versionIDs {
		if i > 0 {
			ids += ","
		}
		ids += fmt.Sprintf("\"%s\"", id)
	}
	ids += "]"
	query.Set("ids", ids)

	var result []models.Version
	if err := c.get("/versions", query, &result); err != nil {
		return nil, fmt.Errorf("get versions: %w", err)
	}
	return result, nil
}

func (c *Client) GetVersionDependencies(projectID string) (*models.Version, []models.Version, []models.Project, error) {
	path := fmt.Sprintf("/project/%s/dependencies", url.PathEscape(projectID))
	var result struct {
		Version  models.Version   `json:"version"`
		Versions []models.Version `json:"versions"`
		Projects []models.Project `json:"projects"`
	}
	if err := c.get(path, nil, &result); err != nil {
		return nil, nil, nil, fmt.Errorf("get version dependencies: %w", err)
	}
	return &result.Version, result.Versions, result.Projects, nil
}
