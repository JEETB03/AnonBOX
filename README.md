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
   git clone https://github.com/JEETB03/AnonBox.git
   cd AnonBox
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
- `connect <multiaddr>`: Connect to a peer manually (WAN/Internet).
  - Example: `connect /ip4/1.2.3.4/tcp/4001/p2p/QmPeerID...`
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

## 🌐 How to Connect

### Local Network (LAN)
AnonBOX uses **mDNS** for automatic discovery. If you and your friend are on the same Wi-Fi or LAN, simply start the application on both devices. You will automatically discover each other and appear in the `peers` list.

### Internet (WAN)
To connect with someone over the internet:
1. **Port Forwarding**: Ensure the host allows incoming connections on the P2P port (random by default, or configured). *Note: Currently AnonBOX binds to a random port. For WAN, you may need to check your logs for the listening address.*
2. **Share Address**: One peer must share their **Multiaddress** with the other.
   - The address looks like: `/ip4/YOUR_PUBLIC_IP/tcp/PORT/p2p/YOUR_PEER_ID`
3. **Connect**: The other peer runs the `connect` command with this address.
   ```bash
   connect /ip4/203.0.113.1/tcp/12345/p2p/Qm...
   ```

## 🔐 Working Mechanism

1. **Initialization**: When AnonBOX starts, it generates a unique Libp2p host and identity. It starts listening on a random TCP port.
2. **Discovery**:
   - **mDNS**: Broadcasts presence on the local network to find peers.
   - **DHT/Manual**: Can connect to specific peers via multiaddress.
3. **Secure Channel**:
   - **Transport Layer**: Uses Noise or TLS to secure the connection stream.
   - **Application Layer**: If a password is provided, all messages and files are **AES-256-GCM** encrypted before sending. Only peers with the same password can decrypt the content.
4. **Amnesic Storage**:
   - Messages are stored in a Go slice (RAM).
   - When the application closes, the OS reclaims the memory, effectively wiping all traces. No data is ever written to disk (except downloaded files).

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
