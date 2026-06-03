package cache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/models"
	_ "modernc.org/sqlite"
)

type Cache struct {
	db *sql.DB
}

type Config struct {
	DataDir string
}

func Open(cfg Config) (*Cache, error) {
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(cfg.DataDir, "mint.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	c := &Cache{db: db}
	if err := c.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return c, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) DB() *sql.DB {
	return c.db
}

func (c *Cache) migrate() error {
	if _, err := c.db.Exec(schema); err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	var currentVer int
	_ = c.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&currentVer)

	if currentVer < 1 {
		_, _ = c.db.Exec(`INSERT OR IGNORE INTO schema_migrations (version) VALUES (1)`)
	}

	if currentVer < 2 {
		_, _ = c.db.Exec(`DROP TABLE IF EXISTS downloads`)
		_, _ = c.db.Exec(`DROP TABLE IF EXISTS installed_modpacks`)
		if _, err := c.db.Exec(`CREATE TABLE IF NOT EXISTS downloads (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_name TEXT NOT NULL DEFAULT '',
			project_id TEXT NOT NULL DEFAULT '',
			version_id TEXT NOT NULL DEFAULT '',
			version_number TEXT NOT NULL DEFAULT '',
			filename TEXT NOT NULL DEFAULT '',
			file_path TEXT NOT NULL DEFAULT '',
			url TEXT NOT NULL DEFAULT '',
			size INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'queued',
			is_mrpack INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			installed_at TEXT
		)`); err != nil {
			return fmt.Errorf("migrate v2 downloads: %w", err)
		}
		if _, err := c.db.Exec(`CREATE TABLE IF NOT EXISTS installed_modpacks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			version_id TEXT NOT NULL DEFAULT '',
			game_version TEXT NOT NULL DEFAULT '',
			loader TEXT NOT NULL DEFAULT '',
			path TEXT NOT NULL DEFAULT '',
			installed_at TEXT NOT NULL DEFAULT (datetime('now')),
			file_count INTEGER NOT NULL DEFAULT 0,
			verified INTEGER NOT NULL DEFAULT 0
		)`); err != nil {
			return fmt.Errorf("migrate v2 installed_modpacks: %w", err)
		}
		_, _ = c.db.Exec(`INSERT OR IGNORE INTO schema_migrations (version) VALUES (2)`)
	}

	return nil
}

func (c *Cache) CacheProject(project *models.Project) error {
	categories, _ := json.Marshal(project.Categories)
	loaders, _ := json.Marshal(project.Loaders)
	gameVersions, _ := json.Marshal(project.GameVersions)
	donationURLs, _ := json.Marshal(project.DonationURLs)
	gallery, _ := json.Marshal(project.Gallery)

	_, err := c.db.Exec(`
		INSERT INTO projects (
			id, slug, title, description, body, project_type,
			published, updated, downloads, followers,
			categories, loaders, game_versions, icon_url,
			issues_url, source_url, license_id, license_name, license_url,
			client_side, server_side, status, color,
			donation_urls, gallery
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			slug=excluded.slug, title=excluded.title, description=excluded.description,
			body=excluded.body, project_type=excluded.project_type,
			published=excluded.published, updated=excluded.updated,
			downloads=excluded.downloads, followers=excluded.followers,
			categories=excluded.categories, loaders=excluded.loaders,
			game_versions=excluded.game_versions, icon_url=excluded.icon_url,
			issues_url=excluded.issues_url, source_url=excluded.source_url,
			license_id=excluded.license_id, license_name=excluded.license_name,
			license_url=excluded.license_url, client_side=excluded.client_side,
			server_side=excluded.server_side, status=excluded.status,
			color=excluded.color, donation_urls=excluded.donation_urls,
			gallery=excluded.gallery, fetched_at=datetime('now')
	`, project.ID, project.Slug, project.Title, project.Description, project.Body, project.ProjectType,
		project.Published.Format(time.RFC3339), project.Updated.Format(time.RFC3339),
		project.Downloads, project.Followers,
		string(categories), string(loaders), string(gameVersions), project.IconURL,
		project.IssuesURL, project.SourceURL,
		project.License.ID, project.License.Name, project.License.URL,
		project.ClientSide, project.ServerSide, project.Status, project.Color,
		string(donationURLs), string(gallery))

	return err
}

func (c *Cache) GetProject(id string) (*models.Project, error) {
	row := c.db.QueryRow(`
		SELECT id, slug, title, description, body, project_type,
			published, updated, downloads, followers,
			categories, loaders, game_versions, icon_url,
			issues_url, source_url, license_id, license_name, license_url,
			client_side, server_side, status, color,
			donation_urls, gallery
		FROM projects WHERE id = ? OR slug = ?
	`, id, id)

	var p models.Project
	var published, updated string
	var categoriesJSON, loadersJSON, gameVersionsJSON string
	var donationURLsJSON, galleryJSON string
	var issuesURL, sourceURL sql.NullString
	var color sql.NullInt64

	err := row.Scan(
		&p.ID, &p.Slug, &p.Title, &p.Description, &p.Body, &p.ProjectType,
		&published, &updated, &p.Downloads, &p.Followers,
		&categoriesJSON, &loadersJSON, &gameVersionsJSON, &p.IconURL,
		&issuesURL, &sourceURL,
		&p.License.ID, &p.License.Name, &p.License.URL,
		&p.ClientSide, &p.ServerSide, &p.Status, &color,
		&donationURLsJSON, &galleryJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	p.Published, _ = time.Parse(time.RFC3339, published)
	p.Updated, _ = time.Parse(time.RFC3339, updated)

	if issuesURL.Valid {
		p.IssuesURL = &issuesURL.String
	}
	if sourceURL.Valid {
		p.SourceURL = &sourceURL.String
	}
	if color.Valid {
		c := int(color.Int64)
		p.Color = &c
	}

	_ = json.Unmarshal([]byte(categoriesJSON), &p.Categories)
	_ = json.Unmarshal([]byte(loadersJSON), &p.Loaders)
	_ = json.Unmarshal([]byte(gameVersionsJSON), &p.GameVersions)
	_ = json.Unmarshal([]byte(donationURLsJSON), &p.DonationURLs)
	_ = json.Unmarshal([]byte(galleryJSON), &p.Gallery)

	return &p, nil
}

func (c *Cache) CacheVersions(projectID string, versions []models.Version) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, v := range versions {
		gameVersions, _ := json.Marshal(v.GameVersions)
		loaders, _ := json.Marshal(v.Loaders)
		files, _ := json.Marshal(v.Files)
		deps, _ := json.Marshal(v.Dependencies)
		featured := 0
		if v.Featured {
			featured = 1
		}

		_, err := tx.Exec(`
			INSERT INTO versions (
				id, project_id, name, version_number, changelog,
				date_published, downloads, version_type, status, featured,
				game_versions, loaders, files, dependencies
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				project_id=excluded.project_id, name=excluded.name,
				version_number=excluded.version_number, changelog=excluded.changelog,
				date_published=excluded.date_published, downloads=excluded.downloads,
				version_type=excluded.version_type, status=excluded.status,
				featured=excluded.featured, game_versions=excluded.game_versions,
				loaders=excluded.loaders, files=excluded.files,
				dependencies=excluded.dependencies, fetched_at=datetime('now')
		`, v.ID, projectID, v.Name, v.VersionNumber, v.Changelog,
			v.DatePublished.Format(time.RFC3339), v.Downloads,
			string(v.VersionType), v.Status, featured,
			string(gameVersions), string(loaders), string(files), string(deps))
		if err != nil {
			return err
		}

		for _, dep := range v.Dependencies {
			_, err := tx.Exec(`
				INSERT INTO dependencies (version_id, project_id, dep_version_id, file_name, dependency_type)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT DO NOTHING
			`, v.ID, dep.ProjectID, dep.VersionID, dep.FileName, dep.DependencyType)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (c *Cache) GetVersions(projectID string) ([]models.Version, error) {
	rows, err := c.db.Query(`
		SELECT id, project_id, name, version_number, changelog,
			date_published, downloads, version_type, status, featured,
			game_versions, loaders, files, dependencies
		FROM versions WHERE project_id = ?
		ORDER BY date_published DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.Version
	for rows.Next() {
		var v models.Version
		var datePublished string
		var gameVersionsJSON, loadersJSON, filesJSON, depsJSON string
		var featured int

		err := rows.Scan(
			&v.ID, &v.ProjectID, &v.Name, &v.VersionNumber, &v.Changelog,
			&datePublished, &v.Downloads, &v.VersionType, &v.Status, &featured,
			&gameVersionsJSON, &loadersJSON, &filesJSON, &depsJSON,
		)
		if err != nil {
			return nil, err
		}

		v.Featured = featured == 1
		v.DatePublished, _ = time.Parse(time.RFC3339, datePublished)
		_ = json.Unmarshal([]byte(gameVersionsJSON), &v.GameVersions)
		_ = json.Unmarshal([]byte(loadersJSON), &v.Loaders)
		_ = json.Unmarshal([]byte(filesJSON), &v.Files)
		_ = json.Unmarshal([]byte(depsJSON), &v.Dependencies)

		versions = append(versions, v)
	}

	return versions, rows.Err()
}

func (c *Cache) CacheSearch(query string, response *models.SearchResponse, ttl time.Duration) error {
	hash := hashQuery(query)
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	expires := time.Now().Add(ttl).Format(time.RFC3339)

	_, err = c.db.Exec(`
		INSERT INTO search_cache (query_hash, query, response, expires_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(query_hash) DO UPDATE SET
			query=excluded.query, response=excluded.response,
			created_at=datetime('now'), expires_at=excluded.expires_at
	`, hash, query, string(data), expires)

	return err
}

func (c *Cache) GetSearch(query string) (*models.SearchResponse, error) {
	hash := hashQuery(query)

	row := c.db.QueryRow(`
		SELECT response, expires_at FROM search_cache
		WHERE query_hash = ? AND expires_at > datetime('now')
		ORDER BY created_at DESC LIMIT 1
	`, hash)

	var responseJSON, expiresAt string
	if err := row.Scan(&responseJSON, &expiresAt); err != nil {
		return nil, nil
	}

	var result models.SearchResponse
	if err := json.Unmarshal([]byte(responseJSON), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Cache) AddRecentlyViewed(entityType, entityID, title string) {
	_, _ = c.db.Exec(`
		INSERT INTO recently_viewed (entity_type, entity_id, title)
		VALUES (?, ?, ?)
	`, entityType, entityID, title)

	_, _ = c.db.Exec(`DELETE FROM recently_viewed WHERE id NOT IN (
		SELECT id FROM recently_viewed ORDER BY viewed_at DESC LIMIT 50
	)`)
}

func (c *Cache) GetRecentlyViewed(entityType string, limit int) ([]struct {
	EntityType string
	EntityID   string
	Title      string
	ViewedAt   time.Time
}, error) {
	rows, err := c.db.Query(`
		SELECT entity_type, entity_id, title, viewed_at
		FROM recently_viewed
		WHERE entity_type = ?
		ORDER BY viewed_at DESC LIMIT ?
	`, entityType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []struct {
		EntityType string
		EntityID   string
		Title      string
		ViewedAt   time.Time
	}

	for rows.Next() {
		var item struct {
			EntityType string
			EntityID   string
			Title      string
			ViewedAt   time.Time
		}
		var viewedAt string
		if err := rows.Scan(&item.EntityType, &item.EntityID, &item.Title, &viewedAt); err != nil {
			return nil, err
		}
		item.ViewedAt, _ = time.Parse(time.RFC3339, viewedAt)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (c *Cache) GetSetting(key string) (string, error) {
	row := c.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key)
	var value string
	if err := row.Scan(&value); err != nil {
		return "", err
	}
	return value, nil
}

func (c *Cache) SetSetting(key, value string) error {
	_, err := c.db.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}

func (c *Cache) Stats() (projects, versions int, cacheSize int64, err error) {
	_ = c.db.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&projects)
	_ = c.db.QueryRow(`SELECT COUNT(*) FROM versions`).Scan(&versions)

	var size int64
	_ = c.db.QueryRow(`SELECT IFNULL(SUM(LENGTH(response)), 0) FROM search_cache`).Scan(&size)
	cacheSize = size

	return
}

func hashQuery(query string) string {
	h := sha256.Sum256([]byte(query))
	return hex.EncodeToString(h[:])
}

func (c *Cache) SaveDownload(projectName, projectID, versionID, versionNumber, filename, filePath, url string, size int64, status string, isMRPack bool) (int64, error) {
	res, err := c.db.Exec(`
		INSERT INTO downloads (project_name, project_id, version_id, version_number, filename, file_path, url, size, status, is_mrpack, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`, projectName, projectID, versionID, versionNumber, filename, filePath, url, size, status, boolToInt(isMRPack))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (c *Cache) UpdateDownloadStatus(id int64, status string) error {
	_, err := c.db.Exec(`UPDATE downloads SET status = ? WHERE id = ?`, status, id)
	return err
}

func (c *Cache) ListDownloads() ([]downloads.DownloadRecord, error) {
	rows, err := c.db.Query(`
		SELECT id, project_name, project_id, version_id, version_number, filename, file_path, url, size, status, is_mrpack, created_at, installed_at
		FROM downloads ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []downloads.DownloadRecord
	for rows.Next() {
		var r downloads.DownloadRecord
		var isMRPack int
		if err := rows.Scan(&r.ID, &r.ProjectName, &r.ProjectID, &r.VersionID,
			&r.VersionNum, &r.Filename, &r.FilePath, &r.URL, &r.Size,
			&r.Status, &isMRPack, &r.CreatedAt, &r.InstalledAt); err != nil {
			return nil, err
		}
		r.IsMRPack = isMRPack != 0
		records = append(records, r)
	}
	return records, nil
}

func (c *Cache) DeleteAllDownloads() error {
	_, err := c.db.Exec(`DELETE FROM downloads`)
	return err
}

func (c *Cache) DeleteDownload(id int64) error {
	_, err := c.db.Exec(`DELETE FROM downloads WHERE id = ?`, id)
	return err
}

func (c *Cache) MarkInstalled(downloadID int64) error {
	_, err := c.db.Exec(`UPDATE downloads SET installed_at = datetime('now') WHERE id = ?`, downloadID)
	return err
}

func (c *Cache) ClearCache() error {
	_, err := c.db.Exec(`
		DELETE FROM projects;
		DELETE FROM versions;
		DELETE FROM dependencies;
		DELETE FROM search_cache;
		DELETE FROM recently_viewed;
	`)
	return err
}

func (c *Cache) ResetSettings() error {
	_, err := c.db.Exec(`DELETE FROM settings`)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
