package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "encoding/json"
    "os"
    supabase "github.com/supabase-community/supabase-go"
    "github.com/joho/godotenv"
)

type Reading struct {
	CO2 uint16 `json:"co2" valid:"required"`
    Timestamp time.Time `json:"timestamp" valid:"required"`
}

var supabaseTable string = "readings"

// Global supabase client that is passed to endpoints
var supabaseClient *supabase.Client

func getSupabaseConnection() *supabase.Client {

    if supabaseClient != nil {
        return supabaseClient
    }

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    API_URL := os.Getenv("API_URL")
    API_KEY := os.Getenv("API_KEY")
    fmt.Println("Initializing Supabase Connection")


    client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})
    if err != nil {
        log.Fatal("Failed to initalize the client: ", err)
    }

    supabaseClient = client
    return supabaseClient
}


func recordReading(w http.ResponseWriter, r *http.Request) {
    /*
    Receives CO2 readings from home server
    */

    if r.Method != http.MethodPost {
        fmt.Println("Non POST request received")
        http.Error(w, "Non POST request received", http.StatusMethodNotAllowed)
        return
    }

    // Check if json content type is being sent 
    if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
        return
    }

    var reading Reading;
    err := json.NewDecoder(r.Body).Decode(&reading)

    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("CO2: %d, Time: %s\n", reading.CO2, reading.Timestamp)

    // Write Reading to Supabase
    var client *supabase.Client = getSupabaseConnection()
    data, _, err := client.From(supabaseTable).Insert(reading, false, "", "", "").Execute()

    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("Successfully inserted reading %s", string(data))

    // Sets 201 status
    w.WriteHeader(http.StatusCreated)
}


func getReadings(w http.ResponseWriter, r *http.Request) {
    /*
    Retrieve readings from supabase for frontend
    */

    // Retrieve readings from supabase
    var client *supabase.Client = getSupabaseConnection()

    data, _, err := client.From(supabaseTable).Select("*", "", false).Execute()

    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func handleRequests() {
    http.HandleFunc("POST /reading", recordReading)
    http.HandleFunc("GET /reading", getReadings)
    http.Handle("/", http.FileServer(http.Dir("./static")))

    fmt.Println("Listening on Port :10000")
    log.Fatal(http.ListenAndServe(":10000", nil))
}


func main() {
    handleRequests()
}