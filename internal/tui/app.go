package tui

import (
	"fmt"
	"os"
	"strings"

	"ublcli/internal/config"
	"ublcli/internal/orchestrator"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	inputs     []textinput.Model
	focused    int
	err        error
	generating bool
	done       bool
	status     string
}

const (
	excelIdx = iota
	sheetIdx
	configIdx
	templateIdx
	outputIdx
	excelTemplateIdx
	numberIdx
)

func NewModel() *model {
	m := &model{
		inputs: make([]textinput.Model, 7),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()

		switch i {
		case excelIdx:
			t.Placeholder = "Path to Excel file (e.g., input.xlsx)"
			t.Focus()
		case sheetIdx:
			t.Placeholder = "Sheet name (e.g., Data)"
			t.SetValue("Data")
		case configIdx:
			t.Placeholder = "Path to Config file (e.g., billing.yaml)"
			if path, err := config.GetDefaultConfigPath(); err == nil {
				if _, err := os.Stat(path); err == nil {
					t.SetValue(path)
				}
			}
		case templateIdx:
			t.Placeholder = "HTML/Word Template (e.g., invoice, ./tpl.html). Optional if configured."
		case outputIdx:
			t.Placeholder = "Output directory (e.g., ./dist)"
			t.SetValue("./dist")
		case excelTemplateIdx:
			t.Placeholder = "Excel Template (.xlsx, .xlsm). Optional."
		case numberIdx:
			t.Placeholder = "Invoice number (4 digits, ex: 0001)"
			t.CharLimit = 4
		}

		m.inputs[i] = t
	}

	return m
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs)-1 {
				m.generating = true
				m.status = "Generating..."
				return m, m.generateCmd()
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs)-1 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}

	case generateDoneMsg:
		m.generating = false
		if msg.err != nil {
			m.err = msg.err
			m.status = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.done = true
			m.status = "Success! Files generated in " + m.inputs[outputIdx].Value()
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

type generateDoneMsg struct {
	err error
}

func (m *model) generateCmd() tea.Cmd {
	return func() tea.Msg {
		opts := orchestrator.GenerateOptions{
			TimesheetInPath:   m.inputs[excelIdx].Value(),
			SheetName:         m.inputs[sheetIdx].Value(),
			ConfigPath:        m.inputs[configIdx].Value(),
			TemplatePath:      m.inputs[templateIdx].Value(),
			ExcelTemplatePath: m.inputs[excelTemplateIdx].Value(),
			OutputDir:         m.inputs[outputIdx].Value(),
			InvoiceNum:        m.inputs[numberIdx].Value(),
		}

		if len(opts.InvoiceNum) != 4 {
			return generateDoneMsg{err: fmt.Errorf("invoice number must be exactly 4 digits")}
		}
		err := orchestrator.Process(opts)
		return generateDoneMsg{err: err}
	}
}

func (m *model) View() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("UBLCLI - Invoice Generator"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(fmt.Sprintf(
			"%s\n%s\n\n",
			m.inputLabel(i),
			m.inputs[i].View(),
		))
	}

	b.WriteString("\n")
	b.WriteString(m.status)
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("press ENTER to generate • ESC to quit"))

	return b.String()
}

func (m *model) inputLabel(i int) string {
	labels := []string{
		"Timesheet Excel File",
		"Sheet Name",
		"Configuration File (YAML)",
		"HTML Template",
		"Output Directory",
		"Excel Template (Optional)",
		"Invoice Number (4 digits)",
	}
	style := lipgloss.NewStyle()
	if i == m.focused {
		style = style.Foreground(lipgloss.Color("205")).Bold(true)
	}
	return style.Render(labels[i])
}

func Start() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}
