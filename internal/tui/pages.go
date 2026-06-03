package tui

type navPage int

const (
	navHome navPage = iota
	navSearch
	navProject
	navVersions
	navDependencies
	navDownloads
	navCache
	navSettings
	navHelp
)

type searchPage int

const (
	searchList searchPage = iota
	searchProject
	searchVersion
	searchDependency
	searchDownload
)

type projectPage int

const (
	projectContent projectPage = iota
	projectSearch
	projectVersion
	projectDependency
	projectDownload
)

type versionPage int

const (
	versionList versionPage = iota
	versionDetail
	versionProject
	versionDependency
	versionSearch
	versionDownload
)

type dependencyPage int

const (
	dependencyBrowse dependencyPage = iota
)

func (p navPage) supportsDownload() bool {
	return p == navHome || p == navSearch || p == navProject
}
