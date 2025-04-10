# GoInspectorGadget Installation Guide

This guide provides detailed instructions for installing and configuring the GoInspectorGadget system.

## System Requirements

### Hardware Requirements

- **Processor**: 2.0 GHz dual-core processor or better
- **Memory**: 4 GB RAM minimum (8 GB recommended for larger investigations)
- **Disk Space**: 1 GB available space for the application (additional space required for case data)
- **Graphics**: Basic integrated graphics (dedicated GPU recommended for faster audio processing)

### Software Requirements

- **Operating System**:
  - Linux: Ubuntu 20.04 LTS or newer, CentOS 8+
  - macOS: 11.0 (Big Sur) or newer
  - Windows: Windows 10 64-bit or Windows 11
- **Go**: Version 1.23.0 or higher
- **FFmpeg**: Latest stable version
- **Git**: For source code management and installation

## Prerequisites Installation

### Installing Go

#### Linux (Ubuntu/Debian)

```bash
# Download the latest Go package
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz

# Extract the archive
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add Go to your PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verify installation
go version
```

#### macOS

```bash
# Using Homebrew (recommended)
brew install go

# Verify installation
go version
```

#### Windows

1. Download the installer from [golang.org/dl/](https://golang.org/dl/)
2. Run the installer and follow the prompts
3. Open Command Prompt and verify installation:
   ```
   go version
   ```

### Installing FFmpeg

#### Linux (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install ffmpeg

# Verify installation
ffmpeg -version
```

#### macOS

```bash
brew install ffmpeg

# Verify installation
ffmpeg -version
```

#### Windows

1. Download the latest build from [ffmpeg.org/download.html](https://ffmpeg.org/download.html)
2. Extract the archive to a folder (e.g., `C:\ffmpeg`)
3. Add FFmpeg to your PATH:
   - Right-click "This PC" and select "Properties"
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Edit the "Path" variable to include the `bin` folder (e.g., `C:\ffmpeg\bin`)
4. Open Command Prompt and verify installation:
   ```
   ffmpeg -version
   ```

## GoInspectorGadget Installation

### Option 1: Installing from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/jth/claude/GoInspectorGadget.git
   cd GoInspectorGadget
   ```

2. Build the binaries:
   ```bash
   # Create bin directory
   mkdir -p bin
   
   # Build investigator tool
   go build -o bin/investigator cmd/investigator/main.go
   
   # Build document processor
   go build -o bin/docprocessor cmd/docprocessor/*.go
   ```

3. Add the binaries to your PATH:
   
   #### Linux/macOS
   ```bash
   # Temporary (for current session)
   export PATH=$PATH:$(pwd)/bin
   
   # Permanent (add to your shell profile)
   echo 'export PATH=$PATH:'"$(pwd)/bin" >> ~/.bashrc  # or ~/.zshrc
   source ~/.bashrc  # or ~/.zshrc
   ```
   
   #### Windows
   ```cmd
   # Temporary (for current session)
   set PATH=%PATH%;%cd%\bin
   
   # Permanent (System Properties > Environment Variables)
   # Add the full path to the bin directory
   ```

4. Verify installation:
   ```bash
   investigator help
   docprocessor --help
   ```

### Option 2: Installing Pre-built Binaries (Coming Soon)

Pre-built binaries for major platforms will be available in future releases.

## Configuration

### Setting Up Working Directory

By default, GoInspectorGadget stores its data in:
- Linux/macOS: `$HOME/investigator-simulator`
- Windows: `%USERPROFILE%\investigator-simulator`

You can customize this location by setting the `INVESTIGATOR_HOME` environment variable:

#### Linux/macOS
```bash
export INVESTIGATOR_HOME=/path/to/data/directory
```

#### Windows
```cmd
set INVESTIGATOR_HOME=C:\path\to\data\directory
```

### Optional Components

#### Language Models for Speech Recognition

For optimal speech recognition performance, download these language models:

1. Venezuelan Spanish model:
   ```bash
   mkdir -p $HOME/investigator-simulator/models
   wget https://example.com/models/venezuelan-spanish.bin -O $HOME/investigator-simulator/models/venezuelan-spanish.bin
   ```

2. American English with police terminology:
   ```bash
   wget https://example.com/models/american-police.bin -O $HOME/investigator-simulator/models/american-police.bin
   ```

3. Configure model paths:
   ```bash
   # Create config directory if it doesn't exist
   mkdir -p $HOME/.config/inspectorgadget
   
   # Create config file
   cat > $HOME/.config/inspectorgadget/config.json << EOF
   {
     "models": {
       "venezuelan_spanish": "$HOME/investigator-simulator/models/venezuelan-spanish.bin",
       "american_english": "$HOME/investigator-simulator/models/american-police.bin"
     }
   }
   EOF
   ```

## Verifying Installation

Run the following commands to verify that the system is properly installed:

```bash
# Check investigator tool
investigator help

# Check document processor
docprocessor --type audio --help
```

If both commands display their help messages, the installation is successful.

## Troubleshooting

### Common Installation Issues

#### Go Module Issues

**Issue**: `go: cannot find module providing package github.com/...`

**Solution**:
```bash
go mod tidy
```

#### FFmpeg Not Found

**Issue**: `Error: ffmpeg command not found`

**Solution**: Ensure FFmpeg is installed and in your PATH. You can check with:
```bash
ffmpeg -version
```

#### Permission Issues

**Issue**: `Permission denied` when running executables

**Solution**: Make sure the binaries are executable:
```bash
chmod +x bin/investigator bin/docprocessor
```

#### Path Issues

**Issue**: `Command not found` when trying to run the tools

**Solution**: Ensure the bin directory is in your PATH:
```bash
echo $PATH  # Check if the bin directory is listed
```

## Next Steps

After installation, refer to the [User Manual](USER_MANUAL.md) for detailed usage instructions, or the [Quick Reference Guide](QUICK_REFERENCE.md) for common commands. 