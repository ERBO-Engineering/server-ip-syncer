# IP Address Tracker

This Go application monitors your machine's public IP address and stores it in Firebase Firestore. It checks for IP changes every 5 minutes and creates a new document in the Firestore collection when changes are detected.

## Prerequisites

- Go 1.21 or later
- Firebase project with Firestore enabled
- Firebase service account key

## Setup

1. Clone this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Create a Firebase project and enable Firestore
4. Generate a service account key from Firebase Console:
   - Go to Project Settings > Service Accounts
   - Click "Generate New Private Key"
   - Save the JSON file as `key.json` in the project directory

## Running the Application

### Manual Run
```bash
go run main.go
```

### Running as a System Service (systemd)

1. Copy the service file to systemd directory:
   ```bash
   sudo cp ip-leaker.service /etc/systemd/system/
   ```

2. Replace the user placeholder in the service file:
   ```bash
   sudo sed -i "s/%i/$(whoami)/g" /etc/systemd/system/ip-leaker.service
   ```

3. Reload systemd to recognize the new service:
   ```bash
   sudo systemctl daemon-reload
   ```

4. Enable the service to start on boot:
   ```bash
   sudo systemctl enable ip-leaker
   ```

5. Start the service:
   ```bash
   sudo systemctl start ip-leaker
   ```

6. Check the service status:
   ```bash
   sudo systemctl status ip-leaker
   ```

7. View the logs:
   ```bash
   journalctl -u ip-leaker -f
   ```

## Firestore Structure

The application stores data in a collection called `ip_addresses` with the following structure:
- `ip`: string (the public IP address)
- `timestamp`: timestamp (when the IP was recorded) 