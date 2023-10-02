package handlers

import (
    "html/template"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

func ServePlaybackPage(c *gin.Context, db *gorm.DB, videoFileName string) {
    videoSrc := "/uploads/" + videoFileName

    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Video Playback</title>
    </head>
    <body>
        <video controls>
            <source src="{{ .VideoSrc }}" type="video/mp4">
            Your browser does not support the video tag.
        </video>
    </body>
    </html>
    `

    tmpl, err := template.New("playback").Parse(html)
    if err != nil {
        http.Error(c.Writer, "Error rendering playback page", http.StatusInternalServerError)
        return
    }

    data := struct {
        VideoSrc string
    }{
        VideoSrc: videoSrc,
    }

    err = tmpl.Execute(c.Writer, data)
    if err != nil {
        http.Error(c.Writer, "Error executing page", http.StatusInternalServerError)
        return
    }
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
