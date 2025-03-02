# 🥭 Mango Bubbletea

![Alt Text](https://github.com/executionreverted/mango-bubbletea/blob/master/views.gif)

Hey there! 👋 This is my little terminal UI app that I built while learning [Bubble Tea](https://github.com/charmbracelet/bubbletea). It's still a work in progress as I figure stuff out, but I wanted to share it for example purposes.

## What's This?

I wanted to learn Golang and how to make cool terminal UIs and stumbled upon Bubble Tea. This project is basically me experimenting with it and trying to build something useful along the way. I'm still learning Go and TUI development, so expect some rough edges!

## Cool Stuff It (kinda) Can Do

- Multiple pages (Home, Settings, About) that you can switch between
- Fancy modal dialogs that pop up over your content
- Keyboard shortcuts you can customize
- Looks good even when you resize your terminal

## Project Structure

```
bubbletea-app/
├── app/
│   ├── actions/
│   ├── components/
│   ├── config/
│   ├── global/
│   ├── pages/
│   └── styles/
└── main.go
```

(There's more files in there but you get the idea)

## Custom Keybindings

You can set up your own keyboard shortcuts by creating a file at `~/.config/sleek/keymap.json`.

### Keys You Can Change (WIP, idk)

```
Up, Down, Left, Right - Navigation
Help - Show help screen
Quit - Exit the app
Home, Settings, About - Jump to pages
Enter - Confirm stuff
Esc - Get out of things
Back - Go back
J, K - Vi-style Navigation

check @config/Keybindings.go
```

### Example Config -> ~/.config/sleek/keymap.json

```json
{
  "Enter": "-",
}
```

## Modal Windows

I finally got modals working! They overlay on top of the current content instead of replacing it. Had a bit of trouble centering them properly.

Use left/right arrows to pick buttons and Enter to confirm. Esc to bail out.
Modals can be called from anywhere, with any callback

## Running It

```bash
# Clone this thing
git clone https://github.com/yourusername/bubbletea-app.git

# Go to the folder
cd bubbletea-app

# Install packages.. 

go mod tidy

# Run
go run . 
```

## Dependencies - install with  "go mod tidy"

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The main framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - For making things pretty
- [Bubbles](https://github.com/charmbracelet/bubbles) - For UI components

## What's Next?

I'm still learning and adding features as I go! Some things I want to add:

- Better theme support
- More interactive components
- Proper error handling (oops)
- Documentation (maybe someday)

Let me know if you have any cool ideas or find bugs! This is a fun project for me to learn Go and terminal UIs.

## License

Do whatever you want with it! (MIT)
