#!/usr/bin/env python3

import os
import subprocess
import sys
import shutil
from pathlib import Path

def get_pi_version():
    """Prompt user to select Raspberry Pi version"""
    while True:
        print("\nSelect Raspberry Pi version:")
        print("1) Raspberry Pi 3")
        print("2) Raspberry Pi 4")
        choice = input("Enter choice (1 or 2): ").strip()
        
        if choice in ['1', '2']:
            return choice
        print("Invalid choice. Please enter 1 or 2.")

def check_requirements():
    """Check if required tools are installed"""
    required_tools = ['go']
    missing_tools = []
    
    for tool in required_tools:
        if not shutil.which(tool):
            missing_tools.append(tool)
    
    if missing_tools:
        print("Error: Missing required tools:")
        for tool in missing_tools:
            print(f"  - {tool}")
        print("\nPlease install the missing tools:")
        print("  sudo apt-get install gcc-arm-linux-gnueabihf")
        sys.exit(1)

def setup_cross_compile(pi_version):
    """Set up Go cross-compilation environment"""
    # Set Go environment variables for cross-compilation
    os.environ['GOOS'] = 'linux'
    
    if pi_version == '1':  # Raspberry Pi 3
        print("Configuring for Raspberry Pi 3 (arm64)...")
        os.environ['GOARCH'] = 'arm64'
    else:  # Raspberry Pi 4
        print("Configuring for Raspberry Pi 4 (armv7)...")
        os.environ['GOARCH'] = 'arm'
        os.environ['GOARM'] = '7'
    
    os.environ['CGO_ENABLED'] = '0'  # Disable CGO
    
    # Create build directory if it doesn't exist
    build_dir = Path('build')
    build_dir.mkdir(exist_ok=True)

def build_binary():
    """Build the binary for Raspberry Pi"""
    print("Building binary for Raspberry Pi...")
    
    # Build the binary with static linking
    cmd = [
        'go', 'build',
        '-o', 'build/ip_viber',
        '-ldflags', '-s -w',  # Strip debug information
        'main.go'
    ]
    
    try:
        subprocess.run(cmd, check=True)
        print("Build successful!")
        print("\nBinary location: build/ip_viber")
        print("\nTo deploy to Raspberry Pi:")
        print("  scp build/ip_viber erwin@10.8.0.3:/home/erwin/")
    except subprocess.CalledProcessError as e:
        print(f"Build failed with error: {e}")
        sys.exit(1)

def main():
    print("Starting cross-compilation for Raspberry Pi...")
    
    # Get Pi version from user
    pi_version = get_pi_version()
    
    # Check requirements
    check_requirements()
    
    # Set up cross-compilation environment
    setup_cross_compile(pi_version)
    
    # Build the binary
    build_binary()

if __name__ == '__main__':
    main()