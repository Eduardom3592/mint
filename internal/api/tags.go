package api

import (
	"fmt"

	"github.com/programmersd21/mint/internal/models"
)

func (c *Client) GetCategories() ([]models.Category, error) {
	var result []models.Category
	if err := c.get("/tag/category", nil, &result); err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	return result, nil
}

func (c *Client) GetLoaders() ([]models.Loader, error) {
	var result []models.Loader
	if err := c.get("/tag/loader", nil, &result); err != nil {
		return nil, fmt.Errorf("get loaders: %w", err)
	}
	return result, nil
}

func (c *Client) GetGameVersions() ([]models.GameVersion, error) {
	var result []models.GameVersion
	if err := c.get("/tag/game_version", nil, &result); err != nil {
		return nil, fmt.Errorf("get game versions: %w", err)
	}
	return result, nil
}
