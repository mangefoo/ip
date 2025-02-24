package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type PageData struct {
	IP        string
	Timestamp string
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func ipHandlerHTML(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	currentTime := time.Now().Format("2006-01-02 15:04:05 MST")

	data := PageData{
		IP:        clientIP,
		Timestamp: currentTime,
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your IP Address</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background-color: #f0f0f0;
        }
        .container {
            background-color: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            text-align: center;
        }
        h1 {
            color: #333;
        }
        .ip {
            font-size: 1.5rem;
            color: #0066cc;
            margin: 1rem 0;
        }
        .timestamp {
            font-size: 0.9rem;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Your IP Address</h1>
        <p class="ip">{{.IP}}</p>
        <p class="timestamp">Accessed on: {{.Timestamp}}</p>
    </div>
</body>
</html>
`

	t, err := template.New("ippage").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ipHandlerPlainText(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, clientIP)
}

func main() {
	http.HandleFunc("/ip", ipHandlerHTML)
	http.HandleFunc("/ip/plain", ipHandlerPlainText)
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
