# IP Viber

A Go application that monitors Firebase for IP changes and automatically updates WireGuard configuration.

## Prerequisites

- Go 1.21 or later
- Firebase project with Firestore database
- WireGuard installed and configured
- sudo privileges for WireGuard operations

## Setup

1. Install the required dependencies:
```bash
go mod tidy
```

2. Create a Firebase project and download the service account credentials:
   - Go to Firebase Console
   - Create a new project or select an existing one
   - Go to Project Settings > Service Accounts
   - Generate a new private key
   - Save the downloaded JSON file as `firestore.json` in the project root
   - Note: This file is gitignored for security reasons

3. Configure the application:
   - Edit `config.json` to set your specific paths:
     - `firebase_config_path`: Path to your Firebase credentials file
     - `wireguard_config_path`: Path to your WireGuard configuration file
     - `firebase_document_path`: Path to the Firestore document containing the IP

4. Set up the Firestore document:
   - Create a document in your Firestore database
   - Add a field named `ip` with the current WireGuard IP address

## Building and Running

1. Build the application:
```bash
go build
```

2. Run the application with sudo privileges:
```bash
sudo ./IP_viber
```

The application will:
- Monitor the Firebase document every 5 minutes
- Compare the IP in Firebase with the current WireGuard configuration
- Update the WireGuard configuration and restart the service if a mismatch is detected

## Installation as a Service

To install the application as a systemd service:

1. Make sure you have all required files:
   - `IP_viber` (compiled binary)
   - `config.json`
   - `firestore.json`
   - `ip-viber.service`
   - `install.sh`

2. Run the installation script:
```bash
./install.sh
```

After installation, you can:
- Check the service status: `sudo systemctl status ip-viber`
- View logs: `sudo journalctl -u ip-viber -f`
- Stop the service: `sudo systemctl stop ip-viber`
- Start the service: `sudo systemctl start ip-viber`
- Restart the service: `sudo systemctl restart ip-viber`

## Git Setup

This project uses git for version control. The following files are ignored:
- `firestore.json` (Firebase credentials)
- `IP_viber` (compiled binary)
- IDE and OS-specific files

To set up git:
```bash
git init
git add .
git commit -m "Initial commit"
```

## Security Considerations

- The application requires sudo privileges to modify WireGuard configuration
- Keep your Firebase credentials secure and never commit them to version control
- Consider using environment variables for sensitive configuration 