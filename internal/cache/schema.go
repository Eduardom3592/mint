package cache

const schema = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    project_type TEXT NOT NULL DEFAULT '',
    published TEXT NOT NULL DEFAULT '',
    updated TEXT NOT NULL DEFAULT '',
    downloads INTEGER NOT NULL DEFAULT 0,
    followers INTEGER NOT NULL DEFAULT 0,
    categories TEXT NOT NULL DEFAULT '[]',
    loaders TEXT NOT NULL DEFAULT '[]',
    game_versions TEXT NOT NULL DEFAULT '[]',
    icon_url TEXT NOT NULL DEFAULT '',
    issues_url TEXT,
    source_url TEXT,
    license_id TEXT NOT NULL DEFAULT '',
    license_name TEXT NOT NULL DEFAULT '',
    license_url TEXT NOT NULL DEFAULT '',
    client_side TEXT NOT NULL DEFAULT '',
    server_side TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT '',
    color INTEGER,
    donation_urls TEXT NOT NULL DEFAULT '[]',
    gallery TEXT NOT NULL DEFAULT '[]',
    fetched_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(slug)
);

CREATE TABLE IF NOT EXISTS versions (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    version_number TEXT NOT NULL,
    changelog TEXT NOT NULL DEFAULT '',
    date_published TEXT NOT NULL DEFAULT '',
    downloads INTEGER NOT NULL DEFAULT 0,
    version_type TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT '',
    featured INTEGER NOT NULL DEFAULT 0,
    game_versions TEXT NOT NULL DEFAULT '[]',
    loaders TEXT NOT NULL DEFAULT '[]',
    files TEXT NOT NULL DEFAULT '[]',
    dependencies TEXT NOT NULL DEFAULT '[]',
    fetched_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id TEXT NOT NULL,
    project_id TEXT,
    dep_version_id TEXT,
    file_name TEXT,
    dependency_type TEXT NOT NULL DEFAULT 'required',
    FOREIGN KEY (version_id) REFERENCES versions(id),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_dependencies_version ON dependencies(version_id);
CREATE INDEX IF NOT EXISTS idx_dependencies_project ON dependencies(project_id);

CREATE TABLE IF NOT EXISTS search_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_hash TEXT NOT NULL,
    query TEXT NOT NULL,
    response TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_search_cache_query ON search_cache(query_hash);

CREATE TABLE IF NOT EXISTS recently_viewed (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    title TEXT NOT NULL,
    viewed_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_recently_viewed_type ON recently_viewed(entity_type);
CREATE INDEX IF NOT EXISTS idx_recently_viewed_at ON recently_viewed(viewed_at);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS downloads (
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
);

CREATE TABLE IF NOT EXISTS installed_modpacks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    version_id TEXT NOT NULL DEFAULT '',
    game_version TEXT NOT NULL DEFAULT '',
    loader TEXT NOT NULL DEFAULT '',
    path TEXT NOT NULL DEFAULT '',
    installed_at TEXT NOT NULL DEFAULT (datetime('now')),
    file_count INTEGER NOT NULL DEFAULT 0,
    verified INTEGER NOT NULL DEFAULT 0
);
`
