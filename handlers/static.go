package handlers

import (
	"io"
	"path"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/labstack/echo"
)

var contentTypes = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "text/javascript",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"txt":  "text/plain",
}

const defaultContentType = "application/octet-stream"

func Static(root string, fs http.FileSystem) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			p := c.Request().URL.Path
			p, err = echo.PathUnescape(p)
			if err != nil {
				err = echo.NewHTTPError(http.StatusBadRequest, "unable to unescape")
				return
			}
			name := filepath.Join(root, path.Clean("/"+p))
			f, err := fs.Open(name)
			if err != nil {
				err = echo.NewHTTPError(http.StatusNotFound, "not found")
				return
			}
			defer f.Close()
			return c.Stream(200, detectContentType(name, f), f)
		}
	}
}

func detectContentType(name string, f http.File) string {
	// detect from the file extension
	ext := filepath.Ext(name)
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	contentType := contentTypes[ext]
	if len(contentType) > 0 {
		return contentType
	}

	// detect from the content
	buf := make([]byte, 0x100)
	n, err := f.Read(buf)
	defer func() {
		f.Seek(0, 0)
		f.Close()
	}()
	if err == nil || err == io.EOF {
		return http.DetectContentType(buf[:n])
	}

	return defaultContentType
}
