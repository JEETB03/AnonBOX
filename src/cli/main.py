import cmd
import time
import threading
import base64
import os
from ..core.network import NetworkManager
from ..core.security import SecurityManager

class AnonCLI(cmd.Cmd):
    intro = 'Welcome to AnonBOX CLI. Type help or ? to list commands.\n'
    prompt = '(anonbox) '
    
    def __init__(self, password=None, username=None):
        super().__init__()
        self.security = SecurityManager(password)
        self.nm = NetworkManager(self.security, username)
        self.nm.start(self.on_message)
        self.prompt = f"({self.nm.username}) "

    def on_message(self, msg):
        sender_id = msg.get('sender_id', 'unknown')[:8]
        sender_name = msg.get('sender_name', sender_id)
        content = msg.get('content', '')
        msg_type = msg.get('type', 'chat')
        
        if msg_type == 'chat':
            print(f"\n[{sender_name}]: {content}\n{self.prompt}", end='', flush=True)
        elif msg_type == 'file':
             filename = msg.get('filename', 'unknown_file')
             print(f"\n[{sender_name}] sent file: {filename}\n{self.prompt}", end='', flush=True)
             try:
                 file_data = base64.b64decode(content)
                 with open(f"received_{filename}", "wb") as f:
                     f.write(file_data)
                 print(f"(Saved as received_{filename})\n{self.prompt}", end='', flush=True)
             except Exception as e:
                 print(f"(Error saving file: {e})\n{self.prompt}", end='', flush=True)

    def do_peers(self, arg):
        'List connected peers'
        print("\nConnected Peers:")
        for name, info in self.nm.peers.items():
            print(f"- {info['username']} ({name}) at {info['address']}:{info['port']}")
        print("")

    def do_chat(self, arg):
        'Send a chat message: chat <peer_name_or_partial_id> <message>'
        parts = arg.split(' ', 1)
        if len(parts) < 2:
            print("Usage: chat <peer_id> <message>")
            return
        
        target_id_part = parts[0]
        message = parts[1]
        
        target = self._find_peer(target_id_part)
        if target:
            if self.nm.send_message(target['address'], target['port'], content=message):
                print(f"Sent to {target['username']}")
            else:
                print("Failed to send.")
        else:
            print("Peer not found.")

    def do_share(self, arg):
        'Send a file: share <peer_id> <filename>'
        parts = arg.split(' ', 1)
        if len(parts) < 2:
            print("Usage: share <peer_id> <filename>")
            return
            
        target_id_part = parts[0]
        filename = parts[1]
        
        if not os.path.exists(filename):
            print("File not found.")
            return

        target = self._find_peer(target_id_part)
        if target:
            try:
                with open(filename, "rb") as f:
                    file_data = f.read()
                
                print(f"Reading file ({len(file_data)} bytes)...")
                encoded_data = base64.b64encode(file_data).decode('utf-8')
                short_name = os.path.basename(filename)
                
                if self.nm.send_message(target['address'], target['port'], message_type='file', content=encoded_data, filename=short_name):
                    print(f"Sent file '{short_name}' to {target['username']}")
                else:
                    print("Failed to send file.")
            except Exception as e:
                print(f"Error processing file: {e}")
        else:
            print("Peer not found.")

    def do_broadcast(self, arg):
        'Broadcast a message to all peers: broadcast <message>'
        if not arg:
             print("Usage: broadcast <message>")
             return
        self.nm.broadcast(arg)
        print("Broadcast sent.")

    def do_exit(self, arg):
        'Exit the application'
        print("Exiting...")
        self.nm.stop()
        return True

    def _find_peer(self, partial_id):
        # Exact match
        if partial_id in self.nm.peers:
            return self.nm.peers[partial_id]
        
        # Partial match
        for name, info in self.nm.peers.items():
            if partial_id in name or partial_id in info['username']:
                return info
        return None

def run_cli(password=None, username=None):
    try:
        if password:
            print("🔒 Encryption Enabled.")
        else:
            print("⚠️  No password provided. Running in plain text mode.")
            
        AnonCLI(password, username).cmdloop()
    except KeyboardInterrupt:
        print("\nExiting...")
