package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/carlfung1003/ssh-portfolio/internal/app"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
)

func main() {
	data := content.Load()

	// SSH server mode
	if len(os.Args) > 1 && os.Args[1] == "--serve" {
		runSSHServer(data)
		return
	}

	// Local TUI mode
	m := app.NewModel(data, 0, 0)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runSSHServer(data content.Data) {
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2222"
	}

	host := os.Getenv("SSH_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	keyPath := os.Getenv("SSH_HOST_KEY")
	if keyPath == "" {
		keyPath = ".ssh/id_ed25519"
	}

	handler := func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, _ := sess.Pty()
		m := app.NewModel(data, pty.Window.Width, pty.Window.Height)
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bm.Middleware(handler),
			activeterm.Middleware(),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error("Server error", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error("Shutdown error", "error", err)
	}
}
