package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func Connect(cmd *cobra.Command, args []string) {
	user, err := cmd.Flags().GetString("user")
	if err != nil {
		fmt.Println("Error get user flag:", err)
	}

	url := fmt.Sprintf("ws://localhost:1804/ws?user=%s", user)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error connect to server: %v", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error read message:", err)
				conn.Close()
				break
			}
			fmt.Println(string(message))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("send message: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading scanner:", err)
			}
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Message cannot be empty!")
			continue
		}

		if input == "exit" || input == "quit" {
			conn.Close()
			fmt.Println("Exiting application. Goodbye!")
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(input)); err != nil {
			fmt.Println("Error write message:", err)
			continue
		}
	}
}

type webSocketMsg struct {
	msg string
}

type model struct {
	conn        *websocket.Conn
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func ConnectWithBubbleTea(cmd *cobra.Command, args []string) {
	user, err := cmd.Flags().GetString("user")
	if err != nil {
		fmt.Println("Error get user flag:", err)
	}

	url := fmt.Sprintf("ws://localhost:1804/ws?user=%s", user)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error connect to server: %v", err)
	}
	defer conn.Close()

	p := tea.NewProgram(initialModel(conn))

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error read message:", err)
				conn.Close()
				break
			}
			p.Send(webSocketMsg{msg: string(message)})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

func initialModel(conn *websocket.Conn) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "|"
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome!
Type a message and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		conn:        conn,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

type WebSocketMsg string

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case webSocketMsg:
		m.messages = append(m.messages, msg.msg)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			m.conn.Close()
			return m, tea.Quit
		case "enter":
			if err := m.conn.WriteMessage(websocket.TextMessage, []byte(m.textarea.Value())); err != nil {
				fmt.Println("Error write message:", err)
			}
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.textarea.View())
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
