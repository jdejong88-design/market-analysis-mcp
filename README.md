# Market Analysis MCP

**Zero-cost market analysis powered by local LLMs**

Analyze competitor positioning, pricing, voice/tone, and design systems — all locally, with no API costs.

## Features

✅ **Website Analysis**
- Positioning & value proposition extraction
- ICP (Ideal Customer Profile) detection  
- Pricing/offer structure
- Voice & tone patterns

✅ **Design System Extraction**
- Color palette detection
- Typography analysis
- Spacing & layout patterns

✅ **Competitive Intelligence**
- Market positioning benchmarking
- Differentiation gaps
- Actionable recommendations

✅ **Zero Token Costs**
- Runs locally with Ollama
- No API calls required
- Your data stays private

## Quick Start

### 1. Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Ollama (auto-started via compose)

### 2. Run as MCP Server

```bash
docker-compose up market-analysis
```

Endpoint: `http://localhost:8080`

### 3. CLI Usage

```bash
go run main.go --url https://djlabs.nl

# Output: JSON market analysis
```

## API (MCP Protocol)

**Request:**
```json
{
  "method": "analyze_market",
  "params": {
    "url": "https://competitor.com"
  }
}
```

**Response:**
```json
{
  "url": "https://competitor.com",
  "positioning": {
    "value_proposition": "...",
    "ideal_customer_profile": ["..."],
    "pain_points": ["..."],
    "differentiation": "..."
  },
  "offers": [
    {"name": "Starter", "price": "$99", "description": "..."}
  ],
  "voice_tone": {
    "tone": "Professional & Direct",
    "keywords": ["..."],
    "copy_patterns": ["..."]
  },
  "design_system": {
    "primary_color": "#000000",
    "accent_color": "#0066CC",
    "background_color": "#FFFFFF",
    "font_family": "System UI"
  },
  "recommendations": ["..."],
  "analyzed_at": "2026-03-20T..."
}
```

## Integration with Claude

Add to `claude_desktop_config.json`:

```json
{
  "mcps": {
    "market-analysis": {
      "command": "docker",
      "args": ["run", "--rm", "-p", "8080:8080", "market-analysis-mcp:latest"]
    }
  }
}
```

Then use in Claude:

```
/analyze_market https://competitor.com
```

## Development

```bash
# Build
go build -o market-analysis-mcp

# Test
go test ./...

# Run locally
./market-analysis-mcp --url https://example.com

# Run server
./market-analysis-mcp --serve --port 8080
```

## Docker Build

```bash
docker build -t market-analysis-mcp .
docker run -p 8080:8080 market-analysis-mcp
```

## Advanced: With Ollama

```bash
# Start Ollama + MCP together
docker-compose up

# Run analysis with local LLM enrichment
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{"url": "https://djlabs.nl"}'
```

## GitHub Stars Goals

- ⭐ Practical: Agencies + freelancers use this daily
- ⭐ Zero-cost: No API subscriptions required
- ⭐ Transparent: Open source, easy to fork & extend
- ⭐ Fast: Single binary, ~50MB Docker image
- ⭐ Extensible: Hook into your own analysis pipeline

## Roadmap

- [ ] Ollama integration (semantic analysis)
- [ ] Sentiment analysis via local LLM
- [ ] Competitive benchmarking dashboard
- [ ] Batch analysis (100+ sites)
- [ ] Export to PDF reports

## License

MIT — Use freely, fork away.

---

**Want to contribute?** Fork, build, send a PR. We're looking for:
- Better CSS parsing
- More advanced positioning detection
- Visualization dashboard
- Benchmark data
