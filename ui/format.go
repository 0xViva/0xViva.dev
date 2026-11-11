package ui

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"time"
)

type LanguageSegment struct {
	Name    string
	Bytes   int
	Percent float64
	Color   string
}

const (
	maxLanguageSegments = 5
	maxCommitMessage    = 72
	daysBeforeDate      = 30
)

var languageColors = map[string]string{
	"Go":          "#00add8",
	"TypeScript":  "#3178c6",
	"JavaScript":  "#f1e05a",
	"HTML":        "#e34c26",
	"CSS":         "#7c5cbf",
	"Python":      "#3572a5",
	"Rust":        "#dea584",
	"C":           "#8b949e",
	"C++":         "#f34b7d",
	"C#":          "#178600",
	"Java":        "#b07219",
	"Kotlin":      "#a97bff",
	"Swift":       "#f05138",
	"Shell":       "#89e051",
	"Dockerfile":  "#4f7b8a",
	"SQL":         "#e38c00",
	"Nix":         "#7e7eff",
	"Lua":         "#61a5c2",
	"Zig":         "#ec915c",
	"Ruby":        "#cc342d",
	"PHP":         "#6577b3",
	"Dart":        "#00b4ab",
	"Elixir":      "#a074c4",
	"Haskell":     "#8f7ec4",
	"Vue":         "#41b883",
	"Svelte":      "#ff3e00",
	"Assembly":    "#a6824f",
	"Makefile":    "#6e8b3d",
	"PowerShell":  "#4f8bcf",
	"Objective-C": "#438eff",
	"Scala":       "#c22d40",
	"Clojure":     "#db5855",
	"OCaml":       "#3be133",
	"F#":          "#b845fc",
	"Perl":        "#0298c3",
	"R":           "#198ce7",
	"Julia":       "#a270ba",
	"Terraform":   "#844fba",
	"YAML":        "#cb171e",
	"JSON":        "#8b949e",
	"Markdown":    "#8b949e",
	"Other":       "#4a5f6a",
}

var fallbackColors = []string{
	"#00d9ff",
	"#ffb000",
	"#ff2e63",
	"#00ff9c",
	"#a78bfa",
	"#f0883e",
	"#39d353",
	"#7a919c",
}

func LanguageSegments(languages map[string]int) []LanguageSegment {
	type language struct {
		name  string
		bytes int
	}

	total := 0
	entries := make([]language, 0, len(languages))
	for name, bytes := range languages {
		if bytes <= 0 {
			continue
		}
		entries = append(entries, language{name: name, bytes: bytes})
		total += bytes
	}
	if total == 0 {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].bytes != entries[j].bytes {
			return entries[i].bytes > entries[j].bytes
		}
		return entries[i].name < entries[j].name
	})

	segments := make([]LanguageSegment, 0, maxLanguageSegments+1)
	for i, entry := range entries {
		if i == maxLanguageSegments {
			break
		}
		segments = append(segments, LanguageSegment{
			Name:    entry.name,
			Bytes:   entry.bytes,
			Percent: float64(entry.bytes) / float64(total) * 100,
			Color:   LanguageColor(entry.name),
		})
	}

	other := 0
	for _, entry := range entries[min(maxLanguageSegments, len(entries)):] {
		other += entry.bytes
	}
	if other > 0 {
		segments = append(segments, LanguageSegment{
			Name:    "Other",
			Bytes:   other,
			Percent: float64(other) / float64(total) * 100,
			Color:   LanguageColor("Other"),
		})
	}

	return segments
}

func LanguageColor(name string) string {
	if color, ok := languageColors[name]; ok {
		return color
	}

	hash := fnv.New32a()
	hash.Write([]byte(strings.ToLower(name)))
	return fallbackColors[int(hash.Sum32())%len(fallbackColors)]
}

func TimeAgo(t time.Time) string {
	return timeAgoAt(t, time.Now())
}

func timeAgoAt(t, now time.Time) string {
	elapsed := now.Sub(t)
	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		return fmt.Sprintf("%dm ago", int(elapsed.Minutes()))
	case elapsed < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(elapsed.Hours()))
	case elapsed < daysBeforeDate*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(elapsed.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}

func ShortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

func CommitMessage(message string) string {
	line := strings.TrimSpace(strings.SplitN(message, "\n", 2)[0])
	runes := []rune(line)
	if len(runes) <= maxCommitMessage {
		return line
	}
	return string(runes[:maxCommitMessage-3]) + "..."
}

func RepoIndex(i int) string {
	if i < 0 {
		i = 0
	}
	return fmt.Sprintf("%02d", i)
}
