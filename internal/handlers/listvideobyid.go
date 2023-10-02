package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

func GetVideoByID(c *gin.Context, db *gorm.DB) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing video ID"})
        return
    }

    videoID, err := strconv.Atoi(id)
    if err != nil || videoID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
        return
    }

    var video Video
    if err := db.First(&video, videoID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video"})
        return
    }

    c.JSON(http.StatusOK, video)
}
