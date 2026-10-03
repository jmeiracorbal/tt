package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dankkomcg/tilt-tui/internal/api"
	"github.com/dankkomcg/tilt-tui/internal/ui"
)

var version = "dev"

func main() {
	port := flag.Int("port", api.DefaultPort, "Tilt server port")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	addr := fmt.Sprintf("%s:%d", api.DefaultHost, *port)
	p := tea.NewProgram(ui.New(addr), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
