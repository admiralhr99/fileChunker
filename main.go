package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ChunkConfig struct {
	InputFile   string
	OutputDir   string
	ChunkType   string // "lines", "chars", "tokens"
	ChunkSize   int
	OverlapSize int
	AddMetadata bool
	Prefix      string
}

type Chunker struct {
	config ChunkConfig
}

func NewChunker(config ChunkConfig) *Chunker {
	return &Chunker{config: config}
}

func (c *Chunker) ChunkByLines() error {
	file, err := os.Open(c.config.InputFile)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var currentChunk []string
	var previousOverlap []string
	chunkNumber := 1
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Start new chunk with overlap from previous chunk
		if len(currentChunk) == 0 && len(previousOverlap) > 0 {
			currentChunk = append(currentChunk, previousOverlap...)
		}

		currentChunk = append(currentChunk, line)

		// Check if chunk is full
		if len(currentChunk) >= c.config.ChunkSize {
			if err := c.writeChunk(currentChunk, chunkNumber, lineNumber-len(currentChunk)+1, lineNumber); err != nil {
				return err
			}

			// Prepare overlap for next chunk
			if c.config.OverlapSize > 0 && len(currentChunk) > c.config.OverlapSize {
				previousOverlap = currentChunk[len(currentChunk)-c.config.OverlapSize:]
			} else {
				previousOverlap = nil
			}

			currentChunk = nil
			chunkNumber++
		}
	}

	// Write remaining lines as final chunk
	if len(currentChunk) > 0 {
		if err := c.writeChunk(currentChunk, chunkNumber, lineNumber-len(currentChunk)+1, lineNumber); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (c *Chunker) ChunkByCharacters() error {
	content, err := os.ReadFile(c.config.InputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	text := string(content)
	chunkNumber := 1
	start := 0

	for start < len(text) {
		end := start + c.config.ChunkSize
		if end > len(text) {
			end = len(text)
		}

		// Try to break at word boundary
		if end < len(text) {
			for i := end; i > start && i > end-100; i-- {
				if text[i] == ' ' || text[i] == '\n' || text[i] == '\t' {
					end = i
					break
				}
			}
		}

		chunk := text[start:end]

		if err := c.writeTextChunk(chunk, chunkNumber, start, end); err != nil {
			return err
		}

		// Move start position with overlap
		if c.config.OverlapSize > 0 {
			start = end - c.config.OverlapSize
			if start < 0 {
				start = end
			}
		} else {
			start = end
		}

		chunkNumber++
	}

	return nil
}

func (c *Chunker) ChunkByTokens() error {
	content, err := os.ReadFile(c.config.InputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Simple token approximation: split by whitespace and punctuation
	text := string(content)
	tokens := c.tokenize(text)

	chunkNumber := 1
	start := 0

	for start < len(tokens) {
		end := start + c.config.ChunkSize
		if end > len(tokens) {
			end = len(tokens)
		}

		chunkTokens := tokens[start:end]
		chunk := strings.Join(chunkTokens, " ")

		if err := c.writeTextChunk(chunk, chunkNumber, start, end); err != nil {
			return err
		}

		// Move start position with overlap
		if c.config.OverlapSize > 0 {
			start = end - c.config.OverlapSize
			if start < 0 {
				start = end
			}
		} else {
			start = end
		}

		chunkNumber++
	}

	return nil
}

func (c *Chunker) tokenize(text string) []string {
	// Simple tokenization - split on whitespace and keep punctuation
	var tokens []string
	var current strings.Builder

	for _, char := range text {
		switch {
		case char == ' ' || char == '\t' || char == '\n':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case char == '.' || char == ',' || char == ';' || char == ':' ||
			char == '!' || char == '?' || char == '(' || char == ')' ||
			char == '[' || char == ']' || char == '{' || char == '}':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(char))
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (c *Chunker) writeChunk(lines []string, chunkNumber, startLine, endLine int) error {
	filename := fmt.Sprintf("%s_chunk_%03d.txt", c.config.Prefix, chunkNumber)
	filepath := filepath.Join(c.config.OutputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating chunk file: %v", err)
	}
	defer file.Close()

	if c.config.AddMetadata {
		fmt.Fprintf(file, "=== CHUNK %d ===\n", chunkNumber)
		fmt.Fprintf(file, "Source: %s\n", c.config.InputFile)
		fmt.Fprintf(file, "Lines: %d-%d\n", startLine, endLine)
		fmt.Fprintf(file, "Total lines in chunk: %d\n", len(lines))
		fmt.Fprintf(file, "=== CONTENT ===\n\n")
	}

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	fmt.Printf("Created chunk %d: %s (lines %d-%d)\n", chunkNumber, filename, startLine, endLine)
	return nil
}

func (c *Chunker) writeTextChunk(content string, chunkNumber, start, end int) error {
	filename := fmt.Sprintf("%s_chunk_%03d.txt", c.config.Prefix, chunkNumber)
	filepath := filepath.Join(c.config.OutputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating chunk file: %v", err)
	}
	defer file.Close()

	if c.config.AddMetadata {
		fmt.Fprintf(file, "=== CHUNK %d ===\n", chunkNumber)
		fmt.Fprintf(file, "Source: %s\n", c.config.InputFile)
		fmt.Fprintf(file, "Range: %d-%d\n", start, end)
		fmt.Fprintf(file, "=== CONTENT ===\n\n")
	}

	fmt.Fprint(file, content)

	fmt.Printf("Created chunk %d: %s\n", chunkNumber, filename)
	return nil
}

func (c *Chunker) Process() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(c.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	switch c.config.ChunkType {
	case "lines":
		return c.ChunkByLines()
	case "chars":
		return c.ChunkByCharacters()
	case "tokens":
		return c.ChunkByTokens()
	default:
		return fmt.Errorf("unsupported chunk type: %s", c.config.ChunkType)
	}
}

func main() {
	var config ChunkConfig

	flag.StringVar(&config.InputFile, "input", "", "Input file to chunk (required)")
	flag.StringVar(&config.OutputDir, "output", "chunks", "Output directory for chunks")
	flag.StringVar(&config.ChunkType, "type", "lines", "Chunk type: lines, chars, or tokens")
	flag.IntVar(&config.ChunkSize, "size", 1000, "Size of each chunk")
	flag.IntVar(&config.OverlapSize, "overlap", 50, "Overlap size between chunks")
	flag.BoolVar(&config.AddMetadata, "metadata", true, "Add metadata to chunks")
	flag.StringVar(&config.Prefix, "prefix", "", "Prefix for output files (defaults to input filename)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Chunk large files for AI processing.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -input large_file.js -type lines -size 500 -overlap 25\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -input document.txt -type chars -size 4000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -input code.py -type tokens -size 1500 -output ./chunks\n", os.Args[0])
	}

	flag.Parse()

	if config.InputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Input file is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Set default prefix to input filename without extension
	if config.Prefix == "" {
		base := filepath.Base(config.InputFile)
		config.Prefix = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// Validate input file exists
	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file does not exist: %s\n", config.InputFile)
		os.Exit(1)
	}

	// Validate chunk type
	validTypes := map[string]bool{"lines": true, "chars": true, "tokens": true}
	if !validTypes[config.ChunkType] {
		fmt.Fprintf(os.Stderr, "Error: Invalid chunk type. Must be: lines, chars, or tokens\n")
		os.Exit(1)
	}

	chunker := NewChunker(config)

	fmt.Printf("Chunking file: %s\n", config.InputFile)
	fmt.Printf("Chunk type: %s\n", config.ChunkType)
	fmt.Printf("Chunk size: %d\n", config.ChunkSize)
	fmt.Printf("Overlap: %d\n", config.OverlapSize)
	fmt.Printf("Output directory: %s\n", config.OutputDir)
	fmt.Println()

	if err := chunker.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nChunking completed successfully!")
}
