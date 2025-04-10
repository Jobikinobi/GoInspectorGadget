package interview

import (
	"fmt"
	"time"
)

// InterviewType categorizes the interview context
type InterviewType string

const (
	TypeWitness  InterviewType = "WITNESS"
	TypeSuspect  InterviewType = "SUSPECT"
	TypeVictim   InterviewType = "VICTIM"
	TypeExpert   InterviewType = "EXPERT"
	TypeInformer InterviewType = "INFORMER"
)

// Interview represents an interview in an investigation
type Interview struct {
	ID             string
	CaseID         string
	Title          string
	InterviewType  InterviewType
	InterviewerID  string // ID of interviewer
	IntervieweeID  string // ID of interviewee
	Date           time.Time
	Location       string
	Status         string // e.g., "SCHEDULED", "COMPLETED", "CANCELLED"
	Duration       time.Duration
	RecordingPath  string
	Notes          string
	MediaType      string // e.g., "Audio", "Video", "Written"
	TranscriptID   string // ID of associated transcript
	KeyPoints      []string
	IsConfidential bool
	CreatedAt      time.Time
	CreatedBy      string
	UpdatedAt      time.Time
}

// Transcript represents a transcription of an interview
type Transcript struct {
	ID          string
	InterviewID string
	Content     string
	Language    string
	IsAutomated bool
	Segments    []Segment
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Segment represents a segment of a transcript
type Segment struct {
	SpeakerRole string
	StartTime   time.Duration
	EndTime     time.Duration
	Text        string
	Confidence  float64
}

// SpeechRecognitionOptions contains options for speech recognition
type SpeechRecognitionOptions struct {
	Language     string
	MultiSpeaker bool
	Diarization  bool
}

// SpeechRecognizer provides interface for speech recognition
type SpeechRecognizer interface {
	Initialize() error
	Transcribe(audioPath string, options SpeechRecognitionOptions) (*Transcript, error)
	Close() error
}

// InterviewRepository defines the interface for interview storage
type InterviewRepository interface {
	Save(i *Interview) error
	Find(id string) (*Interview, error)
	FindByCase(caseID string) ([]*Interview, error)
	Search(query string) ([]*Interview, error)
	Update(i *Interview) error
	Delete(id string) error
}

// TranscriptRepository defines the interface for transcript storage
type TranscriptRepository interface {
	Save(t *Transcript) error
	Find(id string) (*Transcript, error)
	FindByInterview(interviewID string) (*Transcript, error)
	Update(t *Transcript) error
	Delete(id string) error
}

// InterviewService provides business logic for interview management
type InterviewService struct {
	interviewRepo  InterviewRepository
	transcriptRepo TranscriptRepository
	recognizer     SpeechRecognizer
}

// NewInterviewService creates a new interview service
func NewInterviewService(interviewRepo InterviewRepository, transcriptRepo TranscriptRepository, recognizer SpeechRecognizer) *InterviewService {
	return &InterviewService{
		interviewRepo:  interviewRepo,
		transcriptRepo: transcriptRepo,
		recognizer:     recognizer,
	}
}

// CreateInterview creates a new interview
func (s *InterviewService) CreateInterview(i *Interview) error {
	if i.ID == "" {
		i.ID = generateID()
	}

	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now

	return s.interviewRepo.Save(i)
}

// GetInterview retrieves an interview by ID
func (s *InterviewService) GetInterview(id string) (*Interview, error) {
	return s.interviewRepo.Find(id)
}

// TranscribeInterview transcribes an interview using speech recognition
func (s *InterviewService) TranscribeInterview(interviewID string, options SpeechRecognitionOptions) (*Transcript, error) {
	// Retrieve the interview
	interview, err := s.interviewRepo.Find(interviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to find interview: %w", err)
	}

	if interview.RecordingPath == "" {
		return nil, fmt.Errorf("interview has no recording path")
	}

	// Transcribe the audio
	transcript, err := s.recognizer.Transcribe(interview.RecordingPath, options)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe interview: %w", err)
	}

	// Set interview ID
	transcript.InterviewID = interviewID

	// Save the transcript
	if err := s.transcriptRepo.Save(transcript); err != nil {
		return nil, fmt.Errorf("failed to save transcript: %w", err)
	}

	// Update the interview with the transcript ID
	interview.TranscriptID = transcript.ID
	interview.UpdatedAt = time.Now()

	if err := s.interviewRepo.Update(interview); err != nil {
		return nil, fmt.Errorf("failed to update interview: %w", err)
	}

	return transcript, nil
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("INT-%d", time.Now().UnixNano())
}
