# File Chunker for AI Processing

A fast and flexible Go tool to split large files into smaller chunks suitable for AI processing and analysis. Perfect for handling massive codebases, documents, and text files that exceed AI context limits.

## 🚀 Features

- **Multiple Chunking Strategies**: Split by lines, characters, or tokens
- **Smart Overlap**: Maintain context between chunks with configurable overlap
- **Boundary Respect**: Character chunking respects word boundaries
- **Metadata Headers**: Optional metadata with source info and chunk ranges
- **Flexible Output**: Configurable output directories and file naming
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Zero Dependencies**: Pure Go implementation

## 📦 Installation

### Pre-built Binaries
Download the latest release from the [Releases](https://github.com/admiralhr99/fileChunker/releases) page.

### Build from Source
```bash
git clone https://github.com/admiralhr99/fileChunker.git
cd fileChunker
go build -o file-chunker main.go
```

### Install with Go
```bash
go install github.com/admiralhr99/fileChunker@latest
```

## 🛠️ Usage

### Basic Examples

```bash
# Chunk a large JavaScript file by lines
./file-chunker -input large_script.js -type lines -size 1000 -overlap 50

# Chunk a document by characters with word boundaries
./file-chunker -input document.txt -type chars -size 4000 -overlap 200

# Chunk by estimated tokens for AI context limits
./file-chunker -input code.py -type tokens -size 1500 -overlap 100
```

### Advanced Usage

```bash
# Custom output directory and file prefix
./file-chunker -input massive_codebase.js \
               -output ./ai_chunks \
               -prefix "codebase" \
               -type lines \
               -size 800 \
               -overlap 40

# No overlap, no metadata headers
./file-chunker -input data.txt -overlap 0 -metadata false
```

## 📋 Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-input` | Input file to chunk (required) | - |
| `-output` | Output directory for chunks | `chunks` |
| `-type` | Chunking strategy: `lines`, `chars`, or `tokens` | `lines` |
| `-size` | Size of each chunk | `1000` |
| `-overlap` | Overlap size between chunks | `50` |
| `-metadata` | Add metadata headers to chunks | `true` |
| `-prefix` | Prefix for output filenames | Input filename |

## 🎯 Chunking Strategies

### Lines (`-type lines`)
- **Best for**: Source code, structured text files
- **Unit**: Number of lines per chunk
- **Use case**: Breaking down large codebases for AI code review

### Characters (`-type chars`)
- **Best for**: Plain text, documentation, books
- **Unit**: Number of characters per chunk
- **Features**: Respects word boundaries to avoid cutting words
- **Use case**: Processing large documents while maintaining readability

### Tokens (`-type tokens`)
- **Best for**: AI processing with strict token limits
- **Unit**: Estimated tokens (whitespace + punctuation splitting)
- **Use case**: Preparing text for language models with specific context windows

## 📁 Output Format

The tool creates numbered chunk files in the specified output directory:

```
chunks/
├── myfile_chunk_001.txt
├── myfile_chunk_002.txt
├── myfile_chunk_003.txt
└── ...
```

Each chunk includes optional metadata headers:
```
=== CHUNK 1 ===
Source: large_script.js
Lines: 1-1000
Total lines in chunk: 1000
=== CONTENT ===

[actual file content here]
```

## 💡 Real-World Examples

### Processing a Large Codebase (100k+ lines)
```bash
# Chunk for AI code review with context preservation
./file-chunker -input massive_app.js -type lines -size 800 -overlap 40
```

### Preparing Documentation for AI Summarization
```bash
# Chunk by characters to fit AI context windows
./file-chunker -input user_manual.txt -type chars -size 3500 -overlap 200
```

### Analyzing Log Files
```bash
# Chunk log files for AI analysis
./file-chunker -input application.log -type lines -size 500 -overlap 25
```

## 🔧 Integration Examples

### With Claude/ChatGPT
```bash
# Chunk code for AI review (fits most context windows)
./file-chunker -input src/app.js -type tokens -size 1500 -overlap 100
```

### With Local AI Models
```bash
# Smaller chunks for resource-constrained models
./file-chunker -input data.txt -type chars -size 2000 -overlap 150
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Issues and Feature Requests

Found a bug or have a feature request? Please open an issue on the [Issues](https://github.com/yourusername/file-chunker/issues) page.

## ⭐ Show Your Support

If this tool helps you, please consider giving it a star on GitHub! It helps others discover the project.

---

**Perfect for**: AI developers, data scientists, content creators, and anyone working with large files that need to be processed by AI systems with context limitations.