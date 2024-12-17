package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
)

func open_workflow(workflow string) {
	cmd := exec.Command("kitty", "--detach", "tmux")

	cmd.Dir = workflow

	err := cmd.Start()
	if err != nil {
		fmt.Println("err with starting workflow", err)
	}
	err = cmd.Wait()

	if err != nil {
		fmt.Println("err with workflow", err)
	}

	os.Exit(0)
}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func get_conf_file() string {
	file, err := os.Open("/home/hitori/.config/ttyfio/workflows.ttyfio")

	file_data := make([]byte, 1024)
	file.Read(file_data)
	if err != nil {
		return ""
	}
	return string(file_data)
}
func get_title() string {
	return strings.Split(get_conf_file(), "<title>")[1]
}

func read_from_workflows_file() []string {
	data := get_conf_file()
	workflows := strings.Fields(strings.Split(data, "<title>")[2])
	return workflows[:len(workflows)-1]

}

func initialModel() model {
	return model{
		choices:  read_from_workflows_file(),
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
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			open_workflow(m.choices[m.cursor])
		}
	}
	return m, nil
}

func (m model) View() string {
	s := get_title()

	s += "\n\n\n"
	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
		}

		s += fmt.Sprintf("\n                                                 %s [%s] %s\n\n", cursor, checked, choice)

	}

	s += "\n(q) There is peace in death.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
