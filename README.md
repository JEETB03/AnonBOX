# AnonBOX 🏴‍☠️

**AnonBOX** is a secure, decentralized, P2P, amnesic chatroom and file sharing application. It is designed for users who prioritize privacy, security, and anonymity.

## 🚀 Features

- **P2P Architecture**: No central server. All communication is direct between peers using `libp2p`.
- **Amnesic Security**: No logs, no database. All chat history is stored in RAM and wiped instantly upon exit.
- **Military-Grade Encryption**: End-to-end AES-256-GCM encryption for all messages and file transfers.
- **Secure File Sharing**: Transfer files (up to 100MB+) directly between peers with encryption.
- **Cross-Platform**: Runs on Windows, Linux, macOS, and Android.
- **Dual Interface**:
  - **CLI**: For power users and headless environments.
  - **GUI**: A beautiful, simple interface built with Fyne.

## 📦 Installation

### Prerequisites
- **Go** (v1.16+)
- **GCC** (Required for GUI build only)

### Build from Source

1. **Clone the Repository**
   ```bash
   git clone https://github.com/JEETB03/AnonBOX.git
   cd AnonBOX
   ```

2. **Build**
   - **Windows**: Run `build.bat`
   - **Linux/macOS**: Run `./build.sh`

   *Note: If GCC is not installed, the script will automatically build only the CLI version.*

## 🛠 Usage

### CLI (`anonbox-cli`)

Start the node:
```bash
./anonbox-cli start -p "YourSecretPassword"
```

**Commands**:
- `peers`: List connected peers.
- `chat <peerID> <message>`: Send a secure message.
- `share <peerID> <filePath>`: Send a file securely.
- `broadcast <message>`: Broadcast a message to all connected peers.
- `exit`: Close the node and wipe all data.

### GUI (`anonbox-gui`)

1. Launch the application.
2. Enter a "Vault Password" (optional) to enable encryption.
3. Use the tabs to navigate:
   - **Chat**: Real-time secure messaging.
   - **Peers**: View and manage connected peers.
   - **Files**: Send files to peers.

## 📱 Android Support

AnonBOX is compatible with Android. To build the APK:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os android -appID com.anonbox.app -icon icon.png
```

## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
