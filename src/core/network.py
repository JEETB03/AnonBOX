import socket
import threading
import json
import time
import uuid
import struct
import logging
from zeroconf import Zeroconf, ServiceInfo, ServiceBrowser, ServiceListener
from .security import SecurityManager

# Configuration
PORT = 0 # Random port
SERVICE_TYPE = "_anonbox._tcp.local."
BUF_SIZE = 4096

class PeerListener(ServiceListener):
    def __init__(self, network_manager):
        self.nm = network_manager

    def remove_service(self, zc, type, name):
        if name in self.nm.peers:
           del self.nm.peers[name]

    def add_service(self, zc, type, name):
        info = zc.get_service_info(type, name)
        if info:
            self.nm.add_peer(name, info)

    def update_service(self, zc, type, name):
        pass

class NetworkManager:
    def __init__(self, security_manager: SecurityManager, username: str = None):
        self.security = security_manager
        self.peers = {} # name -> {address, port, id, username}
        self.my_id = str(uuid.uuid4())
        self.username = username if username else f"Anon-{self.my_id[:6]}"
        self.running = False
        self.server_socket = None
        self.port = 0
        self.zeroconf = Zeroconf()
        self.msg_callback = None

        # Setup server
        self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.server_socket.bind(('0.0.0.0', 0))
        self.port = self.server_socket.getsockname()[1]
        self.server_socket.listen(5)
        
        logging.basicConfig(level=logging.INFO)
        self.logger = logging.getLogger("AnonBOX")


    def start(self, callback):
        self.msg_callback = callback
        self.running = True
        
        # Start TCP Server
        threading.Thread(target=self._accept_loop, daemon=True).start()

        # Register mDNS service
        props = {'id': self.my_id, 'user': self.username}
        info = ServiceInfo(
            SERVICE_TYPE,
            f"AnonPeer-{self.my_id[:8]}.{SERVICE_TYPE}",
            addresses=[socket.inet_aton(self._get_local_ip())],
            port=self.port,
            properties=props,
            server=f"anonbox-{self.my_id[:8]}.local."
        )
        self.zeroconf.register_service(info)

        # Browse for peers
        self.browser = ServiceBrowser(self.zeroconf, SERVICE_TYPE, PeerListener(self))
        self.logger.info(f"Started on port {self.port}. ID: {self.my_id}, User: {self.username}")

    def add_peer(self, name, info):
        properties = info.properties
        peer_id = ""
        peer_user = "Unknown"
        
        if properties:
            if b'id' in properties:
                peer_id = properties[b'id'].decode('utf-8')
            if b'user' in properties:
                peer_user = properties[b'user'].decode('utf-8')

        if peer_id == self.my_id:
             return

        address = socket.inet_ntoa(info.addresses[0])
        port = info.port
        self.peers[name] = {'address': address, 'port': port, 'id': peer_id, 'username': peer_user}
        self.logger.info(f"Found peer: {peer_user} ({name}) at {address}:{port}")

    def _accept_loop(self):
        while self.running:
            try:
                client, addr = self.server_socket.accept()
                threading.Thread(target=self._handle_client, args=(client,), daemon=True).start()
            except Exception as e:
                if self.running:
                    self.logger.error(f"Accept error: {e}")

    def _recv_all(self, sock, count):
        buf = b''
        while count:
            newbuf = sock.recv(count)
            if not newbuf: return None
            buf += newbuf
            count -= len(newbuf)
        return buf

    def _handle_client(self, client_sock):
        try:
            # Read 4-byte length header
            length_bytes = self._recv_all(client_sock, 4)
            if not length_bytes:
                return
            
            msg_len = struct.unpack('>I', length_bytes)[0]
            
            # Read full message
            encrypted_data = self._recv_all(client_sock, msg_len)
            if not encrypted_data:
                return

            # Decrypt
            try:
                decrypted = self.security.decrypt(encrypted_data)
                # Parse JSON
                msg = json.loads(decrypted.decode('utf-8'))
                
                if self.msg_callback:
                    self.msg_callback(msg)
            except Exception as e:
                 self.logger.error(f"Decryption/Parse error: {e}")

        except Exception as e:
            self.logger.error(f"Client error: {e}")
        finally:
            client_sock.close()

    def send_message(self, target_ip, target_port, message_type="chat", content=None, filename=None):
        msg_payload = {
            'sender_id': self.my_id,
            'sender_name': self.username,
            'type': message_type, 
            'content': content,
            'filename': filename,
            'timestamp': time.time()
        }
        
        payload_bytes = json.dumps(msg_payload).encode('utf-8')
        encrypted = self.security.encrypt(payload_bytes)
        
        # Add length header
        length_header = struct.pack('>I', len(encrypted))

        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((target_ip, target_port))
            s.sendall(length_header + encrypted)
            s.close()
            return True
        except Exception as e:
            self.logger.error(f"Send error: {e}")
            return False

    def broadcast(self, message):
         for peer in self.peers.values():
             self.send_message(peer['address'], peer['port'], content=message)

    def _get_local_ip(self):
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            s.connect(('8.8.8.8', 80))
            IP = s.getsockname()[0]
            s.close()
            return IP
        except:
            return "127.0.0.1"

    def stop(self):
        self.running = False
        self.zeroconf.close()
        self.server_socket.close()

