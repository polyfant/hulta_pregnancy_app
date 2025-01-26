package middleware

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

func StaticFileMiddleware(frontendBuildPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if requesting a static file
		requestPath := c.Request.URL.Path
		
		// Strip leading slash
		if len(requestPath) > 0 && requestPath[0] == '/' {
			requestPath = requestPath[1:]
		}

		filePath := filepath.Join(frontendBuildPath, requestPath)

		// Check if file exists
		if _, err := os.Stat(filePath); err == nil {
			// Serve static file
			c.File(filePath)
			c.Abort()
			return
		}

		// Fallback to index.html for SPA routing
		indexPath := filepath.Join(frontendBuildPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
			c.Abort()
			return
		}

		// Log missing file attempts
		logger.Info("Static file not found", 
			"path", requestPath, 
			"frontend_build_path", frontendBuildPath)

		c.Next()
	}
}
