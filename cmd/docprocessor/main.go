package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jth/claude/GoInspectorGadget/pkg/speech"
)

func main() {
	// Define command-line flags for choosing the processor
	processorType := flag.String("type", "audio", "Type of processing: 'audio' or 'interview'")
	flag.Parse()

	// Choose the processor based on the type flag
	switch *processorType {
	case "interview":
		processInterviewAudio()
	case "audio":
		processAudio()
	default:
		fmt.Printf("Unknown processor type: %s\n", *processorType)
		flag.Usage()
		os.Exit(1)
	}
}

func processAudio() {
	// Define command-line flags
	audioFilePtr := flag.String("audio", "", "Path to the audio file to process")
	outputFilePtr := flag.String("output", "", "Path to save the transcription (default: same as input with .txt extension)")
	accentTypePtr := flag.String("accent", "auto", "Accent type: 'auto', 'venezuelan', 'american', or 'generic'")
	modelPathPtr := flag.String("model", "", "Path to the model (optional)")
	useGPUPtr := flag.Bool("gpu", false, "Use GPU acceleration if available")

	// Parse command-line flags
	flag.Parse()

	// Validate input file
	if *audioFilePtr == "" {
		fmt.Println("Error: Audio file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(*audioFilePtr); os.IsNotExist(err) {
		fmt.Printf("Error: Audio file not found: %s\n", *audioFilePtr)
		os.Exit(1)
	}

	// Set default output file if not specified
	outputFile := *outputFilePtr
	if outputFile == "" {
		outputFile = strings.TrimSuffix(*audioFilePtr, filepath.Ext(*audioFilePtr)) + ".txt"
	}

	// Determine accent type
	var accentType speech.AccentType
	switch strings.ToLower(*accentTypePtr) {
	case "venezuelan", "spanish":
		accentType = speech.VenezuelanSpanish
	case "american", "police":
		accentType = speech.AmericanEnglish
	case "generic":
		accentType = speech.GenericAccent
	case "auto":
		// Automatically detect accent
		detectedType, confidence, err := speech.DetectAccent(*audioFilePtr)
		if err != nil {
			fmt.Printf("Warning: Failed to detect accent: %v. Using generic accent.\n", err)
			accentType = speech.GenericAccent
		} else {
			accentType = detectedType
			fmt.Printf("Detected accent: %v with %.2f%% confidence\n",
				formatAccentType(detectedType), confidence*100)
		}
	default:
		fmt.Printf("Warning: Unknown accent type '%s'. Using automatic detection.\n", *accentTypePtr)
		// Automatically detect accent
		detectedType, _, err := speech.DetectAccent(*audioFilePtr)
		if err != nil {
			accentType = speech.GenericAccent
		} else {
			accentType = detectedType
		}
	}

	// Create recognizer options
	options := speech.RecognizerOptions{
		ModelPath:    *modelPathPtr,
		AccentType:   accentType,
		UseGPU:       *useGPUPtr,
		SampleRate:   16000,
		NumChannels:  1,
		SpecialWords: getSpecialWordsForAccent(accentType),
	}

	// Create appropriate recognizer based on accent
	recognizer, err := speech.RecognizerFactory(accentType, options)
	if err != nil {
		fmt.Printf("Error creating recognizer: %v\n", err)
		os.Exit(1)
	}
	defer recognizer.Close()

	// Initialize the recognizer
	fmt.Println("Initializing speech recognizer...")
	startTime := time.Now()
	if err := recognizer.Initialize(); err != nil {
		fmt.Printf("Error initializing recognizer: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Initialization completed in %.2f seconds\n", time.Since(startTime).Seconds())

	// Process the audio file
	fmt.Printf("Processing audio file: %s\n", *audioFilePtr)
	startTime = time.Now()
	transcription, err := recognizer.ProcessAudio(*audioFilePtr)
	if err != nil {
		fmt.Printf("Error processing audio: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Processing completed in %.2f seconds\n", time.Since(startTime).Seconds())

	// Save transcription to file
	if err := os.WriteFile(outputFile, []byte(transcription), 0644); err != nil {
		fmt.Printf("Error saving transcription: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transcription saved to: %s\n", outputFile)

	// Print a preview of the transcription
	previewLines := strings.Split(transcription, "\n")
	if len(previewLines) > 5 {
		previewLines = previewLines[:5]
	}

	fmt.Println("\nTranscription preview:")
	fmt.Println("----------------------")
	for _, line := range previewLines {
		fmt.Println(line)
	}
	if len(strings.Split(transcription, "\n")) > 5 {
		fmt.Println("...")
	}
}

// getSpecialWordsForAccent returns a list of special words for the given accent type
func getSpecialWordsForAccent(accentType speech.AccentType) []string {
	switch accentType {
	case speech.VenezuelanSpanish:
		// Common Venezuelan Spanish words and phrases that might need special attention
		return []string{
			"Venezuela", "Caracas", "Maracaibo", "Barquisimeto",
			"pana", "chamo", "chévere", "papearse", "burda", "fino",
			"guarandinga", "guachimán", "arrocero", "arrecho", "vaina",
		}
	case speech.AmericanEnglish:
		// Common police terminology and phrases
		return []string{
			"10-4", "10-20", "Roger that", "Copy that", "Over",
			"Dispatch", "Suspect", "Vehicle", "License plate",
			"Miranda rights", "K-9 unit", "Backup", "Officer down",
			"APB", "BOLO", "Perp", "Vic", "Code blue", "Code red",
		}
	default:
		return []string{}
	}
}

// formatAccentType returns a human-readable string for the accent type
func formatAccentType(accentType speech.AccentType) string {
	switch accentType {
	case speech.VenezuelanSpanish:
		return "Venezuelan Spanish"
	case speech.AmericanEnglish:
		return "American English"
	default:
		return "Generic"
	}
}
