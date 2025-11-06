package logging

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// Logger captures structured, persona-driven logging.
type Logger interface {
	Info(cmd *cobra.Command, format string, args ...interface{})
}

// SheldonLogger renders progress updates with Sheldon's unique flair.
type SheldonLogger struct {
	rng    *rand.Rand
	prefix []string
}

// NewSheldonLogger constructs a logger that channels Sheldon Cooper.
func NewSheldonLogger() *SheldonLogger {
	return &SheldonLogger{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		prefix: []string{
			"Observation:",
			"Newsflash:",
			"Hierarchy update:",
			"Cerebral bulletin:",
			"Routine verification:",
			"Minor inconvenience:",
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
	fmt.Fprintf(writer, "Sheldon Cooper %s %s Bazinga.\n", prefix, msg)
}
