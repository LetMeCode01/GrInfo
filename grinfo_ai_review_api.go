package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type GrInfoAIWrongQuestion struct {
	Text          string `json:"text"`
	Category      string `json:"category"`
	UserAnswer    string `json:"userAnswer"`
	CorrectAnswer string `json:"correctAnswer"`
}

type GrInfoAIReviewRequest struct {
	CurrentElo     float64                 `json:"currentElo"`
	WrongQuestions []GrInfoAIWrongQuestion `json:"wrongQuestions"`
}

type geminiCandidatePart struct {
	Text string `json:"text"`
}

type geminiCandidate struct {
	Content struct {
		Parts []geminiCandidatePart `json:"parts"`
	} `json:"content"`
}

type geminiGenerateResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func apiGrInfoAIReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GrInfoAIReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.WrongQuestions) == 0 {
		response := map[string]string{
			"review": fmt.Sprintf(
				"Incurajare: Ai un ELO de %.1f. Felicitari pentru progres!\n\nPuncte slabe: Nu am detectat raspunsuri gresite in acest quiz, continua ritmul bun.\n\nRecomandari:\n- [Grafuri orientate explicate](https://www.youtube.com/results?search_query=grafuri+orientate+explicate)\n- [Grafuri neorientate exercitii](https://www.youtube.com/results?search_query=grafuri+neorientate+exercitii)",
				req.CurrentElo,
			),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		jsonError(w, "GEMINI_API_KEY is not configured", http.StatusInternalServerError)
		return
	}

	questionLines := make([]string, 0, len(req.WrongQuestions))
	for i, q := range req.WrongQuestions {
		if i >= 20 {
			break
		}
		line := fmt.Sprintf(
			"%d) Intrebare: %q | Categorie: %q | Raspuns utilizator: %q | Raspuns corect: %q",
			i+1,
			strings.TrimSpace(q.Text),
			strings.TrimSpace(q.Category),
			strings.TrimSpace(q.UserAnswer),
			strings.TrimSpace(q.CorrectAnswer),
		)
		questionLines = append(questionLines, line)
	}

	prompt := fmt.Sprintf(`
	
Ești un profesor universitar carismatic, empatic și mentor de informatică pentru platforma educațională GrInfo.
  Utilizatorul a terminat un quiz de grafuri și are un ELO actual de %.1f.

  Scrie o recenzie personalizată, fluidă și caldă, respectând cu strictețe aceste reguli:
  1. NU folosi sub nicio formă titluri sau secțiuni rigide precum "Încurajare:", "Puncte slabe:" sau "Recomandări:". Textul trebuie să fie un mesaj unitar, continuu, ca un email personal de feedback.
  2. Comentează scurt scorul lui ELO (explică-i prietenos dacă este un nivel de start - sub 1200, intermediar - 1200-1499 sau avansat - 1500+ în interpretarea grafurilor și ce indică acest scor despre parcursul lui).
  3. Analizează-i întrebările greșite listate mai jos și explică-i, pe un ton constructiv, ce erori de logică a făcut și ce concepte specifice l-au încurcat. Păstrează explicațiile concentrate pe informatică, fără încurajări exagerate sau text inutil.
  4. La finalul textului, oferă-i direct în paragrafe un număr dinamic de recomandări (între 2 și 5 link-uri, în funcție de câte subiecte unice a greșit). 

Reguli absolute de formatare pentru link-uri:
- Fiecare recomandare trebuie să fie un hyperlink Markdown valid, care trimite către o căutare pe YouTube pentru acel subiect greșit. 
- Structura URL-ului trebuie să fie EXACT așa: [Nume tutorial informatică](https://youtube.com)
- Înlocuiește textul după egal (=) cu termenii extrași din greșelile lui (de exemplu: grafuri+neorientate+diametru).
- NU lăsa paranteze goale "()" și NU folosi ID-uri fixe de videoclipuri (fără /watch?v=).

Intrebari gresite:
%s
`, req.CurrentElo, strings.Join(questionLines, "\n"))

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		jsonError(w, "Failed to prepare AI request", http.StatusInternalServerError)
		return
	}

	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + url.QueryEscape(apiKey)
	httpClient := &http.Client{Timeout: 25 * time.Second}
	resp, err := httpClient.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		jsonError(w, "AI service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		fmt.Printf("Gemini API error status=%d body=%s\n", resp.StatusCode, string(raw))
		jsonError(w, "AI service returned an error", http.StatusBadGateway)
		return
	}

	var gemResp geminiGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gemResp); err != nil {
		jsonError(w, "Failed to decode AI response", http.StatusBadGateway)
		return
	}
	if gemResp.Error != nil {
		jsonError(w, "AI service error: "+gemResp.Error.Message, http.StatusBadGateway)
		return
	}

	reviewText := ""
	if len(gemResp.Candidates) > 0 {
		for _, p := range gemResp.Candidates[0].Content.Parts {
			if strings.TrimSpace(p.Text) != "" {
				if reviewText != "" {
					reviewText += "\n"
				}
				reviewText += strings.TrimSpace(p.Text)
			}
		}
	}

	if strings.TrimSpace(reviewText) == "" {
		jsonError(w, "AI response was empty", http.StatusBadGateway)
		return
	}

	sanitized := enforceSafeYouTubeLinks(reviewText, req.WrongQuestions)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"review": sanitized})
}

func enforceSafeYouTubeLinks(review string, wrongQuestions []GrInfoAIWrongQuestion) string {
	cleaned := strings.TrimSpace(review)

	// Remove non-YouTube raw URLs from text, keep markdown links handling below.
	rawURLRe := regexp.MustCompile(`https?://[^\s)]+`)
	cleaned = rawURLRe.ReplaceAllStringFunc(cleaned, func(u string) string {
		if isSafeYouTubeSearchURL(u) {
			return u
		}
		return ""
	})

	// Calculate dynamic number of recommendations: 1 min, 5 max, based on wrong questions count.
	maxRecs := 5
	if len(wrongQuestions) < 5 {
		maxRecs = len(wrongQuestions)
	}
	if maxRecs < 1 {
		maxRecs = 1
	}

	safeLinks := extractSafeMarkdownYouTubeLinks(cleaned)
	if len(safeLinks) < maxRecs {
		fallback := buildFallbackYouTubeLinks(wrongQuestions)
		for _, f := range fallback {
			if len(safeLinks) >= maxRecs {
				break
			}
			safeLinks = append(safeLinks, f)
		}
	}

	if !strings.Contains(cleaned, "Recomandari:") {
		cleaned += "\n\nRecomandari:"
	}

	// Replace existing markdown links block with a safe one.
	lines := strings.Split(cleaned, "\n")
	output := make([]string, 0, len(lines)+3)
	for _, line := range lines {
		if strings.Contains(line, "](http") || strings.Contains(line, "](https") {
			continue
		}
		output = append(output, line)
	}

	for i, l := range safeLinks {
		if i >= maxRecs {
			break
		}
		output = append(output, "- "+l)
	}

	return strings.TrimSpace(strings.Join(output, "\n"))
}

func extractSafeMarkdownYouTubeLinks(text string) []string {
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]+)\)`)
	matches := linkRe.FindAllStringSubmatch(text, -1)

	out := make([]string, 0, len(matches))
	for _, m := range matches {
		title := strings.TrimSpace(m[1])
		rawURL := strings.TrimSpace(m[2])
		if title == "" || !isSafeYouTubeSearchURL(rawURL) {
			continue
		}
		out = append(out, fmt.Sprintf("[%s](%s)", title, rawURL))
	}

	return out
}

func buildFallbackYouTubeLinks(wrongQuestions []GrInfoAIWrongQuestion) []string {
	keywords := extractTopKeywords(wrongQuestions)
	if len(keywords) == 0 {
		keywords = []string{"grafuri", "algoritmi"}
	}

	links := make([]string, 0, 2)
	for i := 0; i < len(keywords) && len(links) < 2; i++ {
		kw := keywords[i]
		query := url.QueryEscape("grafuri " + kw)
		title := strings.Title(kw)
		links = append(links, fmt.Sprintf("[%s - explicatii si exercitii](https://www.youtube.com/results?search_query=%s)", title, query))
	}

	for len(links) < 2 {
		query := url.QueryEscape("grafuri bac informatica")
		links = append(links, fmt.Sprintf("[Grafuri BAC - recomandare %d](https://www.youtube.com/results?search_query=%s)", len(links)+1, query))
	}

	return links
}

func extractTopKeywords(wrongQuestions []GrInfoAIWrongQuestion) []string {
	stopwords := map[string]bool{
		"care": true, "este": true, "sunt": true, "din": true, "pentru": true, "despre": true,
		"cand": true, "cum": true, "ce": true, "la": true, "si": true, "sau": true, "ale": true,
		"unei": true, "unui": true, "intr": true, "intre": true, "prin": true, "fara": true,
		"graf": true, "grafuri": true,
	}

	sanitizeRe := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	freq := map[string]int{}

	for _, q := range wrongQuestions {
		source := strings.ToLower(strings.TrimSpace(q.Text + " " + q.Category))
		source = sanitizeRe.ReplaceAllString(source, " ")
		parts := strings.Fields(source)
		for _, p := range parts {
			if len(p) < 4 || stopwords[p] {
				continue
			}
			freq[p]++
		}
	}

	type item struct {
		keyword string
		count   int
	}
	items := make([]item, 0, len(freq))
	for k, v := range freq {
		items = append(items, item{keyword: k, count: v})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].count == items[j].count {
			return items[i].keyword < items[j].keyword
		}
		return items[i].count > items[j].count
	})

	out := make([]string, 0, 3)
	for i := 0; i < len(items) && i < 3; i++ {
		out = append(out, items[i].keyword)
	}
	return out
}

func isSafeYouTubeSearchURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "https" {
		return false
	}
	host := strings.ToLower(parsed.Host)
	if host != "youtube.com" && host != "www.youtube.com" {
		return false
	}
	if parsed.Path != "/results" {
		return false
	}
	q := parsed.Query().Get("search_query")
	return strings.TrimSpace(q) != ""
}
