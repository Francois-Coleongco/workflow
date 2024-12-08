package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func open_workflow(workflow string) error {
	cmd := exec.Command("kitty", "--detach", "tmux")

	cmd.Dir = workflow

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}

	return nil
}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func read_from_workflows_file() []string {

	file, err := os.Open("/home/hitori/.config/ttyfio/workflows.ttyfio")

	file_data := make([]byte, 1024)

	file.Read(file_data)

	workflows := strings.Split(string(file_data), "\n")

	if err != nil {
		return []string{}
	}

	workflows_count := len(workflows) - 1 // minus one becuase it counts the \n for the new line at the end of the file
	return workflows[:workflows_count]
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
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}
	}
	return m, nil
}

func (m model) View() string {
	s := `                                                         ╋╋╋┏┓╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋┏┓╋╋╋┏┓
                                                        ┏┳┳┫┗┳━┓┏━┓┏┳┳━┓┏┳┳━┳┳┓┃┗┳━┳┛┣━┓┏┳┓
                                                        ┃┃┃┃┃┃╋┃┃╋┗┫┏┫┻┫┃┃┃╋┃┃┃┃┏┫╋┃╋┃╋┗┫┃┃
                                                        ┗━━┻┻┻━┛┗━━┻┛┗━┛┣┓┣━┻━┛┗━┻━┻━┻━━╋┓┃
                                                        ╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋╋┗━┛╋╋╋╋╋╋╋╋╋╋╋╋╋┗━┛`
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
