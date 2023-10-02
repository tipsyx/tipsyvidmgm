package transcriber

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/streadway/amqp"
)

func TranscribeAudio(filename string, ch *amqp.Channel) (string, error) {
    audio, err := ioutil.ReadFile(filename)
    if err != nil {
        return "", err
    }

    apiKey := "YOUR_WHISPER_API_KEY"

    req, err := http.NewRequest("POST", "https://api.whisper.ai/v1/recognize", bytes.NewReader(audio))
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "audio/wav")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Whisper API returned status code %d", resp.StatusCode)
    }

    transcription, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(transcription), nil
}
