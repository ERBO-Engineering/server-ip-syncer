package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type Config struct {
	FirebaseConfigPath   string `json:"firebase_config_path"`
	WireGuardConfigPath  string `json:"wireguard_config_path"`
	FirebaseDocumentPath string `json:"firebase_document_path"`
}

type IPData struct {
	IP string `json:"ip"`
}

func main() {
	// Read configuration
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Initialize Firebase
	opt := option.WithCredentialsFile(config.FirebaseConfigPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	defer firestore.Close()

	log.Println("Starting IP monitoring service...")

	for {
		// Read current WireGuard IP
		currentIP, err := getWireGuardIP(config.WireGuardConfigPath)
		if err != nil {
			log.Printf("Error reading WireGuard IP: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		// Get IP from Firebase
		firebaseIP, err := getFirebaseIP(firestore, config.FirebaseDocumentPath)
		if err != nil {
			log.Printf("Error reading Firebase IP: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		// Compare IPs and update if necessary
		if currentIP != firebaseIP {
			log.Printf("IP mismatch detected. Current: %s, Firebase: %s", currentIP, firebaseIP)
			if err := updateWireGuardConfig(config.WireGuardConfigPath, firebaseIP); err != nil {
				log.Printf("Error updating WireGuard config: %v", err)
				time.Sleep(5 * time.Minute)
				continue
			}
			if err := restartWireGuard(); err != nil {
				log.Printf("Error restarting WireGuard: %v", err)
				time.Sleep(5 * time.Minute)
				continue
			}
			log.Printf("Successfully updated WireGuard IP to %s", firebaseIP)
		}

		time.Sleep(5 * time.Minute)
	}
}

func getWireGuardIP(configPath string) (string, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Endpoint = ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				// Return just the IP part without the port
				return strings.Split(parts[2], ":")[0], nil
			}
		}
	}
	return "", fmt.Errorf("endpoint IP not found in WireGuard config")
}

func getFirebaseIP(firestore *firestore.Client, docPath string) (string, error) {
	log.Printf("Reading IP from Firestore document: %s", docPath)
	doc, err := firestore.Doc(docPath).Get(context.Background())
	if err != nil {
		return "", fmt.Errorf("error reading Firestore document: %v", err)
	}

	var ipData IPData
	if err := doc.DataTo(&ipData); err != nil {
		return "", fmt.Errorf("error parsing Firestore data: %v", err)
	}

	if ipData.IP == "" {
		return "", fmt.Errorf("no IP found in Firestore document")
	}

	log.Printf("Successfully read IP from Firestore: %s", ipData.IP)
	return ipData.IP, nil
}

func updateWireGuardConfig(configPath string, newIP string) error {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.Contains(line, "Endpoint = ") {
			// Keep the port number from the existing endpoint
			port := ":51820" // Default port
			if strings.Contains(line, ":") {
				port = line[strings.LastIndex(line, ":"):]
			}
			newLines = append(newLines, fmt.Sprintf("Endpoint = %s%s", newIP, port))
		} else {
			newLines = append(newLines, line)
		}
	}

	return ioutil.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func restartWireGuard() error {
	// Stop WireGuard
	stopCmd := exec.Command("sudo", "wg-quick", "down", "wg0")
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("error stopping WireGuard: %v", err)
	}

	// Start WireGuard
	startCmd := exec.Command("sudo", "wg-quick", "up", "wg0")
	if err := startCmd.Run(); err != nil {
		return fmt.Errorf("error starting WireGuard: %v", err)
	}

	return nil
}
