package main

import (
	"fmt"
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

	f, err := os.Open(m.path)

	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := f.ReadDir(-1)

	f.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
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
			splitPath := strings.Split(m.path, "/")
			splitPath = splitPath[:(len(splitPath) - 1)]
			m.path = strings.Join(splitPath, "/")
			m.loadFile()
			m.cursor = 0

		case "l":
			m.cursor = 0
			m.path = filepath.Join(m.files[m.cursor])
			m.loadFile()

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
