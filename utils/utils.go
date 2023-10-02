package utils

import (
    "github.com/gerow/thumbs"
    "io"
    "net/http"
)

func GenerateVideoThumbnail(videoPath string, thumbnailPath string) error {
    config := thumbs.Config{
        SeekPercentage: 10,
        Width:          320,
        Height:         240,
        Output:         thumbnailPath,
    }

    err := thumbs.Create(videoPath, config)
    if err != nil {
        return err
    }

    return nil
}

func validateScreenRecording(file io.Reader, filename string) bool {
    ext := filepath.Ext(filename)
    switch ext {
    case ".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv":
        return true
    } 
    
    buffer := make([]byte, 512)
    _, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return false
    }
    mimeType := http.DetectContentType(buffer)
    switch mimeType {
    case "video/mp4", "video/avi", "video/quicktime", "video/x-matroska":
        return true
    default:
        return false
    }
}
