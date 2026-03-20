package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

// Market represents analyzed market data
type Market struct {
	URL            string                 `json:"url"`
	Positioning    Positioning            `json:"positioning"`
	Offers         []Offer                `json:"offers"`
	VoiceTone      VoiceTone              `json:"voice_tone"`
	DesignSystem   DesignSystem           `json:"design_system"`
	Competitors    []CompetitorInsight    `json:"competitors"`
	Recommendations []string              `json:"recommendations"`
	AnalyzedAt     time.Time              `json:"analyzed_at"`
}

type Positioning struct {
	ValueProp      string   `json:"value_proposition"`
	ICP            []string `json:"ideal_customer_profile"`
	PainPoints     []string `json:"pain_points"`
	Differentiation string  `json:"differentiation"`
}

type Offer struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

type VoiceTone struct {
	Tone         string   `json:"tone"`
	Keywords     []string `json:"keywords"`
	CopyPatterns []string `json:"copy_patterns"`
}

type DesignSystem struct {
	PrimaryColor   string `json:"primary_color"`
	AccentColor    string `json:"accent_color"`
	BackgroundColor string `json:"background_color"`
	FontFamily     string `json:"font_family"`
}

type CompetitorInsight struct {
	Name           string `json:"name"`
	Strength       string `json:"strength"`
	Opportunity    string `json:"opportunity"`
}

// analyzeWebsite scrapes and analyzes a website
func analyzeWebsite(ctx context.Context, url string) (*Market, error) {
	// Create chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Screenshot + HTML
	var buf []byte
	var html string

	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.CaptureScreenshot(&buf),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract data
	market := &Market{
		URL:        url,
		AnalyzedAt: time.Now(),
	}

	// Positioning (h1, p tags)
	market.Positioning = extractPositioning(doc)

	// Offers (pricing section)
	market.Offers = extractOffers(doc)

	// Voice & Tone (from copy)
	market.VoiceTone = extractVoiceTone(doc)

	// Design System (CSS variables, colors)
	market.DesignSystem = extractDesignSystem(doc)

	// Generate recommendations
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

	// Extract repeated words (keywords)
	text := doc.Find("body").Text()
	words := strings.Fields(strings.ToLower(text))
	wordCount := make(map[string]int)

	for _, w := range words {
		if len(w) > 5 {
			wordCount[w]++
		}
	}

	// Top 5 words
	for w, count := range wordCount {
		if count > 3 {
			vt.Keywords = append(vt.Keywords, w)
		}
		if len(vt.Keywords) >= 5 {
			break
		}
	}

	return vt
}

func extractDesignSystem(doc *goquery.Document) DesignSystem {
	ds := DesignSystem{
		FontFamily:      "System UI",
		PrimaryColor:    "#000000",
		AccentColor:     "#0066CC",
		BackgroundColor: "#FFFFFF",
	}

	// Extract from CSS variables or computed styles
	doc.Find("[style*='--color'], [style*='color:']").First().Each(func(i int, s *goquery.Selection) {
		style, _ := s.Attr("style")
		if strings.Contains(style, "primary") {
			ds.PrimaryColor = extractColorFromStyle(style)
		}
		if strings.Contains(style, "accent") {
			ds.AccentColor = extractColorFromStyle(style)
		}
	})

	return ds
}

func extractColorFromStyle(style string) string {
	// Simple color extraction
	parts := strings.Split(style, "#")
	if len(parts) > 1 {
		hex := "#" + strings.Split(parts[1], ";")[0]
		if len(hex) >= 7 {
			return hex[:7]
		}
	}
	return "#000000"
}

func generateRecommendations(market *Market) []string {
	var recs []string

	// Based on positioning
	if market.Positioning.ValueProp == "" {
		recs = append(recs, "Add clear value proposition in hero section")
	}

	// Based on offers
	if len(market.Offers) == 0 {
		recs = append(recs, "Clearly display pricing/offers")
	} else if len(market.Offers) == 1 {
		recs = append(recs, "Consider multiple pricing tiers for better conversion")
	}

	// Based on design
	if market.DesignSystem.AccentColor == "#0066CC" {
		recs = append(recs, "Consider a distinctive accent color")
	}

	recs = append(recs, "Add social proof (testimonials, client count)")
	recs = append(recs, "Implement clear CTAs on each section")
	recs = append(recs, "Add FAQ section addressing common objections")

	return recs
}

// MCPHandler implements Model Context Protocol
type MCPHandler struct{}

func (h *MCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Method string `json:"method"`
		Params struct {
			URL string `json:"url"`
		} `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Method != "analyze_market" {
		http.Error(w, "Unknown method", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	market, err := analyzeWebsite(ctx, req.Params.URL)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	market, err := analyzeWebsite(ctx, *url)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Output JSON
	b, _ := json.MarshalIndent(market, "", "  ")
	fmt.Println(string(b))
}
