package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		if d.IsDir() {
			s = s + "/"
		}
		if !strings.HasPrefix(s, ".") {
			files = append(files, s)
		}
		return nil
	})
	if err != nil {
		log.Println("Error listing files:", err)
	}

	m.files = files
}

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
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "h":
			m.files = []string{}
			splitPath := strings.Split(m.path, "/")
			splitPath = splitPath[:(len(splitPath) - 1)]
			m.path = strings.Join(splitPath, "/")
			print(m.path)
			m.loadFile()
			m.cursor = 0

		case "l":
			m.path = filepath.Join(m.files[m.cursor])
			m.loadFile()
			m.cursor = 0

		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

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

	return m, nil
}

func (m model) View() string {
	s := "Select file:\n\n"

	for i, choice := range m.files {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	model := initialModel()
	path, err := os.Getwd()
	if err != nil {
		println(err)
	}
	model.path = path
	model.loadFile()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
