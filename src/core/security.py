import os
import hashlib
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

class SecurityManager:
    def __init__(self, password: str = None):
        self.key = None
        if password:
            self.set_password(password)

    def set_password(self, password: str):
        """Derives a 32-byte key from the password using SHA-256."""
        if not password:
            self.key = None
            return
        
        # Use SHA-256 to get a fixed 32-byte key
        digest = hashlib.sha256(password.encode()).digest()
        self.key = digest

    def encrypt(self, data: bytes) -> bytes:
        """Encrypts data using AES-256-GCM."""
        if not self.key:
            return data # Return plain if no encryption set (or handle as error)
        
        aesgcm = AESGCM(self.key)
        nonce = os.urandom(12)
        ciphertext = aesgcm.encrypt(nonce, data, None)
        return nonce + ciphertext

    def decrypt(self, data: bytes) -> bytes:
        """Decrypts data using AES-256-GCM."""
        if not self.key:
            return data
        
        if len(data) < 12:
            raise ValueError("Data too short")
        
        nonce = data[:12]
        ciphertext = data[12:]
        aesgcm = AESGCM(self.key)
        
        return aesgcm.decrypt(nonce, ciphertext, None)
