package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type ConsoleHandler struct {
	opts  slog.HandlerOptions
	out   io.Writer
	attrs []slog.Attr
}

func NewConsoleHandler(out io.Writer) *ConsoleHandler {
	return &ConsoleHandler{
		opts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
		out: out,
	}
}

func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ConsoleHandler) Handle(ctx context.Context, record slog.Record) error {
	var b strings.Builder

	
	b.WriteString(record.Time.Format("15:04:05"))

	b.WriteString(" ")
	b.WriteString(fmt.Sprintf("%-5s", record.Level.String()))

	// Buscar el módulo entre los atributos heredados
	var module string

	for _, attr := range h.attrs {
		if attr.Key == "module" {
			module = attr.Value.String()
			break
		}
	}

	
	if module != "" {
		b.WriteString(" [")
		b.WriteString(module)
		b.WriteString("]")
	}

	
	b.WriteString(" ")
	b.WriteString(record.Message)
  
	// aqui se recorre la lista de los atributos heredados (excepto el modulo)
	for _, attr := range h.attrs {
		if attr.Key == "module" {
			continue
		}

		b.WriteString(" ")
		b.WriteString(attr.Key)
		b.WriteString("=")
		b.WriteString(attr.Value.String())
	}

	// Atributos propios del log
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "module" {
			return true
		}

		b.WriteString(" ")
		b.WriteString(attr.Key)
		b.WriteString("=")
		b.WriteString(attr.Value.String())

		return true
	})

	b.WriteString("\n")

	_, err := h.out.Write([]byte(b.String()))
	return err
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))

	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	return &ConsoleHandler{
		opts:  h.opts,
		out:   h.out,
		attrs: newAttrs,
	}
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}