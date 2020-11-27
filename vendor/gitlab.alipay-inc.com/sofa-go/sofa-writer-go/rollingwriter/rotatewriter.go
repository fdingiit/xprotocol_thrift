package rollingwriter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	DefaultTimeRollingPerDayFormat    = "2006-01-02"
	DefaultTimeRollingPerHourFormat   = "2006-01-02-15"
	DefaultTimeRollingPerMinuteFormat = "2006-01-02-15_04"
	DefaultTimeRollingPerSecondFormat = "2006-01-02-15_04_05"
)

type RotateWriter interface {
	io.WriteCloser
	Rotate(name, rotatename string) error
}

type FileRotateWriter struct {
	file *os.File
}

func NewFileRotateWriter(filename string) (*FileRotateWriter, error) {
	err := os.MkdirAll(filepath.Dir(filename), 0744)
	if err != nil {
		return nil, err
	}

	// TODO: add the truncate flag?
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return &FileRotateWriter{
		file: f,
	}, nil
}

func (w *FileRotateWriter) Close() error {
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *FileRotateWriter) Write(p []byte) (int, error) {
	if w.file == nil {
		return 0, fmt.Errorf("no available file")
	}

	return w.file.Write(p)
}

func (w *FileRotateWriter) Rotate(filename, rotatename string) error {
	if err := w.Close(); err != nil {
		return err
	}

	if err := w.openNew(filename, rotatename); err != nil {
		return err
	}

	return nil
}

func (w *FileRotateWriter) openNew(filename, rotatename string) error {
	// 0755 is safe to logging
	// nolint
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	info, err := os.Stat(filename)
	mode := os.FileMode(0644)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()
		// move the existing file
		if err = os.Rename(filename, rotatename); err != nil {
			return err
		}

		// this is a no-op anywhere but linux
		if err = chown(filename, info); err != nil {
			return err
		}
	}

	// TODO: add the truncate flag?
	f, ferr := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, mode)
	if ferr != nil {
		return ferr
	}
	w.file = f

	return nil
}
