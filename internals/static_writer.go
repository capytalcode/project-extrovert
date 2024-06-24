package internals

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
)

const PERMISSIONS = 0755

type Page struct {
	Path      string
	Component templ.Component
}

type StaticWriter struct {
	DistDir   *string
	StaticDir *string
	Pages     []Page
	Context   context.Context
	Logger    log.Logger
}

func (w *StaticWriter) WritePage(path string, writer func(ctx context.Context, w io.Writer) error) error {
	directory := filepath.Dir(path)
	err := os.MkdirAll(directory, PERMISSIONS)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = writer(w.Context, f)
	return err
}

func (w *StaticWriter) WriteAll() error {
	for _, page := range w.Pages {
		p := filepath.Join(*w.DistDir, page.Path)
		w.Logger.Printf("Writing page %s", p)
		err := w.WritePage(p, page.Component.Render)
		if err != nil {
			return err
		}
	}

	err := filepath.WalkDir(*w.StaticDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if d.IsDir() || path == *w.StaticDir {
			return nil
		}

		f, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		s, err := filepath.Abs(*w.StaticDir)
		if err != nil {
			return err
		}

		err = w.CopyStatic(strings.TrimPrefix(f, s))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (w *StaticWriter) CopyStatic(path string) error {
	c, err := os.ReadFile(filepath.Join(*w.StaticDir, path))
	if err != nil {
		return err
	}

	p := filepath.Join(*w.DistDir, path)
	err = os.MkdirAll(filepath.Dir(p), PERMISSIONS)
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := f.Write(c)
	if err != nil {
		return err
	}
	w.Logger.Printf("Wrote %v bytes in %s", b, p)

	return nil
}
