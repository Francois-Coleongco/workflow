package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type choice struct {
	parent_dir string
	child_dir string
}

type model struct {
	choices  []choice
	cursor   int
	selected map[int]struct{}
}

func open_workflow(workflow choice) {

	home_dir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("err occurred getting home_dir", home_dir)
	}	

	cmd := exec.Command("kitty", "--detach", "tmux")

	cmd.Dir = fmt.Sprintf("%s/%s/%s", home_dir, workflow.parent_dir, workflow.child_dir)

	err = cmd.Start()
	if err != nil {
		fmt.Println("err with starting workflow", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("err with workflow", err)
	}

	os.Exit(0)
}


func get_conf_file() (string, string) {

	home_dir, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("homedir could not be found", err)
		home_dir = ""
	}

	
	conf_file_path := fmt.Sprintf("%s/.config/ttyfio/workflows.ttyfio", home_dir)

	file, err := os.Open(conf_file_path)
	if err != nil {
		fmt.Println("couldn't open conf file", err)
	}
	
	file_data := make([]byte, 1024)

	_, err = file.Read(file_data)
	if err != nil {
		return "", home_dir
	}

	return string(file_data), home_dir
}

func get_title() string {

	data, _ := get_conf_file()
	return strings.Split(data, "<title>")[1]
}

func read_parents() ([]string, string) {
	data, home_dir := get_conf_file()
	parsed := strings.Fields(strings.Split(data, "<title>")[2])[0]
	parent_dirs := strings.Split(parsed, ":")
	return parent_dirs, home_dir
}

func read_children(parent_dir string) []string {
	output, err := os.ReadDir(parent_dir)
	var dirs []string
	if err != nil {
		fmt.Println("couldn't read children", err)
	}

	for idx := range output {
		curr_dir := output[idx]

		if curr_dir.Type().IsDir() {
			dirs = append(dirs, curr_dir.Name())
		}
	}

	return dirs
}

func all_dirs_consolidator() []choice {
	var all_dirs []choice

	parent_dirs, home_dir := read_parents()

	fmt.Println("these are parents", parent_dirs)

	counter := 0

	for par_idx := range parent_dirs {
		parent_dir := parent_dirs[par_idx]
		constructed_dir := fmt.Sprintf("%s/%s/", home_dir, parent_dir)

		children_dirs := read_children(constructed_dir)

		for _, child := range children_dirs {
			all_dirs = append(all_dirs, choice{parent_dir, child})
			counter++
		}
	}

	fmt.Println("this is alldirs", all_dirs)
	fmt.Println("end alldirs")
	return all_dirs
}

func initialModel() model {

	return model{
		choices:  all_dirs_consolidator(),
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
