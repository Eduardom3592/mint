Output docs/assets/downloads.gif
Require go

Set Width 1600
Set Height 900
Set FontSize 18
Set TypingSpeed 55ms

Hide
Type "go run ./..."
Sleep 3s
Enter
Sleep 3s

Show

# tab to search page (search input already focused)
Tab
Sleep 2s

# search for lithium
Type "lithium"
Sleep 1.2s

# results load (simulate real lag)
Sleep 2.5s

# open project
Enter
Sleep 1.5s

# small UI transition delay (important on low-end rendering)
Sleep 1s

# quick download auto-enqueues and navigates to downloads
Type "D"
Sleep 1.5s

# download in progress with progress bar (simulate real network + UI updates)
Sleep 3.5s

# completion render delay
Sleep 1.5s

# switch to history tab
Type "l"
Sleep 800ms

# history render (slower DOM update simulation)
Sleep 1.5s

# quit
Type "q"
