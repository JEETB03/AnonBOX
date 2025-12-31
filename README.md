# AnonBOX 🏴‍☠️

**AnonBOX** is a secure, decentralized, P2P, amnesic chatroom and file sharing application. It is designed for users who prioritize privacy, security, and anonymity.

## 🚀 Features

-   **True P2P Architecture**: No central server. Direct communication between peers using local discovery (mDNS) and TCP sockets.
-   **Amnesic Security**: No logs, no database. All chat history is stored in RAM and wiped instantly upon exit.
-   **Military-Grade Encryption**: Optional end-to-end **AES-256-GCM** encryption for all messages and file transfers (`Vault Password`).
-   **Secure File Sharing**: Transfer files of any size (100MB+) directly between peers.
-   **Cross-Platform**: Runs on Linux and Windows (macOS experimental).
-   **Dual Interface**:
    -   **CLI**: For power users and headless environments.
    -   **GUI**: A beautiful, modern interface built with `CustomTkinter`.

## 📦 Installation & Usage

### Method 1: The Launcher (Recommended)
We provide automated launchers that handle environment setup and dependencies.

**Linux / macOS**:
```bash
./run.sh
```

**Windows**:
```batch
run.bat
```

Select **Option 1** for GUI or **Option 2** for CLI from the menu.

### Method 2: Manual Run
Prerequisites: Python 3.10+

1.  **Install Dependencies**:
    ```bash
    pip install -r requirements.txt
    ```

2.  **Run GUI**:
    ```bash
    python3 main.py gui
    ```
    *Optional: You will be prompted for a Display Name and Vault Password.*

3.  **Run CLI**:
    ```bash
    python3 main.py cli --name "MyUser" --password "MySecret"
    ```

## 🔐 Core Philosophy & Mechanism

1.  **Initialization**: AnonBOX generates a random ephemeral ID and Identity on startup.
2.  **Discovery**: Uses ZeroConf (mDNS) to find other AnonBOX peers on the local network (`.local` domain).
3.  **Secure Channel**:
    -   If a **Vault Password** is provided, a SHA-256 key is derived.
    -   All traffic is encrypted with **AES-256-GCM`** before leaving the device.
    -   Peers without the password cannot decrypt messages.
4.  **Amnesic Storage**:
    -   Messages reside only in volatile memory (RAM).
    -   Closing the app destroys the key and all data.

## 🤝 Contributing
Open source and privacy-focused. Contributions are welcome!

## 📄 License

This project is licensed under the MIT License.
