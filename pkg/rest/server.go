package rest

import (
	"fmt"
	"github.com/joern1811/ai/pkg/core/domain"
	"github.com/joern1811/ai/pkg/core/ports"
	"github.com/joern1811/ai/pkg/core/service"
	adapters "github.com/joern1811/ai/pkg/framework/adapters/notifiers"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"log"
	"net/http"
	"os"
)

type Server struct {
	speachService *service.SpeachService
	notifyAdapter ports.Notifier
}

// NewServer erstellt einen neuen Server
func NewServer() *Server {
	return &Server{
		speachService: service.NewSpeachService(os.Getenv("OPEN_AI_AUTH_TOKEN"), domain.PromptConfig{
			SummarizePrompt: "Fasse die folgende Nachricht stichpunktartig zusammen. Wenn du Aufgaben identifizieren kannst, erstelle einen extra Bereich dafür.",
		}),
		notifyAdapter: adapters.NewTelegramNotifier(adapters.TelegramConfig{
			ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
			Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		}),
	}
}

func (srv Server) Start() {
	http.HandleFunc("/upload/", srv.uploadHandler)

	fmt.Println("Server startet auf Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// uploadHandler verarbeitet den Datei-Upload
func (srv Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob es ein POST-Request ist
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT methods are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Überprüfen des Authorization-Headers
	apiToken := os.Getenv("API_TOKEN")
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+apiToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extrahieren des Dateinamens aus der URL
	fileName := r.URL.Path[len("/upload/"):]
	log.Printf("Datei %s wird hochgeladen...", fileName)

	// Datei aus dem Request extrahieren
	filePath := "./uploads/" + fileName
	log.Printf("Datei wird gespeichert unter %s", filePath)
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Konnte Datei nicht speichern", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Kopiert den Dateiinhalt
	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		http.Error(w, "Fehler beim Speichern der Datei", http.StatusInternalServerError)
		return
	}

	// async process summary
	go srv.processSummary(filePath)

	// Erfolgsnachricht zurücksenden
	_, _ = fmt.Fprintf(w, "OK")
}

func (srv Server) processSummary(filePath string) {
	log.Printf("Zusammenfassung für Datei %s wird erstellt...", filePath)
	summary, err := srv.speachService.SummarizeAudio(filePath)
	if err != nil {
		log.Printf("Fehler beim Erstellen der Zusammenfassung: %s", err)
		return
	}
	err = srv.notifyAdapter.Notify(summary)
	if err != nil {
		log.Printf("Fehler beim Senden der Benachrichtigung: %s", err)
		return
	}
}
