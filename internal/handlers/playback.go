package handlers

import (
    "html/template"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

func PlaybackPage(c *gin.Context, db *gorm.DB) {
    id := c.Param("id")
    if id == ""{
        c.JSON(http.StatusBadRequest, gin.H{"error":"Video not found"})
        return
    }
    videoSrc := filepath.Join("/uploads/",id+".mp4")

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

    playbackTemplate, err := template.New("playback").Parse(html)
    if err != nil {
        http.Error(c.Writer, "Error rendering playback page", http.StatusInternalServerError)
        return
    }

    data := struct {
        VideoSrc string
    }{
        VideoSrc: videoSrc,
    }

    err = playbackTemplate.Execute(c.Writer, data)
    if err != nil {
        http.Error(c.Writer, "Error executing page", http.StatusInternalServerError)
        return
    }
      c.HTML(http.StatusOK, "playback", data)
}
