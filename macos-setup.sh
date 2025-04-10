#!/bin/bash

# Go Document and Media Processing Environment Setup Script for macOS
# This script sets up a comprehensive environment for processing, indexing, and analyzing
# documents, photos, audio, and videos in Go on macOS

echo "Starting Go Media Processing Environment Setup for macOS"
echo "========================================================"

# Create a log file
LOG_FILE="setup_log.txt"
touch $LOG_FILE
exec > >(tee -a "$LOG_FILE") 2>&1

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to check if a Go package is installed
go_package_installed() {
  go list "$1" >/dev/null 2>&1
}

# Install Homebrew if not already installed
echo "[1/10] Checking Homebrew installation..."
if ! command_exists brew; then
  echo "Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
  echo "Homebrew is already installed."
  brew update
fi

echo "[2/10] Installing Go language..."
if command_exists go; then
  echo "Go is already installed. Version: $(go version)"
else
  echo "Installing Go..."
  brew install go
  
  # Set up Go environment variables
  echo 'export GOPATH=$HOME/go' >> ~/.zshrc
  echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
  
  # Apply changes to current session
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin
  
  echo "Go installation completed. Version: $(go version)"
fi

# Set up Go workspace
echo "[3/10] Setting up Go workspace..."
mkdir -p $HOME/go/{bin,pkg,src}

# Create project directory
PROJECT_DIR="$HOME/media-processor"
mkdir -p $PROJECT_DIR
cd $PROJECT_DIR

# Initialize Go module
echo "[4/10] Initializing Go module..."
go mod init github.com/yourusername/media-processor

# Install system dependencies
echo "[5/10] Installing system dependencies..."

# Install essential dependencies using Homebrew
brew install ffmpeg opencv poppler cairo pkg-config cmake portaudio python3

# Install Go packages
echo "[6/10] Installing document and media processing libraries..."

# Document processing
go get -u github.com/jung-kurt/gofpdf
go get -u github.com/pdfcpu/pdfcpu/pkg/api
go get -u go.mozilla.org/pkcs7

# Image processing
go get -u github.com/disintegration/imaging
go get -u gocv.io/x/gocv

# Audio and video processing
go get -u github.com/u2takey/ffmpeg-go
go get -u github.com/go-audio/audio
go get -u github.com/go-audio/wav

# Install Whisper for speech recognition
echo "[7/10] Installing Whisper and Go bindings..."
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
make
# Download medium model for better accent handling instead of base.en
bash ./models/download-ggml-model.sh medium
cd ..

echo "Setting up Go bindings for whisper.cpp..."
mkdir -p $GOPATH/src/github.com/ggerganov/whisper-go
cp -r whisper.cpp/bindings/go/* $GOPATH/src/github.com/ggerganov/whisper-go/

echo "Setting environment variables for whisper.cpp..."
echo 'export CGO_CFLAGS="-I'$PROJECT_DIR'/whisper.cpp"' >> ~/.zshrc
echo 'export CGO_LDFLAGS="-L'$PROJECT_DIR'/whisper.cpp -lwhisper"' >> ~/.zshrc
source ~/.zshrc

# Install DeepSpeech for improved American English recognition
echo "[8/10] Installing Mozilla DeepSpeech..."
brew install sox
go get -u github.com/asticode/go-astideepspeech
mkdir -p $PROJECT_DIR/deepspeech-models
cd $PROJECT_DIR/deepspeech-models
# Download the DeepSpeech model files
curl -LO https://github.com/mozilla/DeepSpeech/releases/download/v0.9.3/deepspeech-0.9.3-models.pbmm
curl -LO https://github.com/mozilla/DeepSpeech/releases/download/v0.9.3/deepspeech-0.9.3-models.scorer
cd ..

# Install VOSK for real-time and Spanish accented speech
echo "[9/10] Installing VOSK speech recognition..."
go get -u github.com/alphacep/vosk-api/go
mkdir -p $PROJECT_DIR/vosk-models
cd $PROJECT_DIR/vosk-models
# Download English and Spanish models for VOSK
curl -LO https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip
curl -LO https://alphacephei.com/vosk/models/vosk-model-small-es-0.42.zip
unzip vosk-model-small-en-us-0.15.zip
unzip vosk-model-small-es-0.42.zip
rm vosk-model-small-en-us-0.15.zip vosk-model-small-es-0.42.zip
cd ..

# Setup diarization with pyannote.audio
echo "[10/10] Setting up pyannote.audio for speaker diarization..."
mkdir -p $PROJECT_DIR/diarization
cd $PROJECT_DIR/diarization
python3 -m venv venv
source venv/bin/activate
pip install pyannote.audio torch
# Create simple Python bridge for Go to call
cat > diarization_bridge.py << 'EOF'
#!/usr/bin/env python3
import sys
import torch
from pyannote.audio import Pipeline

def process_audio(audio_file, output_file):
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization@2.1",
                                       use_auth_token="YOUR_HF_TOKEN")
    diarization = pipeline(audio_file)
    
    with open(output_file, "w") as f:
        for turn, _, speaker in diarization.itertracks(yield_label=True):
            print(f"{turn.start:.3f} {turn.end:.3f} {speaker}", file=f)
    
    return 0

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: diarization_bridge.py <audio_file> <output_file>")
        sys.exit(1)
    
    audio_file = sys.argv[1]
    output_file = sys.argv[2]
    sys.exit(process_audio(audio_file, output_file))
EOF
chmod +x diarization_bridge.py
deactivate
cd ..

# Create a Go wrapper for diarization
mkdir -p $PROJECT_DIR/speech
cat > $PROJECT_DIR/speech/diarization.go << 'EOF'
package speech

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SpeakerSegment struct {
	Start   float64
	End     float64
	Speaker string
}

func DiarizeAudio(audioFilePath, outputPath string) ([]SpeakerSegment, error) {
	cmd := exec.Command(os.Getenv("PROJECT_DIR")+"/diarization/venv/bin/python",
		os.Getenv("PROJECT_DIR")+"/diarization/diarization_bridge.py",
		audioFilePath, outputPath)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("diarization failed: %v", err)
	}

	// Parse the output file
	segments, err := parseDiarizationOutput(outputPath)
	if err != nil {
		return nil, err
	}

	return segments, nil
}

func parseDiarizationOutput(filePath string) ([]SpeakerSegment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var segments []SpeakerSegment
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		
		if len(parts) != 3 {
			continue
		}
		
		start, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		
		end, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}
		
		segments = append(segments, SpeakerSegment{
			Start:   start,
			End:     end,
			Speaker: parts[2],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}
EOF

# Create a sample app that demonstrates multi-accent speech recognition
cat > $PROJECT_DIR/main.go << 'EOF'
package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	// Document processing
	_ "github.com/jung-kurt/gofpdf"
	_ "github.com/pdfcpu/pdfcpu/pkg/api"
	
	// Image processing
	_ "github.com/disintegration/imaging"
	_ "gocv.io/x/gocv"
	
	// FFmpeg
	ffmpeg "github.com/u2takey/ffmpeg-go"
	
	// Audio processing
	_ "github.com/go-audio/audio"
	_ "github.com/go-audio/wav"
	
	// Search and indexing
	_ "github.com/blevesearch/bleve/v2"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("   Multi-Accent Speech Processing System")
	fmt.Println("==============================================")
	
	fmt.Println("\nThis system is configured for processing:")
	fmt.Println("1. Venezuelan-accented Spanish speakers")
	fmt.Println("2. American English speakers (e.g., police officers)")
	
	fmt.Println("\nAvailable speech recognition engines:")
	fmt.Println("- OpenAI Whisper (medium model for improved accent handling)")
	fmt.Println("- Mozilla DeepSpeech (optimized for American English)")
	fmt.Println("- VOSK (with both English and Spanish models)")
	fmt.Println("- pyannote.audio (for speaker diarization/identification)")
	
	fmt.Println("\nTo process an audio file:")
	fmt.Println("1. First split speakers using the diarization tool")
	fmt.Println("2. Process Venezuelan accent with Whisper+VOSK Spanish model")
	fmt.Println("3. Process American accent with DeepSpeech+Whisper")
	fmt.Println("4. Combine results with timestamps from diarization")
	
	fmt.Println("\nModel locations:")
	fmt.Println("- Whisper: " + filepath.Join(os.Getenv("PROJECT_DIR"), "whisper.cpp", "models"))
	fmt.Println("- DeepSpeech: " + filepath.Join(os.Getenv("PROJECT_DIR"), "deepspeech-models"))
	fmt.Println("- VOSK: " + filepath.Join(os.Getenv("PROJECT_DIR"), "vosk-models"))
	
	fmt.Println("\nUsage example (once implemented):")
	fmt.Println("go run cmd/process_interview/main.go --input interview.mp3 --output transcript.txt")
}
EOF

# Update go.mod and download all dependencies
go mod tidy

echo "\nSetup complete! To test your installation, run:"
echo "  cd $PROJECT_DIR && go run main.go"
echo "\nIf the test program runs without errors, your environment is ready."
echo "Check $LOG_FILE for the complete installation log."

echo "\nIMPORTANT: For pyannote.audio to work properly, you need to:"
echo "1. Create a Hugging Face account at https://huggingface.co/"
echo "2. Accept the pyannote/speaker-diarization-2.1 model terms"
echo "3. Create an access token at https://huggingface.co/settings/tokens"
echo "4. Replace 'YOUR_HF_TOKEN' in diarization_bridge.py with your token"

# Optional: Run the test program
read -p "Would you like to run the test program now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    go run main.go
fi 