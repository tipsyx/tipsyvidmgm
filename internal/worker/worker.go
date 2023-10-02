package worker

import (
    "log"

    "github.com/jinzhu/gorm"
    "github.com/streadway/amqp"
)

type TranscriptionWorker struct {
    Channel *amqp.Channel
    DB      *gorm.DB
}

func NewTranscriptionWorker(ch *amqp.Channel, db *gorm.DB) *TranscriptionWorker {
    return &TranscriptionWorker{Channel: ch, DB: db}
}

func (w *TranscriptionWorker) Start() {
    msgs, err := w.Channel.Consume(
        "transcription_queue",
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            filename := string(d.Body)
            transcription, err := transcribeAudio(filename, w.Channel)
            if err != nil {
                log.Printf("Error transcribing audio: %v", err)
            } else {
                err = w.DB.Table("videos").Where("filename = ?", filename).Update("transcription", transcription).Error
                if err != nil {
                    log.Printf("Error updating database with transcription result: %v", err)
                }
            }
        }
    }()

    log.Printf("Transcription worker is waiting for tasks. To exit, press Ctrl+C")
}
