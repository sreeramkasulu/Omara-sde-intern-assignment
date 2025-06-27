package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Server struct {
	db *sql.DB
	geminiClient *genai.GenerativeModel
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Document struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FileName    string    `json:"file_name"`
	StoragePath string    `json:"storage_path"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type DocumentChunk struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	ChunkIndex int    `json:"chunk_index"`
	Content    string `json:"content"`
}

type ChatMessage struct {
	ID             string    `json:"id"`
	DocumentID     string    `json:"document_id"`
	UserID         string    `json:"user_id"`
	MessageType    string    `json:"message_type"`
	MessageContent string    `json:"message_content"`
	Timestamp      time.Time `json:"timestamp"`
}

type AnalyzeRequest struct {
	Query string `json:"query"`
}

type AnalyzeResponse struct {
	Response string `json:"response"`
}

func main() {
	// Database connection
	dbURL := "postgres://postgres:password@localhost/strategic_insight_db?sslmode=disable"
	if os.Getenv("DATABASE_URL") != "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Gemini API client setup
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		log.Fatal(err)
	}
	defer geminiClient.Close()

	model := geminiClient.GenerativeModel("gemini-1.5-flash-latest") // Using gemini-1.5-flash-latest for text generation

	server := &Server{db: db, geminiClient: model}

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Document routes
	api.HandleFunc("/documents/upload", server.uploadDocument).Methods("POST")
	api.HandleFunc("/documents", server.listDocuments).Methods("GET")
	api.HandleFunc("/documents/{id}", server.getDocument).Methods("GET")
	api.HandleFunc("/documents/{id}", server.deleteDocument).Methods("DELETE")
	api.HandleFunc("/documents/{id}/analyze", server.analyzeDocument).Methods("POST")
	api.HandleFunc("/documents/{id}/chat-history", server.getChatHistory).Methods("GET")

	// User routes
	api.HandleFunc("/users", server.createUser).Methods("POST")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	// Create uploads directory
	os.MkdirAll("uploads", 0755)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (id, email, created_at) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, user.ID, user.Email, user.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Server) uploadDocument(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Validate file type
	fileName := header.Filename
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".pdf" && ext != ".txt" && ext != ".docx" {
		http.Error(w, "Only PDF, TXT, and DOCX files are allowed", http.StatusBadRequest)
		return
	}

	// Generate unique document ID
	docID := uuid.New().String()
	storagePath := filepath.Join("uploads", docID+ext)

	// Save file
	dst, err := os.Create(storagePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Extract text content
	content, err := extractTextFromFile(storagePath, ext)
	if err != nil {
		http.Error(w, "Failed to extract text", http.StatusInternalServerError)
		return
	}

	// Save document to database
	doc := Document{
		ID:          docID,
		UserID:      userID,
		FileName:    fileName,
		StoragePath: storagePath,
		UploadedAt:  time.Now(),
	}

	query := `INSERT INTO documents (id, user_id, file_name, storage_path, uploaded_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = s.db.Exec(query, doc.ID, doc.UserID, doc.FileName, doc.StoragePath, doc.UploadedAt)
	if err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	// Chunk the content and save to database
	chunks := chunkText(content, 1000) // 1000 characters per chunk
	for i, chunk := range chunks {
		chunkID := uuid.New().String()
		query := `INSERT INTO document_chunks (id, document_id, chunk_index, content) VALUES ($1, $2, $3, $4)`
		_, err = s.db.Exec(query, chunkID, docID, i, chunk)
		if err != nil {
			log.Printf("Failed to save chunk %d: %v", i, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, user_id, file_name, storage_path, uploaded_at FROM documents WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(&doc.ID, &doc.UserID, &doc.FileName, &doc.StoragePath, &doc.UploadedAt)
		if err != nil {
			log.Printf("Failed to scan document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	query := `SELECT id, user_id, file_name, storage_path, uploaded_at FROM documents WHERE id = $1`
	var doc Document
	err := s.db.QueryRow(query, docID).Scan(&doc.ID, &doc.UserID, &doc.FileName, &doc.StoragePath, &doc.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch document", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	// Get document info first
	var storagePath string
	query := `SELECT storage_path FROM documents WHERE id = $1`
	err := s.db.QueryRow(query, docID).Scan(&storagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch document", http.StatusInternalServerError)
		}
		return
	}

	// Delete from database (cascades to chunks and chat history)
	query = `DELETE FROM documents WHERE id = $1`
	_, err = s.db.Exec(query, docID)
	if err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	// Delete file from storage
	os.Remove(storagePath)

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) analyzeDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Get document chunks
	query := `SELECT content FROM document_chunks WHERE document_id = $1 ORDER BY chunk_index`
	rows, err := s.db.Query(query, docID)
	if err != nil {
		http.Error(w, "Failed to fetch document content", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var content strings.Builder
	for rows.Next() {
		var chunk string
		if err := rows.Scan(&chunk); err != nil {
			continue
		}
		content.WriteString(chunk)
		content.WriteString("\n")
	}

	// Call LLM API
	geminiResponse, err := s.generateInsight(req.Query, content.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate insight: %v", err), http.StatusInternalServerError)
		return
	}

	// Save chat history
	userMsgID := uuid.New().String()
	aiMsgID := uuid.New().String()
	now := time.Now()

	// Save user message
	query = `INSERT INTO chat_history (id, document_id, user_id, message_type, message_content, timestamp) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = s.db.Exec(query, userMsgID, docID, userID, "user", req.Query, now)
	if err != nil {
		log.Printf("Failed to save user chat message: %v", err)
	}

	// Save AI response
	_, err = s.db.Exec(query, aiMsgID, docID, userID, "ai", geminiResponse, now)
	if err != nil {
		log.Printf("Failed to save AI chat message: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AnalyzeResponse{Response: geminiResponse})
}

func (s *Server) getChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	query := `SELECT id, document_id, user_id, message_type, message_content, timestamp FROM chat_history WHERE document_id = $1 ORDER BY timestamp`
	rows, err := s.db.Query(query, docID)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		err := rows.Scan(&msg.ID, &msg.DocumentID, &msg.UserID, &msg.MessageType, &msg.MessageContent, &msg.Timestamp)
		if err != nil {
			log.Printf("Failed to scan message: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func extractTextFromFile(filePath, ext string) (string, error) {
	if ext == ".txt" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	} else if ext == ".pdf" {
		cmd := exec.Command("pdftotext", filePath, "-")
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to extract text from PDF: %w", err)
		}
		return string(output), nil
	} else if ext == ".docx" {
		// Placeholder for DOCX extraction
		return "DOCX content extraction not implemented yet", nil
	}
	return "", fmt.Errorf("unsupported file type: %s", ext)
}

func chunkText(text string, chunkSize int) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

func (s *Server) generateInsight(query, documentContent string) (string, error) {
	ctx := context.Background()

	prompt := []genai.Part{
		genai.Text(fmt.Sprintf("Based ONLY on the following document content, answer the query. If the information is not available in the document, state that. \n\nDocument Content:\n%s\n\nQuery: %s", documentContent, query)),
	}

	resp, err := s.geminiClient.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", fmt.Errorf("Gemini API call failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "No insight generated by AI.", nil
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			sb.WriteString(string(text))
		}
	}
	return sb.String(), nil
}



