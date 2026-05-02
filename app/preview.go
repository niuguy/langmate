package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	PreviewActionReplace = "replace"
	PreviewActionCopy    = "copy"
	PreviewActionCancel  = "cancel"
)

type PreviewResult struct {
	Action string `json:"action"`
	Text   string `json:"text"`
}

type previewRequest struct {
	Original    string `json:"original"`
	Replacement string `json:"replacement"`
	Model       string `json:"model"`
}

// ShowPreview opens the native preview helper and returns the user's decision.
func ShowPreview(original, replacement string, model string) (PreviewResult, error) {
	helperPath, err := findPreviewHelper()
	if err != nil {
		return PreviewResult{}, err
	}
	fmt.Printf("Opening preview helper: %s\n", helperPath)

	payload, err := json.Marshal(previewRequest{
		Original:    original,
		Replacement: replacement,
		Model:       model,
	})
	if err != nil {
		return PreviewResult{}, fmt.Errorf("encode preview request: %w", err)
	}

	cmd := exec.Command(helperPath)
	cmd.Stdin = bytes.NewReader(payload)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return PreviewResult{}, fmt.Errorf("run preview helper: %w: %s", err, msg)
		}
		return PreviewResult{}, fmt.Errorf("run preview helper: %w", err)
	}

	var result PreviewResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return PreviewResult{}, fmt.Errorf("decode preview response: %w", err)
	}
	if result.Action == "" {
		result.Action = PreviewActionCancel
	}
	return result, nil
}

func findPreviewHelper() (string, error) {
	candidates := []string{}

	if configured := os.Getenv("LANGMATE_PREVIEW_HELPER"); configured != "" {
		candidates = append(candidates, configured)
	}

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "langmate-preview"))
	}

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, "langmate-preview"))
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("preview helper not found; rebuild the app bundle or set LANGMATE_PREVIEW_HELPER")
}
