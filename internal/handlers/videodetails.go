package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "html/template"
)

type Video struct {
    gorm.Model
    Filename     string
    Transcription string
}

func ShowVideoDetails(c *gin.Context, db *gorm.DB) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing video ID"})
        return
    }

    var video Video
    if err := db.First(&video, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }

    const html = `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Video Details</title>
    </head>
    <body>
        <h1>Video Details</h1>
        <p>Video ID: {{ .ID }}</p>
        <p>Filename: {{ .Filename }}</p>
        <p>Transcription Status: {{ .Transcription }}</p>
    </body>
    </html>
    `

    tmpl, err := template.New("videoDetails").Parse(html)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error rendering video details"})
        return
    }

    c.Header("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.Execute(c.Writer, video); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing template"})
        return
    }
}
