package main

import (
	"crypto/tls"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"encoding/json"

	"github.com/ben-of-codecraft/workshop-model-viewer/items"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

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