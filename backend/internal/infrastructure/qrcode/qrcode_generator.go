package qrcode

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	// DefaultQRCodeSize is the default size for QR code generation
	DefaultQRCodeSize = 256
	// MaxQRCodeSize is the maximum allowed size for QR code generation
	MaxQRCodeSize = 1024
	// MinQRCodeSize is the minimum allowed size for QR code generation
	MinQRCodeSize = 128
)

// QRCodeGenerator handles QR code generation
type QRCodeGenerator struct{}

// NewQRCodeGenerator creates a new QR code generator instance
func NewQRCodeGenerator() *QRCodeGenerator {
	return &QRCodeGenerator{}
}

// GenerateQRCode generates a QR code from the given content
func (q *QRCodeGenerator) GenerateQRCode(content string, size int) ([]byte, error) {
	// Validate size
	if size < MinQRCodeSize || size > MaxQRCodeSize {
		return nil, fmt.Errorf("invalid QR code size: must be between %d and %d", MinQRCodeSize, MaxQRCodeSize)
	}

	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	pngData, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PNG: %w", err)
	}

	return pngData, nil
}

// GenerateQRCodeBase64 generates a QR code and returns it as a base64 encoded string
func (q *QRCodeGenerator) GenerateQRCodeBase64(content string, size int) (string, error) {
	qrData, err := q.GenerateQRCode(content, size)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(qrData), nil
}

// GenerateQRCodeDataURI generates a QR code and returns it as a data URI
func (q *QRCodeGenerator) GenerateQRCodeDataURI(content string, size int) (string, error) {
	qrData, err := q.GenerateQRCode(content, size)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrData), nil
}
