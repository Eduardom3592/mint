Output docs/assets/modpack.gif
Require go
Set Width 1600
Set Height 900
Set FontSize 18
Set TypingSpeed 40ms

Hide
Type "go run ./..."
Sleep 3s
Enter
Sleep 3s
Show

# home page
Sleep 1.5s

# tab to search page (search input already focused)
Tab
Sleep 1.5s

# search for modpack
Type "fabulously optimized"
Sleep 1.5s

# results load
Sleep 2.5s

# open project
Enter
Sleep 2s

# scroll project info
Down
Sleep 500ms
Down
Sleep 500ms

# open version picker for download
Type "d"
Sleep 1.5s

# version picker shows versions list
Sleep 2s

# select first version
Enter
Sleep 600ms

# download primary file
Enter
Sleep 1.5s

# auto-navigated to downloads, progress
Sleep 3s

# download completing
Sleep 1.5s

# switch to installed modpacks tab (two rights)
Type "l"
Sleep 500ms
Type "l"
Sleep 800ms

# installed modpacks section visible
Sleep 1.5s

# quit
Type "q"
