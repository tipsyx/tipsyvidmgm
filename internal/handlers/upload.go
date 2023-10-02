package handlers

import (
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/streadway/amqp"

    "github.com/tipsyx/tipsyvidmgm/config"
)

func HandleUpload(c *gin.Context, ch *amqp.Channel, db *gorm.DB) {
    maxUploadSize := int64(200 * 1024 * 1024) 

    file, header, err := c.Request.FormFile("video")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving the file"})
        return
    }
    defer file.Close()

    if header.Size > maxUploadSize {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit"})
        return
    }

    uploadDir := "./uploads"
    
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating upload directory"})
        return
    }

    ext := filepath.Ext(header.Filename)
    filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
    filePath := filepath.Join(uploadDir, filename)

    outFile, err := os.Create(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the file"})
        return
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying the file"})
        return
    }

    err = ch.Publish(
        "",                    
        "transcription_queue", 
        false,                
        false,                
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(filePath),
        },
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error enqueuing transcription task"})
        return
    }

    newVideo := Video{Filename: filename, Transcription: "In progress"}
    db.Create(&newVideo)

    c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": filename})
}
