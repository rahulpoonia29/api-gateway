package utils

import (
	"log/slog"

	"github.com/armon/go-radix"
)

// To bypass cyclic import issues
type App struct {
	RouteTree *radix.Tree
	Logger    *slog.Logger
}
