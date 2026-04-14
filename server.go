package main

import (
	"net"

	"search_engine/internal/tui"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"github.com/charmbracelet/ssh"
)

func initServer() (*ssh.Server, error) {
	s, err := wish.NewServer(
		ssh.AllocatePty(),
		wish.WithAddress(net.JoinHostPort(HOST, PORT)),
		wish.WithHostKeyPath(KEYHOST_PATH),
		wish.WithMiddleware(
			// activeterm.Middleware(),
			// logging.Middleware(),
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, ok := s.Pty()
				if !ok {
					return nil, []tea.ProgramOption{}
				}
				return tui.CreatePTYModel(
						pty.Window.Width,
						pty.Window.Height,
						pty.Term),
					[]tea.ProgramOption{}
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}
