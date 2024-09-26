package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	files    []string
	cursor   int
	selected map[int]struct{}
	path     string
}

func initialModel() model {
	return model{
		selected: make(map[int]struct{}),
	}
}

func (m *model) loadFile() {
	var files []string
	err := filepath.WalkDir(m.path, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, s)
		}
		return nil
	})
	if err != nil {
		log.Println("Error listing files:", err)
	}

	m.files = files
}

// makes openning zathura async
func openFileAsync(file string) {
	go func() {
		cmd := exec.Command("zathura", file)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "l":
			m.path = ".."
			m.loadFile()

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

			//o to match the same as yazi
		case "o":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				openFileAsync(m.files[m.cursor])
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select file:\n\n"

	// Iterate over our choices
	for i, choice := range m.files {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	model := initialModel()
	model.path = "."
	model.loadFile()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
