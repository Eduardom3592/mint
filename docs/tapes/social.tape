Output docs/assets/social.gif
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

# home page with top mod cards visible
Sleep 1.5s

# tab to search page (search input already focused)
Tab
Sleep 1.5s

# search for sodium
Type "sodium"
Sleep 1.2s

# results populate
Sleep 2.5s

# open sodium project
Enter
Sleep 2s

# scroll to see project details
Down
Sleep 500ms
Down
Sleep 500ms

# quick download via version picker
Type "D"
Sleep 2s

# auto-navigated to downloads with progress
Sleep 3s

# download completes
Sleep 1.5s

# open theme switcher
Type "t"
Sleep 600ms

# cycle through themes with live preview
Down
Sleep 800ms
Down
Sleep 800ms
Down
Sleep 800ms
Down
Sleep 800ms

# dismiss theme switcher
Escape
Sleep 600ms

# quit
Type "q"
