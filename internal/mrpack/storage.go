package mrpack

import (
	"database/sql"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveInstall(record *InstallRecord) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`
		INSERT INTO installed_modpacks (name, version_id, game_version, loader, path, installed_at, file_count, verified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, record.Name, record.VersionID, record.GameVersion, record.Loader,
		record.Path, record.InstalledAt, record.FileCount, record.Verified)
	return err
}

func (s *Storage) ListInstalled() ([]InstallRecord, error) {
	if s.db == nil {
		return nil, nil
	}
	rows, err := s.db.Query(`
		SELECT id, name, version_id, game_version, loader, path, installed_at, file_count, verified
		FROM installed_modpacks ORDER BY installed_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []InstallRecord
	for rows.Next() {
		var r InstallRecord
		if err := rows.Scan(&r.ID, &r.Name, &r.VersionID, &r.GameVersion, &r.Loader,
			&r.Path, &r.InstalledAt, &r.FileCount, &r.Verified); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

func (s *Storage) GetInstalled(id int64) (*InstallRecord, error) {
	if s.db == nil {
		return nil, nil
	}
	row := s.db.QueryRow(`
		SELECT id, name, version_id, game_version, loader, path, installed_at, file_count, verified
		FROM installed_modpacks WHERE id = ?
	`, id)
	var r InstallRecord
	err := row.Scan(&r.ID, &r.Name, &r.VersionID, &r.GameVersion, &r.Loader,
		&r.Path, &r.InstalledAt, &r.FileCount, &r.Verified)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Storage) DeleteInstalled(id int64) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`DELETE FROM installed_modpacks WHERE id = ?`, id)
	return err
}

type DownloadRecord struct {
	ID          int64
	ProjectName string
	ProjectID   string
	VersionID   string
	VersionNum  string
	Filename    string
	FilePath    string
	URL         string
	Size        int64
	Status      string
	IsMRPack    bool
	CreatedAt   string
	InstalledAt *string
}

func (s *Storage) SaveDownload(d *DownloadRecord) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`
		INSERT INTO downloads (project_name, project_id, version_id, version_number, filename, file_path, url, size, status, is_mrpack, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, d.ProjectName, d.ProjectID, d.VersionID, d.VersionNum, d.Filename,
		d.FilePath, d.URL, d.Size, d.Status, d.IsMRPack, d.CreatedAt)
	return err
}

func (s *Storage) UpdateDownloadStatus(id int64, status string) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`UPDATE downloads SET status = ? WHERE id = ?`, status, id)
	return err
}

func (s *Storage) ListDownloads() ([]DownloadRecord, error) {
	if s.db == nil {
		return nil, nil
	}
	rows, err := s.db.Query(`
		SELECT id, project_name, project_id, version_id, version_number, filename, file_path, url, size, status, is_mrpack, created_at, installed_at
		FROM downloads ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []DownloadRecord
	for rows.Next() {
		var r DownloadRecord
		var installedAt *string
		if err := rows.Scan(&r.ID, &r.ProjectName, &r.ProjectID, &r.VersionID,
			&r.VersionNum, &r.Filename, &r.FilePath, &r.URL, &r.Size,
			&r.Status, &r.IsMRPack, &r.CreatedAt, &installedAt); err != nil {
			return nil, err
		}
		r.InstalledAt = installedAt
		records = append(records, r)
	}
	return records, nil
}

func (s *Storage) DeleteDownload(id int64) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.Exec(`DELETE FROM downloads WHERE id = ?`, id)
	return err
}

func (s *Storage) MarkInstalled(downloadID int64) error {
	if s.db == nil {
		return nil
	}
	now := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec(`UPDATE downloads SET installed_at = ? WHERE id = ?`, now, downloadID)
	return err
}
