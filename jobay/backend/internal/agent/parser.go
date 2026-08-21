package agent

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

// Parser handles CV text extraction
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// ExtractText extracts text from a PDF or DOCX file
func (p *Parser) ExtractText(filePath string) (string, error) {
	if strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return p.extractPDF(filePath)
	}
	if strings.HasSuffix(strings.ToLower(filePath), ".docx") {
		return p.extractDOCX(filePath)
	}
	return "", fmt.Errorf("unsupported file format: %s", filePath)
}

func (p *Parser) extractPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open PDF: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	if buf.Len() == 0 {
		return "", fmt.Errorf("no text extracted from PDF")
	}

	return buf.String(), nil
}

func (p *Parser) extractDOCX(filePath string) (string, error) {
	data, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read DOCX: %w", err)
	}

	content := data.Editable().GetContent()
	if content == "" {
		return "", fmt.Errorf("no text extracted from DOCX")
	}

	return content, nil
}
