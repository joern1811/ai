package rest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/joern1811/ai/internal/core/domain"
	"github.com/joern1811/ai/internal/core/ports"
	"github.com/joern1811/ai/internal/core/service"
	adapters "github.com/joern1811/ai/internal/framework/adapters/notifiers"
	"github.com/joern1811/ai/internal/framework/adapters/utils"
)

const uploadDir = "./uploads"

// sanitizeForLog strips newline characters to prevent log-injection.
func sanitizeForLog(s string) string {
	return strings.NewReplacer("\n", "", "\r", "").Replace(s)
}

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
	mux := http.NewServeMux()
	mux.HandleFunc("/upload/", srv.uploadHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	fmt.Println("Server startet auf Port 8080...")
	log.Fatal(server.ListenAndServe())
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

	// Dateinamen aus der URL extrahieren — filepath.Base verhindert Path-Traversal
	// (Eingaben wie "../etc/passwd" werden auf "passwd" reduziert).
	fileName := filepath.Base(r.URL.Path[len("/upload/"):])
	if fileName == "." || fileName == "/" || fileName == "" {
		http.Error(w, "Ungültiger Dateiname", http.StatusBadRequest)
		return
	}
	log.Printf("Datei %s wird hochgeladen...", sanitizeForLog(fileName)) //nolint:gosec // sanitizeForLog strips newlines

	filePath := filepath.Join(uploadDir, fileName)
	log.Printf("Datei wird gespeichert unter %s", sanitizeForLog(filePath)) //nolint:gosec // sanitizeForLog strips newlines
	outFile, err := os.Create(filePath)                                     //nolint:gosec // filePath is constrained to uploadDir via filepath.Base
	if err != nil {
		log.Println(err)
		http.Error(w, "Konnte Datei nicht speichern", http.StatusInternalServerError)
		return
	}
	defer utils.CloseResource(outFile, &err)

	if _, err = io.Copy(outFile, r.Body); err != nil {
		http.Error(w, "Fehler beim Speichern der Datei", http.StatusInternalServerError)
		return
	}

	go srv.processSummary(filePath)

	_, _ = fmt.Fprintf(w, "OK")
}

func (srv Server) processSummary(filePath string) {
	log.Printf("Zusammenfassung für Datei %s wird erstellt...", sanitizeForLog(filePath)) //nolint:gosec // sanitizeForLog strips newlines
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
