package logging

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Logger captures structured, persona-driven logging.
type Logger interface {
	Info(cmd *cobra.Command, format string, args ...interface{})
}

// SheldonLogger renders progress updates with Sheldon's unique flair.
type SheldonLogger struct {
	rng      *rand.Rand
	identity []string
	prefix   []string
	spinner  []string
	colors   []string
	captions []string
	signoff  []string
	verbs    []string
}

// NewSheldonLogger constructs a logger that channels Sheldon Cooper.
func NewSheldonLogger() *SheldonLogger {
	return &SheldonLogger{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		identity: []string{
			"Sheldon Cooper",
		},
		prefix: []string{
			"Observation:",
			"Newsflash:",
			"Hierarchy update:",
			"Cerebral bulletin:",
			"Routine verification:",
			"Minor inconvenience:",
		},
		spinner: []string{"⠋", "⠙", "⠸", "⠴", "⠦", "⠇"},
		colors: []string{
			"\033[38;5;39m",
			"\033[38;5;45m",
			"\033[38;5;81m",
			"\033[38;5;118m",
			"\033[38;5;214m",
		},
		captions: []string{
			"Aligning sarcasm matrix",
			"Buffering unsolicited advice",
			"Optimizing ego subroutines",
			"Recalibrating bazinga drive",
			"Projecting smug certainty",
		},
		signoff: []string{
			"Bazinga.",
			"Succinct perfection achieved.",
			"Puzzle solved, you're welcome.",
			"Ego stabilized, for now.",
			"Try to keep up.",
			"Controlled brilliance complete.",
		},
		verbs: []string{
			"reports",
			"annotates",
			"notes",
			"declares",
			"broadcasts",
			"observes",
			"documents",
		},
	}
}

// Info emits a formatted line to the command's stderr with a Sheldon-esque prefix.
func (l *SheldonLogger) Info(cmd *cobra.Command, format string, args ...interface{}) {
	if cmd == nil {
		return
	}
	writer := cmd.ErrOrStderr()
	msg := fmt.Sprintf(format, args...)
	prefix := l.prefix[l.rng.Intn(len(l.prefix))]
	l.animate(writer)
	lead := "Sheldon Cooper"
	if len(l.identity) > 0 {
		lead = l.identity[l.rng.Intn(len(l.identity))]
	}
	if l.supportsColor(writer) {
		color := l.colors[l.rng.Intn(len(l.colors))]
		lead = fmt.Sprintf("%s%s\033[0m", color, lead)
	}
	verb := "reports"
	if len(l.verbs) > 0 {
		verb = l.verbs[l.rng.Intn(len(l.verbs))]
	}
	signoff := "Bazinga."
	if len(l.signoff) > 0 {
		signoff = l.signoff[l.rng.Intn(len(l.signoff))]
	}
	fmt.Fprintf(writer, "%s %s %s %s %s\n", lead, verb, prefix, msg, signoff)
}

func (l *SheldonLogger) animate(w io.Writer) {
	if !l.supportsColor(w) || len(l.spinner) == 0 {
		return
	}
	caption := l.captions[l.rng.Intn(len(l.captions))]
	frames := len(l.spinner)
	cycles := l.rng.Intn(2) + 1
	delay := time.Duration(l.rng.Intn(50)+50) * time.Millisecond
	for i := 0; i < cycles*frames; i++ {
		frame := l.spinner[i%frames]
		fmt.Fprintf(w, "\r%s %s", frame, caption)
		time.Sleep(delay)
	}
	fmt.Fprint(w, "\r\033[K")
}

func (l *SheldonLogger) supportsColor(w io.Writer) bool {
	type fdWriter interface {
		Fd() uintptr
	}
	if f, ok := w.(fdWriter); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}
