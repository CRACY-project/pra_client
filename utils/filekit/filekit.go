package filekit

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/utils"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrInvalidHash          = errors.New("hash check failed")
)

type File struct {
	SourceName string
	SourcePath string
	DestName   string
	DestPath   string
}

type FileManager struct {
	SignalServer string
	HTTPClient   *http.Client
}

const timeOut = 2 * time.Second

func NewFileManager(signalServer string) *FileManager {
	return &FileManager{
		SignalServer: signalServer,
		HTTPClient:   utils.GetHttpClient(timeOut),
	}
}
func (fm *FileManager) DownloadAndVerifyFile(file File) error {
	hashFileName := strings.TrimSpace(file.SourceName) + ".sha1"
	tempFileName := file.DestName + ".tmp"
	fullDestPath := filepath.Join(file.DestPath, file.DestName)
	tempDestPath := filepath.Join(file.DestPath, tempFileName)
	hashFullDestPath := filepath.Join(file.DestPath, hashFileName)

	// Download the hash file
	if err := fm.downloadFile(hashFileName, file.SourcePath, hashFullDestPath); err != nil {
		return fmt.Errorf("error downloading hash file: %w", err)
	}
	defer os.Remove(hashFullDestPath)

	// Check if the file exists and compare its hash
	if _, err := os.Stat(fullDestPath); err == nil {
		// File exists, verify its hash
		isValid, err := CheckFileHash(fullDestPath, hashFullDestPath)
		if err != nil {
			return fmt.Errorf("error checking existing file hash: %w", err)
		}
		if isValid {
			logging.Log(logging.INFO, "File", file.DestName, "already exists and hash matches, skipping download")
			return nil
		}
		logging.Log(logging.INFO, "File", file.DestName, "exists but hash does not match, replacing")
	}

	// Download the actual file to a temporary location
	if err := fm.downloadFile(file.SourceName, file.SourcePath, tempDestPath); err != nil {
		os.Remove(tempDestPath)
		return fmt.Errorf("error downloading file: %w", err)
	}

	// Check the hash of the temporary file
	isValid, err := CheckFileHash(tempDestPath, hashFullDestPath)
	if err != nil {
		os.Remove(tempDestPath)
		return fmt.Errorf("error checking file hash: %w", err)
	}
	if !isValid {
		os.Remove(tempDestPath)
		return fmt.Errorf("hash check failed for file %s", file.DestName)
	}

	// Replace the destination file with the verified file
	if err := os.Remove(fullDestPath); err != nil && !os.IsNotExist(err) {
		os.Remove(tempDestPath)
		return fmt.Errorf("error removing existing file: %w", err)
	}
	if err := os.Rename(tempDestPath, fullDestPath); err != nil {
		os.Remove(tempDestPath)
		return fmt.Errorf("error replacing file: %w", err)
	}

	logging.Log(logging.INFO, "File", file.DestName, "downloaded and verified successfully")
	return nil
}

func (fm *FileManager) downloadFile(fileName, sourcePath, destPath string) error {
	logging.Log(logging.INFO, "Downloading file:", fileName)
	baseURL, _ := url.Parse(fm.SignalServer)
	fullURL, _ := baseURL.Parse(path.Join(sourcePath, fileName))

	resp, err := fm.HTTPClient.Get(fullURL.String())
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
	// Create the file and write the body to it
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func CheckFileHash(filePath, hashPath string) (bool, error) {
	calculatedHash, err := CalculateFileHash(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to calculate hash of file %s: %w", filePath, err)
	}
	calculatedHashHex := hex.EncodeToString(calculatedHash[:])

	// Read the expected hash from the downloaded file
	expectedHash, err := os.ReadFile(hashPath)
	if err != nil {
		return false, fmt.Errorf("failed to read hash file %s: %w", hashPath, err)
	}
	expectedHash = bytes.TrimSpace(expectedHash)

	return calculatedHashHex == string(expectedHash), nil
}

func CalculateFileHash(filePath string) ([20]byte, error) {
	binaryFile, err := os.ReadFile(filePath)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to read file: %w", err)
	}

	return sha1.Sum(binaryFile), nil
}
