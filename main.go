package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func downloadVideo(url string, randomStr string) (string, error) {
	if url == "" {
		url = "https://www.instagram.com/reel/DJgyrt0gxHc/"
	}
	cmd := exec.Command("yt-dlp", url, "-o", fmt.Sprintf("out/%s.mp4", randomStr))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error downloading video: %v, output: %s", err, output)
	}
	return string(output), nil
}

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func loadEnv() {
	// parse .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		file, err := os.ReadFile(".env")
		if err != nil {
			log.Fatalf("Error reading .env file: %v", err)
		}
		lines := strings.Split(string(file), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) > 0 && !strings.HasPrefix(line, "#") { // skip empty lines and comments
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					os.Setenv(key, value)
				}
			}
		}
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	loadEnv()

	password := getEnv("PASSWORD", "admin")
	log.Printf("Using password %s\n", password)

	port := getEnv("PORT", "8080")
	log.Printf("Using port %s\n", port)

	host := getEnv("HOST", "127.0.0.1")
	log.Printf("Using host %s\n", host)

	// Serve static files under /static/
	staticFS, _ := fs.Sub(staticFiles, "static")
	http.Handle("/static/", loggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	// serve video files under /out/
	http.Handle("/out/", loggingMiddleware(http.StripPrefix("/out/", http.FileServer(http.Dir("out")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		tmplData, err := templateFiles.ReadFile("templates/index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.New("index").Parse(string(tmplData))
		if err != nil {
			http.Error(w, "template parse error", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "template execute error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseMultipartForm(10 << 20) // 10 MB limit
		url := r.FormValue("url")
		inputPassword := r.FormValue("password")
		if inputPassword != password {
			log.Printf("Unauthorized download attempt with password: %s\n", inputPassword)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		randomStr := randomString(10)
		output, err := downloadVideo(url, randomStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error downloading video: %v", err), http.StatusInternalServerError)
			return
		}
		response := fmt.Sprintf("Video downloaded successfully. Output: %s, File: %s.mp4", output, randomStr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		type jsonResponse struct {
			Message string `json:"message"`
			File    string `json:"file"`
		}
		resp := jsonResponse{
			Message: response,
			File:    randomStr + ".mp4",
		}
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseMultipartForm(10 << 20)
		inputPassword := r.FormValue("password")
		if inputPassword != password {
			log.Printf("Unauthorized clear attempt with password: %s\n", inputPassword)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		files, err := fs.Glob(os.DirFS("."), "out/*")
		if err != nil {
			http.Error(w, "Error listing files", http.StatusInternalServerError)
			return
		}
		for _, file := range files {
			err := os.Remove(file)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error deleting file %s: %v", file, err), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "All files deleted successfully")
	})

	// login route with just password check
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseMultipartForm(10 << 20)
		inputPassword := r.FormValue("password")
		if inputPassword != password {
			log.Printf("Unauthorized login attempt with password: %s, doesn't match %s", inputPassword, password)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If password matches, send a success response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Login successful")
	})

	log.Printf("Starting server on http://%s:%s\n", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
