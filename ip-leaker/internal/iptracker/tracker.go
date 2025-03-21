package iptracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type IPInfo struct {
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

type Tracker struct {
	client       *firestore.Client
	currentIPRef *firestore.DocumentRef
	historyRef   *firestore.DocumentRef
	lastIP       string
}

func NewTracker(credentialsFile string) (*Tracker, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	tracker := &Tracker{
		client:       client,
		currentIPRef: client.Collection("ip_addresses").Doc("current_ip"),
		historyRef:   client.Collection("ip_addresses").Doc("history"),
	}

	// Get initial current IP
	currentDoc, err := tracker.currentIPRef.Get(ctx)
	if err == nil {
		var currentInfo IPInfo
		if err := currentDoc.DataTo(&currentInfo); err == nil {
			tracker.lastIP = currentInfo.IP
		}
	}

	return tracker, nil
}

func (t *Tracker) Close() {
	if t.client != nil {
		t.client.Close()
	}
}

func (t *Tracker) getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.IP, nil
}

func (t *Tracker) Update() error {
	ip, err := t.getPublicIP()
	if err != nil {
		return fmt.Errorf("error getting public IP: %v", err)
	}

	if ip != t.lastIP {
		doc := IPInfo{
			IP:        ip,
			Timestamp: time.Now(),
		}

		// Update the current IP document
		_, err = t.currentIPRef.Set(context.Background(), doc)
		if err != nil {
			return fmt.Errorf("error updating current IP in Firestore: %v", err)
		}

		// Add to history
		_, _, err = t.historyRef.Collection("changes").Add(context.Background(), doc)
		if err != nil {
			return fmt.Errorf("error adding to history in Firestore: %v", err)
		}

		fmt.Printf("IP address changed from %s to %s at %s\n", t.lastIP, ip, doc.Timestamp.Format(time.RFC3339))
		t.lastIP = ip
	}

	return nil
}
