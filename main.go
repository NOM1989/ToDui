package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// to_do_filename = "/home/NOM/todo.txt"
	to_do_filename = "/home/NOM/Code/Go/ToDui/list.txt"
)

type model struct {
	cursor   int
	mode     string
	items    []string
	selected map[int]struct{}
}

func initialModel() model {
	data, err := os.ReadFile(to_do_filename)
	if err != nil {
		panic(err)
	}

	file_items := []string{}
	for _, todo_item := range strings.Split(string(data), "\n") {
		if todo_item != "" {
			file_items = append(file_items, todo_item)
		}
	}

	return model{
		items: file_items,

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if len(m.selected) != 0 {
				items_string := ""
				for index, todo_item := range m.items {
					_, ok := m.selected[index]
					if !ok {
						items_string += todo_item + "\n"
					}
				}

				err := os.WriteFile(to_do_filename, []byte(items_string), 0644)
				if err != nil {
					panic(err)
				}
			}
			return m, tea.Quit
		case "e":
			if m.mode == "" {
				m.mode = "edit"
			} else {
				m.mode = ""
			}
		case "up", "k":
			if m.cursor > 0 {
				if m.mode == "edit" {
					tmp_item := m.items[m.cursor]
					m.items[m.cursor] = m.items[m.cursor-1]
					m.items[m.cursor-1] = tmp_item
				}
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				if m.mode == "edit" {
					tmp_item := m.items[m.cursor]
					m.items[m.cursor] = m.items[m.cursor+1]
					m.items[m.cursor+1] = tmp_item
				}
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := m.mode + fmt.Sprint(m.cursor) + "Your to-do list:\n\n"

	for i, choice := range m.items {
		cursor := " "
		if m.cursor == i {
			if m.mode == "edit" {
				cursor = " >"
			} else {
				cursor = ">"
			}
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to exit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
