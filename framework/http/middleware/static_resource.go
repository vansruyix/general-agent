package middleware

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed all:frontend/dist
var frontendFiles embed.FS

type StaticResource struct {
	HttpFS     http.FileSystem
	FileServer http.Handler
}

func NewStaticResource() *StaticResource {
	staticFS, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		panic(err)
	}
	httpFS := http.FS(staticFS)

	fileServer := http.FileServer(httpFS)
	return &StaticResource{HttpFS: httpFS, FileServer: fileServer}
}

func StaticResourceInit(sr *StaticResource) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		// API 路径放行
		if len(p) >= 4 && p[:4] == "/api" {
			c.Next()
			return
		}
		if p == "/" || p == "" {
			p = "index.html"
		}
		// 尝试提供静态文件
		// fullPath := filepath.Join(distDir, filepath.Clean(p))
		file, err := sr.HttpFS.Open(p)
		if err != nil {
			c.Next()
			return
		}

		if info, err := file.Stat(); err == nil && !info.IsDir() {
			sr.FileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		// SPA fallback
		// indexPath := filepath.Join(distDir, "index.html")
		// if _, err := os.Stat(indexPath); err == nil {
		// 	c.File(indexPath)
		// 	c.Abort()
		// 	return
		// }
		c.Next()
		// path := ctx.Request.URL.Path
		// // API 路径放行
		// if len(path) >= 4 && path[:4] == "/api" {
		// 	ctx.Next()
		// 	return
		// }
		// if path == "/" || path == "" {
		// 	path = "index.html"
		// }

		// // 如果是真实存在的静态文件，直接返回
		// file, err := sr.httpFS.Open(path)
		// if err == nil {
		// 	info, statErr := file.Stat()
		// 	_ = file.Close()

		// 	if statErr == nil && !info.IsDir() {
		// 		ctx.FileFromFS(path, sr.httpFS)
		// 		ctx.Abort()
		// 		return
		// 	}
		// }

		// // Vue Router / React Router history 模式需要 fallback 到 index.html
		// ctx.FileFromFS("index.html", sr.httpFS)
		// ctx.Abort()
	}
}
