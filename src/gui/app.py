import customtkinter as ctk
import threading
import tkinter as tk
from tkinter import filedialog, messagebox
import base64
import os
from ..core.network import NetworkManager
from ..core.security import SecurityManager

ctk.set_appearance_mode("System")
ctk.set_default_color_theme("blue")

class LoginDialog(ctk.CTkToplevel):
    def __init__(self, parent):
        super().__init__(parent)
        self.title("Enter the Cove")
        self.geometry("400x300")
        self.resizable(False, False)
        self.username = None
        self.password = None
        
        self.label = ctk.CTkLabel(self, text="Welcome to AnonBOX", font=ctk.CTkFont(size=20, weight="bold"))
        self.label.pack(pady=20)
        
        self.name_entry = ctk.CTkEntry(self, placeholder_text="Display Name (Optional)")
        self.name_entry.pack(pady=10, padx=20, fill="x")
        
        self.pass_entry = ctk.CTkEntry(self, placeholder_text="Vault Password (Optional)", show="*")
        self.pass_entry.pack(pady=10, padx=20, fill="x")
        
        self.btn = ctk.CTkButton(self, text="Enter", command=self.on_enter)
        self.btn.pack(pady=20)
        
        self.protocol("WM_DELETE_WINDOW", self.on_close)
        self.transient(parent)
        self.grab_set()
        self.parent = parent

    def on_enter(self):
        self.username = self.name_entry.get() or None
        self.password = self.pass_entry.get() or None
        self.destroy()

    def on_close(self):
        self.parent.destroy()

class App(ctk.CTk):
    def __init__(self, cli_password=None, cli_username=None):
        super().__init__()
        self.withdraw() # Hide until login

        self.title("AnonBOX 🏴‍☠️")
        self.geometry("800x600")
        
        # Login Logic
        self.password = cli_password
        self.username = cli_username
        
        if not self.password and not self.username:
             dialog = LoginDialog(self)
             self.wait_window(dialog)
             self.username = dialog.username
             self.password = dialog.password
        
        self.deiconify()

        self.security = SecurityManager(self.password)
        self.nm = NetworkManager(self.security, self.username)
        
        # Grid layout
        self.grid_columnconfigure(1, weight=1)
        self.grid_rowconfigure(0, weight=1)

        # Sidebar
        self.sidebar_frame = ctk.CTkFrame(self, width=140, corner_radius=0)
        self.sidebar_frame.grid(row=0, column=0, sticky="nsew")
        self.logo_label = ctk.CTkLabel(self.sidebar_frame, text="AnonBOX", font=ctk.CTkFont(size=20, weight="bold"))
        self.logo_label.grid(row=0, column=0, padx=20, pady=(20, 10))
        
        self.user_label = ctk.CTkLabel(self.sidebar_frame, text=f"User: {self.nm.username[:10]}", font=ctk.CTkFont(size=12))
        self.user_label.grid(row=1, column=0, padx=20, pady=(0, 10))

        self.peers_label = ctk.CTkLabel(self.sidebar_frame, text="Peers:", anchor="w")
        self.peers_label.grid(row=2, column=0, padx=20, pady=(10, 0))
        
        self.peer_listbox = tk.Listbox(self.sidebar_frame, height=20, bg="#2b2b2b", fg="white", borderwidth=0)
        self.peer_listbox.grid(row=3, column=0, padx=20, pady=10, sticky="nsew")
        self.peer_listbox.bind('<<ListboxSelect>>', self.on_peer_select)
        
        # Main Chat Area
        self.chat_frame = ctk.CTkFrame(self, corner_radius=0)
        self.chat_frame.grid(row=0, column=1, sticky="nsew")
        self.chat_frame.grid_rowconfigure(0, weight=1)
        self.chat_frame.grid_columnconfigure(0, weight=1)

        self.chat_display = ctk.CTkTextbox(self.chat_frame, state="disabled")
        self.chat_display.grid(row=0, column=0, padx=20, pady=20, sticky="nsew")

        self.input_frame = ctk.CTkFrame(self.chat_frame, height=50)
        self.input_frame.grid(row=1, column=0, padx=20, pady=20, sticky="ew")
        
        self.msg_entry = ctk.CTkEntry(self.input_frame, placeholder_text="Type message...")
        self.msg_entry.pack(side="left", fill="both", expand=True, padx=(0, 10))
        self.msg_entry.bind("<Return>", self.send_message)
        
        self.send_btn = ctk.CTkButton(self.input_frame, text="Send", command=self.send_message)
        self.send_btn.pack(side="right")
        
        self.file_btn = ctk.CTkButton(self.input_frame, text="📎", width=40, command=self.send_file)
        self.file_btn.pack(side="right", padx=(0, 10))

        self.selected_peer = None
        self.after(1000, self.update_peers)
        
        # Start Network
        self.nm.start(self.on_network_message)

    def on_network_message(self, msg):
        sender_id = msg.get('sender_id', 'unknown')[:8]
        sender_name = msg.get('sender_name', sender_id)
        content = msg.get('content', '')
        msg_type = msg.get('type', 'chat')
        
        display_text = ""
        if msg_type == 'chat':
            display_text = f"[{sender_name}]: {content}\n"
        elif msg_type == 'file':
             filename = msg.get('filename', 'unknown_file')
             display_text = f"[{sender_name}] sent file: {filename}\n"
             # Auto-save or prompt? For simplicity, auto-save to 'downloads' or current dir
             try:
                 file_data = base64.b64decode(content)
                 with open(f"received_{filename}", "wb") as f:
                     f.write(file_data)
                 display_text += f"(Saved as received_{filename})\n"
             except Exception as e:
                 display_text += f"(Error saving file: {e})\n"
             
        self.chat_display.configure(state="normal")
        self.chat_display.insert("end", display_text)
        self.chat_display.configure(state="disabled")
        self.chat_display.see("end")

    def update_peers(self):
        # Refresh peer list
        current_selection = self.peer_listbox.curselection()
        self.peer_listbox.delete(0, "end")
        
        for name, data in self.nm.peers.items():
            display_name = f"{data['username']} ({name.split('.')[0]})"
            self.peer_listbox.insert("end", display_name)
            
        if current_selection:
             try:
                self.peer_listbox.select_set(current_selection[0])
             except:
                 pass
                 
        self.after(2000, self.update_peers)

    def on_peer_select(self, event):
        selection = self.peer_listbox.curselection()
        if selection:
            display_string = self.peer_listbox.get(selection[0]) # "User (AnonPeer...)"
            # Find key in peers
            for original_name, data in self.nm.peers.items():
                if f"{data['username']} ({original_name.split('.')[0]})" == display_string:
                    self.selected_peer = data
                    break

    def send_message(self, event=None):
        text = self.msg_entry.get()
        if not text:
            return
            
        if not self.selected_peer:
            messagebox.showwarning("Warning", "Select a peer first!")
            return

        if self.nm.send_message(self.selected_peer['address'], self.selected_peer['port'], content=text):
             self.chat_display.configure(state="normal")
             self.chat_display.insert("end", f"[Me]: {text}\n")
             self.chat_display.configure(state="disabled")
             self.msg_entry.delete(0, "end")
        else:
             messagebox.showerror("Error", "Failed to send message.")

    def send_file(self):
        if not self.selected_peer:
            messagebox.showwarning("Warning", "Select a peer first!")
            return
            
        filename = filedialog.askopenfilename()
        if filename:
            try:
                with open(filename, "rb") as f:
                    file_data = f.read()
                    
                encoded_data = base64.b64encode(file_data).decode('utf-8')
                short_name = os.path.basename(filename)
                
                if self.nm.send_message(self.selected_peer['address'], self.selected_peer['port'], message_type='file', content=encoded_data, filename=short_name):
                     self.chat_display.configure(state="normal")
                     self.chat_display.insert("end", f"[Me] sent file: {short_name}\n")
                     self.chat_display.configure(state="disabled")
                else:
                    messagebox.showerror("Error", "Failed to send file.")
            except Exception as e:
                messagebox.showerror("Error", f"File processing failed: {e}")

    def on_closing(self):
        self.nm.stop()
        self.destroy()

def run_gui(password=None, username=None):
    app = App(password, username)
    app.protocol("WM_DELETE_WINDOW", app.on_closing)
    app.mainloop()
