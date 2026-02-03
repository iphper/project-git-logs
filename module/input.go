package module

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// 输入组件
type Input struct {
	textInput textinput.Model
	err       error
	// 记录组件值
	value string
}

// NewInput 创建一个输入组件，并运行它，返回输入的值
func NewInput(placeholder string) string {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = placeholder

	input := Input{
		textInput: ti,
		err:       nil,
	}

	p := tea.NewProgram(&input)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(strings.ReplaceAll(input.value, "\n", ""))
}

// 初始化
func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

// 更新事件监听
func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// 记录值
			m.value = m.textInput.Value()
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// 视图渲染
func (m *Input) View() string {
	return m.textInput.View()
}
