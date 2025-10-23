package git

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Writer writes git events to a JSONL file.
type Writer struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filepath string) (*Writer, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return &Writer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// Write writes a single event to the JSONL file.
func (w *Writer) Write(event *Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if _, err := w.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	if err := w.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// WriteAll writes multiple events to the JSONL file.
func (w *Writer) WriteAll(events []*Event) error {
	for _, event := range events {
		if err := w.Write(event); err != nil {
			return err
		}
	}
	return nil
}

// Flush flushes any buffered data to the file.
func (w *Writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.writer.Flush()
}

// Close flushes and closes the writer.
func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.file.Close()
}
