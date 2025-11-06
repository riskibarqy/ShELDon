package commands

import (
	"github.com/riskiramdan/ShELDon/internal/config"
	"github.com/riskiramdan/ShELDon/internal/git"
	"github.com/riskiramdan/ShELDon/internal/llm"
	"github.com/riskiramdan/ShELDon/internal/logging"
	"github.com/riskiramdan/ShELDon/internal/system"
)

// Dependencies lists all cross-cutting services required by CLI commands.
type Dependencies struct {
	Config *config.Config
	LLM    llm.Client
	Files  system.FileManager
	Git    git.Client
	Shell  system.Shell
	Logger logging.Logger
}
