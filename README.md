# 🌐 Server IP Syncer - Dynamic IP? No Problem! 

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.20)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Never lose connection to your WireGuard VPN server again, even with dynamic IPs! This solution automatically keeps your WireGuard configurations up-to-date by syncing server IP changes through Firebase.

## 🚀 Overview

The Server IP Syncer consists of two main components that work together to ensure your WireGuard VPN connection remains stable even when your server's IP address changes:

### 1. IP Leaker (Server-side)
- 🔍 Continuously monitors the server's public IP address
- 🔄 Detects any IP address changes
- 🔥 Pushes updates to Firebase Firestore in real-time
- 🛡️ Runs as a system service for reliability

### 2. IP Viber (Client-side)
- ⏰ Checks for IP updates every 5 minutes
- 🔄 Automatically updates WireGuard configuration when server IP changes
- 🚦 Ensures VPN connectivity remains stable
- 🛡️ Runs as a system service in the background

> 💫 **Fun Fact**: This entire codebase is 100% "vibecoded" - because sometimes the best code just vibes! 🎵

## ��️ Why This Matters

Dynamic IPs can be a major headache for VPN setups, causing connection drops and requiring manual configuration updates. This solution makes you **IMMUNE** to dynamic IP changes by:

- 🔄 Automatically detecting IP changes
- ⚡ Instantly propagating updates to all clients
- 🔒 Maintaining secure VPN connections
- 🤖 Requiring zero manual intervention

## 🔧 Installation

### Server Setup (IP Leaker)

1. Clone the repository
```bash
git clone https://github.com/yourusername/server-ip-syncer
cd server-ip-syncer/ip_leaker
```

2. Configure Firebase credentials
- Place your `firestore.json` credentials file in the root directory

3. Install the service
```bash
make install
```

### Client Setup (IP Viber)

1. Clone the repository
```bash
git clone https://github.com/yourusername/server-ip-syncer
cd server-ip-syncer/ip_viber
```

2. Configure the client
- Update `config.json` with your Firebase project details
- Place your `firestore.json` credentials file in the directory

3. Install the service
```bash
./install.sh
```

## 🔍 How It Works

1. The IP Leaker service on your server:
   - Monitors the server's public IP address
   - Detects any changes in real-time
   - Updates the IP in Firebase Firestore

2. The IP Viber service on your clients:
   - Polls Firestore every 5 minutes for IP updates
   - When a change is detected, it automatically:
     - Updates the WireGuard configuration
     - Restarts the WireGuard interface
     - Ensures continuous VPN connectivity

## 📝 Configuration

### IP Leaker
Configuration is handled through environment variables or a `.env` file:
- `FIREBASE_PROJECT_ID`: Your Firebase project ID
- `CHECK_INTERVAL`: IP check interval (default: 5 minutes)

### IP Viber
Edit `config.json` with your settings:
```json
{
    "firebase_config_path": "firestore.json",
    "wireguard_config_path": "/etc/wireguard/wg0.conf",
    "firebase_document_path": "ip_addresses/current_ip"
}
```

Note: The IP Viber checks for updates every 5 minutes by default (hardcoded).

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐ Support

If you find this project useful, please give it a star on GitHub! 