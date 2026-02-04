package module

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

// func main() {
// 	m := Progress{
// 		progress: progress.New(progress.WithDefaultGradient()),
// 	}

// 	if _, err := tea.NewProgram(m).Run(); err != nil {
// 		fmt.Println("Oh no!", err)
// 		os.Exit(1)
// 	}
// }

type tickMsg time.Time

type Progress struct {
	progress progress.Model
	count    float64
	current  float64
}

func NewProgress() *Progress {
	return &Progress{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m *Progress) SetCount(count float64) *Progress {
	m.count = count
	return m
}

func (m *Progress) Add(curr float64) *Progress {
	m.current += curr
	return m
}

func (m *Progress) Init() tea.Cmd {
	return tickCmd()
}

func (m *Progress) Run() error {
	_, e := tea.NewProgram(m).Run()
	return e
}

func (m *Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case tea.KeyMsg: // 不监听按钮事件【等待自动完成】
	// 	return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width
		return m, nil

	case tickMsg:
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.SetPercent(m.current / m.count)
		if m.progress.Percent() == 1 || m.current >= m.count {
			return m, tea.Quit
		}
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m *Progress) View() string {
	return m.progress.View()
}

func tickCmd() tea.Cmd {
	// 10毫秒刷新一次
	return tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
		// 等待调用
		return tickMsg(t)
	})
}
