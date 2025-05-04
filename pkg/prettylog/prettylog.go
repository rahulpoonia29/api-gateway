package prettylog

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	writer io.Writer
	level  slog.Level
	logger *log.Logger

	attrs []slog.Attr
	group string
}

func NewPrettyHandler(w io.Writer, level slog.Level) *PrettyHandler {
	return &PrettyHandler{
		writer: w,
		level:  level,
		logger: log.New(w, "", 0),
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *PrettyHandler) Handle(_ context.Context, record slog.Record) error {
	// Merge static and record attributes
	attrs := make([]slog.Attr, 0, len(h.attrs)+int(record.NumAttrs()))
	attrs = append(attrs, h.attrs...)
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	// Format timestamp
	timestamp := color.New(color.FgHiBlack).Sprintf(record.Time.Format("[15:04:05.000]"))

	// Colorized log level
	levelStr := strings.ToUpper(record.Level.String()) + ":"
	switch record.Level {
	case slog.LevelDebug:
		levelStr = color.MagentaString(levelStr)
	case slog.LevelInfo:
		levelStr = color.BlueString(levelStr)
	case slog.LevelWarn:
		levelStr = color.YellowString(levelStr)
	case slog.LevelError:
		levelStr = color.RedString(levelStr)
	default:
		levelStr = color.WhiteString(levelStr)
	}

	// Colorized message
	message := color.CyanString(record.Message)

	// Print log header
	h.logger.Println(timestamp, levelStr, message)

	// Format and print attributes
	if len(attrs) > 0 {
		attrLine := formatAttrs(attrs, h.group)
		h.logger.Println("   ", attrLine)
	}

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append([]slog.Attr{}, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &PrettyHandler{
		writer: h.writer,
		level:  h.level,
		logger: h.logger,
		attrs:  newAttrs,
		group:  h.group,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newGroup := h.group
	if newGroup != "" {
		newGroup += "." + name
	} else {
		newGroup = name
	}
	return &PrettyHandler{
		writer: h.writer,
		level:  h.level,
		logger: h.logger,
		attrs:  h.attrs,
		group:  newGroup,
	}
}

// formatAttrs returns a single line of key=value pairs, colorized
func formatAttrs(attrs []slog.Attr, group string) string {
	var parts []string
	prefix := ""
	if group != "" {
		prefix = group + "."
	}

	for _, attr := range attrs {
		if attr.Value.Kind() == slog.KindGroup {
			for _, sub := range attr.Value.Group() {
				key := color.GreenString("%s%s.%s", prefix, attr.Key, sub.Key)
				val := color.WhiteString("%v", sub.Value.Any())
				parts = append(parts, fmt.Sprintf("%s=%s", key, val))
			}
		} else {
			key := color.GreenString("%s%s", prefix, attr.Key)
			val := color.WhiteString("%v", attr.Value.Any())
			parts = append(parts, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return strings.Join(parts, " ")
}
