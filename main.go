package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var addedFiles []string
var index int

type model struct {
	files    []string
	cursor   int
	selected map[int]struct{}
}

func (mod model) addFiles(s string) []string {
	mod.files = append(mod.files, s)
	return mod.files
}

func initialModel() model {
	return model{
		files:    addedFiles,
		selected: make(map[int]struct{}),
	}
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		// if s == "../yahtzee/pdf" {
		// 	cmd := exec.Command("zathura", "../yahtzee/pdf")
		// 	err := cmd.Start()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	err = cmd.Wait()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }
		addedFiles = append(addedFiles, s)
		// println(s)
	}
	return nil
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

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				cmd := exec.Command("zathura", m.files[index])
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				err = cmd.Wait()
				if err != nil {
					log.Fatal(err)
				}
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
	s := "What should we buy at the market?\n\n"

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
			index = i
			delete(m.selected, m.cursor)
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
	filepath.WalkDir("..", walk)
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
