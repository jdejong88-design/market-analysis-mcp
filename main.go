package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Market represents analyzed market data
type Market struct {
	URL             string              `json:"url"`
	Positioning     Positioning         `json:"positioning"`
	Offers          []Offer             `json:"offers"`
	VoiceTone       VoiceTone           `json:"voice_tone"`
	Recommendations []string            `json:"recommendations"`
	AnalyzedAt      time.Time           `json:"analyzed_at"`
}

type Positioning struct {
	ValueProp       string   `json:"value_proposition"`
	ICP             []string `json:"ideal_customer_profile"`
	PainPoints      []string `json:"pain_points"`
	Differentiation string   `json:"differentiation"`
}

type Offer struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

type VoiceTone struct {
	Tone     string   `json:"tone"`
	Keywords []string `json:"keywords"`
}

// analyzeWebsite fetches and analyzes a website
func analyzeWebsite(url string) (*Market, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	market := &Market{
		URL:        url,
		AnalyzedAt: time.Now(),
	}

	// Extract data
	market.Positioning = extractPositioning(doc)
	market.Offers = extractOffers(doc)
	market.VoiceTone = extractVoiceTone(doc)
	market.Recommendations = generateRecommendations(market)

	return market, nil
}

func extractPositioning(doc *goquery.Document) Positioning {
	pos := Positioning{}

	// Find main value proposition (usually in hero)
	h1 := doc.Find("h1").First().Text()
	if h1 != "" {
		pos.ValueProp = h1
	}

	// Extract from paragraphs
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if len(text) > 50 && len(text) < 300 {
			pos.PainPoints = append(pos.PainPoints, text)
		}
	})

	return pos
}

func extractOffers(doc *goquery.Document) []Offer {
	var offers []Offer

	// Find pricing/offerings section
	doc.Find(".pricing-card, [class*='pricing'], [class*='offer']").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3, h4, .title").First().Text()
		price := s.Find(".price, [class*='price']").First().Text()
		desc := s.Find("p, .description").First().Text()

		if title != "" {
			offers = append(offers, Offer{
				Name:        strings.TrimSpace(title),
				Price:       strings.TrimSpace(price),
				Description: strings.TrimSpace(desc),
			})
		}
	})

	return offers
}

func extractVoiceTone(doc *goquery.Document) VoiceTone {
	vt := VoiceTone{
		Tone:     "Professional & Direct",
		Keywords: []string{},
	}

	// Extract text
	text := doc.Find("body").Text()
	words := strings.Fields(strings.ToLower(text))
	wordCount := make(map[string]int)

	for _, w := range words {
		if len(w) > 5 {
			wordCount[w]++
		}
	}

	// Top 5 words
	for w := range wordCount {
		if len(vt.Keywords) < 5 {
			vt.Keywords = append(vt.Keywords, w)
		}
	}

	return vt
}

func generateRecommendations(market *Market) []string {
	var recs []string

	if market.Positioning.ValueProp == "" {
		recs = append(recs, "Add clear value proposition in hero section")
	}

	if len(market.Offers) == 0 {
		recs = append(recs, "Clearly display pricing/offers")
	} else if len(market.Offers) == 1 {
		recs = append(recs, "Consider multiple pricing tiers for better conversion")
	}

	recs = append(recs, "Add social proof (testimonials, client count)")
	recs = append(recs, "Implement clear CTAs on each section")

	return recs
}

// MCPHandler implements Model Context Protocol
type MCPHandler struct{}

func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var req struct {
		Method string `json:"method"`
		Params struct {
			URL string `json:"url"`
		} `json:"params"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Method != "analyze_market" {
		http.Error(w, "Unknown method", http.StatusBadRequest)
		return
	}

	market, err := analyzeWebsite(req.Params.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(market)
}

func main() {
	url := flag.String("url", "", "Website URL to analyze")
	serve := flag.Bool("serve", false, "Run as MCP server")
	port := flag.String("port", "8080", "Server port")

	flag.Parse()

	if *serve {
		log.Printf("Market Analysis MCP listening on :%s\n", *port)
		log.Fatal(http.ListenAndServe(":"+*port, &MCPHandler{}))
	}

	if *url == "" {
		flag.Usage()
		return
	}

	market, err := analyzeWebsite(*url)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Output JSON
	b, _ := json.MarshalIndent(market, "", "  ")
	fmt.Println(string(b))
}
