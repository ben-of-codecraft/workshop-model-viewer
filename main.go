package main

import (
	"crypto/tls"
	"compress/gzip"
	"compress/flate"
	"context"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"


	"github.com/ben-of-codecraft/workshop-model-viewer/items"
)
//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var client = &http.Client{
	Timeout: 10 * time.Second, // Set a timeout for requests
}

// Define a semaphore to limit concurrent requests
var semaphore = make(chan struct{}, 10) // Limit to 10 concurrent requests
var wg sync.WaitGroup                // WaitGroup to wait for goroutines to finish


func main() {

	// Show the web application view and assets 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html.tmpl", nil)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Add APIs used for the webapplication
	http.HandleFunc("/item-lookup", ItemLookUpHandler)
	http.HandleFunc("/get-races", GetRacesHandler)
	http.HandleFunc("/broken", BrokenHandler)
	http.HandleFunc("/proxy/", ProxyHandler)

	// Start the server
	dev := os.Getenv("DEVELOPMENT")
	port := "8080"
	if dev == "false" {

		// Start listening on port 80/443 traffic for production build 
		go func() {
			log.Println("listening on", "80")
			log.Fatal(http.ListenAndServe(":80", nil))
		}()

		// Self-signed certificate jsut for testing
		certFile := "cert.pem" 
		keyFile := "key.pem"   
		httpsServer := &http.Server{
			Addr: ":443",
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	
		log.Println("Listening on port 443")
		log.Fatal(httpsServer.ListenAndServeTLS(certFile, keyFile))
	} else {	
		log.Println("listening on", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

}

func ItemLookUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get itemId from the query string
	qItem := r.URL.Query().Get("item")
	if qItem == "" {
		http.Error(w, "item is required", http.StatusBadRequest)
		return
	}

	// conver to itemId to an int	
	itemId, err := strconv.Atoi(qItem)
	if err != nil {
		http.Error(w, "item must be a number", http.StatusBadRequest)
		return
	}


	displayId, err := items.GetDisplayId(itemId)
	if err != nil {
		http.Error(w, "Error fetching display ID", http.StatusInternalServerError)
		return
	}

	// return the display id in a json response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"displayId": "` + displayId + `"}`))
}

func GetRacesHandler(w http.ResponseWriter, r *http.Request) {

	raceIDs := map[string]int{
		"Human":                 1,
		"Orc":                   2,
		"Dwarf":                 3,
		"Night Elf":             4,
		"Undead":                5,
		"Tauren":                6,
		"Gnome":                 7,
		"Troll":                 8,
		"Goblin":                9,
		"Blood Elf":             10,
		"Draenei":               11,
		"Fel Orc":               12,
		"Naga":                  13,
		"Broken":                14,
		"Skeleton":              15,
		"Vrykul":                16,
		"Tuskarr":               17,
		"Forest Troll":          18,
		"Taunka":                19,
		"Northrend Skeleton":    20,
		"Ice Troll":             21,
		"Worgen":                22,
		"Gilnean":               23,
		"Pandaren (Neutral)":    24,
		"Pandaren (Alliance)":   25,
		"Pandaren (Horde)":      26,
		"Nightborne":            27,
		"Highmountain Tauren":   28,
		"Void Elf":              29,
		"Lightforged Draenei":   30,
	}

	// Marshal the map into JSON
	jsonData, err := json.Marshal(raceIDs)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func BrokenHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "something terrible has happened", http.StatusInternalServerError)
}
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Build the target URL
	targetURL := "https://wow.zamimg.com/modelviewer/live/" + r.URL.Path[len("/proxy/"):]

	// Create a new request to the target URL
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Forward original request headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set Accept header to request JSON
	req.Header.Set("Accept", "application/json")

	// Use context to handle timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Limit concurrent requests
	semaphore <- struct{}{} // Acquire a token
	defer func() {
		<-semaphore // Release the token
	}()

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching data from external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	// Determine the appropriate reader based on Content-Encoding
	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			http.Error(w, "Error creating gzip reader", http.StatusInternalServerError)
			return
		}
		defer gzipReader.Close()
		reader = gzipReader
	case "deflate":
		flateReader := flate.NewReader(resp.Body)
		defer flateReader.Close()
		reader = flateReader
	}

	// Stream the response directly to the client
	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "Error copying response to client", http.StatusInternalServerError)
		return
	}
}