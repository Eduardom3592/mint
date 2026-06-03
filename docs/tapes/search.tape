Output docs/assets/search.gif
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

# tab to search page
Tab
Sleep 1.5s

# type first query (search input already focused)
Type "iris"
Sleep 1s

# results load
Sleep 2.5s

# browse results
Down
Sleep 500ms
Down
Sleep 500ms

# open first result
Up
Sleep 300ms
Up
Sleep 300ms
Enter
Sleep 2s

# project page visible
Sleep 1s

# go back to search
Escape
Sleep 600ms

# focus search input
Type "/"
Sleep 300ms

# clear query
Ctrl+u
Sleep 300ms

# type new search
Type "fabric api"
Sleep 1.5s

# new results appear
Sleep 2.5s

# open result
Enter
Sleep 2s

# quit
Type "q"
