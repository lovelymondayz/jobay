package config

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	// "time"

	// _ "github.com/lib/pq"
	"github.com/unpoller/unifi"
)

var UnifiClient *unifi.Unifi

var UnifiClient2 *http.Client

// // InitDB initializes the database connection
// func InitDB() {
// 	var err error
// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		dbURL = "postgres://wifiuser:wifipass@localhost/wifitracking?sslmode=disable"
// 	}

// 	DB, err = sql.Open("postgres", dbURL)
// 	if err != nil {
// 		log.Fatal("Failed to open database:", err)
// 	}

// 	// Set connection pool settings
// 	DB.SetMaxOpenConns(25)
// 	DB.SetMaxIdleConns(25)
// 	DB.SetConnMaxLifetime(5 * time.Minute)

// 	if err = DB.Ping(); err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	log.Println("Database connection established")
// }

// // CloseDB closes the database connection
// func CloseDB() {
// 	if DB != nil {
// 		DB.Close()
// 	}
// }

// InitUniFiClient initializes the UniFi client
func InitUniFiClient() {
	controllerURL := os.Getenv("UNIFI_CONTROLLER_URL")
	username := os.Getenv("UNIFI_USERNAME")
	password := os.Getenv("UNIFI_PASSWORD")

	if controllerURL == "" || username == "" || password == "" {
		log.Fatal("UniFi credentials not provided. Set UNIFI_CONTROLLER_URL, UNIFI_USERNAME, and UNIFI_PASSWORD")
	}

	config := &unifi.Config{
		User:      username,
		Pass:      password,
		URL:       controllerURL,
		VerifySSL: false,
	}

	var err error
	UnifiClient, err = unifi.NewUnifi(config)
	if err != nil {
		log.Fatal("Failed to create UniFi client:", err)
	}

	log.Println("UniFi client initialized")
}

func InitUniFiClient2() {
	controllerURL := os.Getenv("UNIFI_CONTROLLER_URL")
	username := os.Getenv("UNIFI_USERNAME")
	password := os.Getenv("UNIFI_PASSWORD")

	if controllerURL == "" || username == "" || password == "" {
		log.Fatal("UniFi credentials not provided. Set UNIFI_CONTROLLER_URL, UNIFI_USERNAME, and UNIFI_PASSWORD")
	}

	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	UnifiClient2 = &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	loginURL := fmt.Sprintf("%s/api/login", controllerURL)
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Failed to marshal login payload:", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Failed to create login request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := UnifiClient2.Do(req)
	if err != nil {
		log.Fatal("Login request failed:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatal("Login failed with status", resp.StatusCode, ":", string(bodyBytes))
	}
}
