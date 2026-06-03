# mint

![thumbnail](thumbnail.png)

a terminal-native client for modrinth.

browse, explore, download, and analyze content from modrinth - all from your terminal.

## features

- **home** - top 4 most-downloaded mods displayed as cards on startup
- **search** - query modrinth with keyboard-first ux
- **projects** - view project details, descriptions, downloads, followers, categories
- **versions** - browse versions with file info, hashes, and dependencies
- **version picker** - interactive filter (mc version, loader, channel) when downloading from home/search/project pages
- **download manager** - queue, progress bars, speed/eta, hash verification, cancel, retry, open folder, sqlite persistence
- **downloads tab** - dashboard with active/history/installed/failed sections, delete, open folder, clear all
- **mrpack support** - parse, validate, extract, verify, and install modrinth modpacks (.mrpack) with auto-install on download completion
- **cache tab** - browse cached stats and recently viewed projects, re-open with enter, clear cache
- **editable settings** - download directory, max workers, and api key are editable inline; reset to defaults
- **dependency explorer** - expandable/collapsible dependency tree visualization
- **offline cache** - sqlite-backed cache for projects, versions, and search results
- **mouse support** - click navigation, row selection, scrolling, and common actions across tabs

## demos

![social](docs/assets/social.gif)

browse projects, quick download, theme switcher.

![search](docs/assets/search.gif)

search modrinth, view projects, back to search.

![downloads](docs/assets/downloads.gif)

search, quick download, progress bar, history tab.

![themes](docs/assets/themes.gif)

cycle through and apply all built-in themes.

![modpack](docs/assets/modpack.gif)

search modpack, version picker, download, installed tab.

## installation

### from releases

download the latest binary from [github releases](https://github.com/programmersd21/mint/releases).

### from source

```bash
git clone https://github.com/programmersd21/mint.git
cd mint
make build
```

## usage

```bash
mint
```

### keybindings

| key | action |
|-----|--------|
| `j/k` or `up/down` | navigate lists / cycle settings |
| `h/l` or `left/right` | cycle filter / switch tabs |
| `g` / `G` | top / bottom |
| `/` | search |
| `tab` / `shift+tab` | switch pages / done filtering |
| `enter` | select / open / confirm edit |
| `esc` | go back / cancel edit |
| `d` | download (opens version picker) |
| `D` | quick download latest version |
| `i` | inspect mrpack metadata |
| `I` | download & install mrpack |
| `r` | retry failed download |
| `c` | cancel download |
| `o` | open url / reveal file / open downloads folder |
| `delete` | remove download from history |
| `X` | clear all downloaded files and records |
| `C` | clear cache (projects, versions, search, recently viewed) |
| `t` | theme switcher |
| `?` | help |
| `q` / `ctrl+c` | quit |

## development

```bash
make build     # build binary
make test      # run tests
make test-race # run race-enabled tests
make lint      # lint
make cover     # coverage
```

## architecture

```
cmd/mint/           entry point
internal/api/       modrinth api client
internal/cache/     sqlite cache layer + download persister
internal/downloads/ download manager
internal/models/    strongly-typed data models
internal/mrpack/    .mrpack parser, validator, extractor, installer
internal/platform/  cross-platform open commands
internal/tui/       bubble tea tui
```

## license

mit
