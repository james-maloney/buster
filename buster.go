package buster

import (
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	noCacheHeader = "no-store"
	cacheHeader   = "max-age=31536000"
)

func NewFileServer(rootDir, prefix string) *FileServer {
	return &FileServer{
		RootDir:         rootDir,
		Prefix:          prefix,
		Buster:          Buster,
		BuildBusterFunc: DefaultBuildBuster,
		StripBusterFunc: DefaultStripBuster,
		BuildURLFunc:    DefaultBuildURL,
	}
}

type FileServer struct {
	RootDir      string // File system directory that will be served
	Prefix       string // URL prefix
	Buster       string // Cache buster that will be added to the Prefix
	Host         string // include for absolute urls (example http://example.com)
	DisableCache bool   // don't cache

	prefix          string
	BuildBusterFunc func(s *FileServer) string                      // func used to build the path prefix
	StripBusterFunc func(path string, s *FileServer) (string, bool) // func used to strip the prefix
	BuildURLFunc    func(path string, s *FileServer) string         // func use to create URLs with the cache buster
}

// BuildURL
func (fs *FileServer) BuildURL(url string) string {
	return fs.BuildURLFunc(url, fs)
}

// Buster is used as the cache buster slug
var Buster = random(10)

func DefaultBuildBuster(s *FileServer) string {
	return s.Prefix + s.Buster + "/"
}

func DefaultStripBuster(p string, s *FileServer) (string, bool) {
	l := len(s.prefix)
	if len(p) < l {
		return "", false
	}
	return path.Join(s.RootDir, p[l:]), true
}

func DefaultBuildURL(path string, s *FileServer) string {
	keys := strings.Split(path, "/")
	if len(keys) == 0 || len(keys) == 1 {
		return s.Buster
	}
	if len(keys) <= 2 {
		return path + "/" + s.Buster
	}

	keys = append(keys, "")
	copy(keys[3:], keys[2:])
	keys[2] = s.Buster

	if len(s.Host) > 0 {
		return s.Host + strings.Join(keys, "/")
	}
	return strings.Join(keys, "/")
}

// GinFunc returns a gin middleware func
func (s *FileServer) GinFunc() gin.HandlerFunc {
	s.prefix = s.BuildBusterFunc(s)
	return func(ctx *gin.Context) {
		if !strings.HasPrefix(ctx.Request.URL.Path, s.Prefix) {
			ctx.Next()
			return
		}

		filePath, ok := s.StripBusterFunc(ctx.Request.URL.Path, s)
		if !ok {
			ctx.Next()
			return
		}

		stats, err := os.Stat(filePath)
		if err != nil {
			ctx.Next()
			return
		}
		if stats.IsDir() {
			ctx.Next()
			return
		}

		if s.DisableCache {
			ctx.Writer.Header().Set("Cache-Control", noCacheHeader)
		} else {
			ctx.Writer.Header().Set("Cache-Control", cacheHeader)
		}
		http.ServeFile(ctx.Writer, ctx.Request, filePath)
		ctx.Abort()
	}
}

// vars used to generate a random string
var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rnd         = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func random(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rnd.Intn(len(letterRunes))]
	}
	return string(b)
}
