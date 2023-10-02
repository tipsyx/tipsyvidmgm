package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
)

func ListUploadedVideos(c *gin.Context, db *gorm.DB) {
    page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    itemsPerPage, err := strconv.Atoi(c.DefaultQuery("itemsPerPage", "10"))
    if err != nil || itemsPerPage <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid itemsPerPage value"})
        return
    }

    offset, limit := calculatePaginationValues(page, itemsPerPage)

    var videos []Video
    if err := db.Offset(offset).Limit(limit).Find(&videos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
        return
    }

    c.JSON(http.StatusOK, videos)
}

func calculatePaginationValues(page, itemsPerPage int) (int, int) {
    if page < 1 {
        page = 1
    }

    if itemsPerPage <= 0 {
        itemsPerPage = 10
    }

    offset := (page - 1) * itemsPerPage
    limit := itemsPerPage

    return offset, limit
}
