import sys
import argparse
from src.cli.main import run_cli

def main():
    parser = argparse.ArgumentParser(description="AnonBOX: P2P Amnesic Chat")
    parser.add_argument('mode', choices=['cli', 'gui'], nargs='?', default='gui', help='Mode to run: cli or gui (default)')
    parser.add_argument('--password', '-p', type=str, help='Vault password for encryption (Optional)')
    parser.add_argument('--name', '-n', type=str, help='Display Name (Optional)')
    
    args = parser.parse_args()
    
    if args.mode == 'cli':
        run_cli(args.password, args.name)
    else:
        # GUI Import inside function to avoid dependency issues if just running CLI
        try:
            from src.gui.app import run_gui
            run_gui(args.password, args.name)
        except ImportError as e:
            print(f"Failed to load GUI: {e}")
            print("Ensure customtkinter is installed or run in CLI mode.")
            sys.exit(1)

if __name__ == "__main__":
    main()
