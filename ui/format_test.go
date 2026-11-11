package ui

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestLanguageSegments(t *testing.T) {
	tests := []struct {
		name         string
		languages    map[string]int
		wantNames    []string
		wantPercents []float64
	}{
		{
			name:      "empty",
			languages: map[string]int{},
		},
		{
			name:         "single language",
			languages:    map[string]int{"Go": 100},
			wantNames:    []string{"Go"},
			wantPercents: []float64{100},
		},
		{
			name:         "drops non-positive values",
			languages:    map[string]int{"Go": 10, "Rust": 0, "C": -5},
			wantNames:    []string{"Go"},
			wantPercents: []float64{100},
		},
		{
			name: "keeps top five and merges the rest",
			languages: map[string]int{
				"Go":         500,
				"TypeScript": 300,
				"JavaScript": 100,
				"HTML":       50,
				"CSS":        30,
				"Shell":      15,
				"Dockerfile": 5,
			},
			wantNames:    []string{"Go", "TypeScript", "JavaScript", "HTML", "CSS", "Other"},
			wantPercents: []float64{50, 30, 10, 5, 3, 2},
		},
		{
			name: "ties sort by name before merging",
			languages: map[string]int{
				"Go":     100,
				"Rust":   100,
				"C":      100,
				"Python": 100,
				"Lua":    100,
				"Zig":    100,
			},
			wantNames:    []string{"C", "Go", "Lua", "Python", "Rust", "Other"},
			wantPercents: []float64{100.0 / 6, 100.0 / 6, 100.0 / 6, 100.0 / 6, 100.0 / 6, 100.0 / 6},
		},
		{
			name: "five or fewer gets no other segment",
			languages: map[string]int{
				"Go": 5, "Rust": 4, "C": 3, "Lua": 2, "Zig": 1,
			},
			wantNames:    []string{"Go", "Rust", "C", "Lua", "Zig"},
			wantPercents: []float64{33.333333333333336, 26.666666666666668, 20, 13.333333333333334, 6.666666666666667},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := LanguageSegments(tt.languages)
			if len(segments) != len(tt.wantNames) {
				t.Fatalf("got %d segments, want %d: %+v", len(segments), len(tt.wantNames), segments)
			}

			sum := 0.0
			for i, segment := range segments {
				if segment.Name != tt.wantNames[i] {
					t.Errorf("segment %d name = %q, want %q", i, segment.Name, tt.wantNames[i])
				}
				if math.Abs(segment.Percent-tt.wantPercents[i]) > 1e-9 {
					t.Errorf("segment %d percent = %v, want %v", i, segment.Percent, tt.wantPercents[i])
				}
				if segment.Color == "" {
					t.Errorf("segment %d has no color", i)
				}
				sum += segment.Percent
			}
			if len(segments) > 0 && math.Abs(sum-100) > 1e-9 {
				t.Errorf("percent total = %v, want 100", sum)
			}
		})
	}
}

func TestLanguageColor(t *testing.T) {
	if got := LanguageColor("Go"); got != languageColors["Go"] {
		t.Errorf("known language color = %q, want %q", got, languageColors["Go"])
	}
	first := LanguageColor("Some Unknown Language")
	second := LanguageColor("Some Unknown Language")
	if first != second {
		t.Errorf("fallback colors are not stable: %q != %q", first, second)
	}
	if first == "" {
		t.Error("fallback color is empty")
	}
}

func TestShortSHA(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"abc", "abc"},
		{"1234567", "1234567"},
		{"1234567890abcdef", "1234567"},
	}
	for _, tt := range tests {
		if got := ShortSHA(tt.in); got != tt.want {
			t.Errorf("ShortSHA(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestTimeAgoAt(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		then time.Time
		want string
	}{
		{"future", now.Add(time.Minute), "just now"},
		{"seconds", now.Add(-30 * time.Second), "just now"},
		{"minutes", now.Add(-5 * time.Minute), "5m ago"},
		{"hours", now.Add(-3 * time.Hour), "3h ago"},
		{"days", now.Add(-50 * time.Hour), "2d ago"},
		{"last relative day", now.Add(-29 * 24 * time.Hour), "29d ago"},
		{"older than a month", now.Add(-31 * 24 * time.Hour), "2026-10-11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeAgoAt(tt.then, now); got != tt.want {
				t.Errorf("timeAgoAt(%v) = %q, want %q", tt.then, got, tt.want)
			}
		})
	}
}

func TestCommitMessage(t *testing.T) {
	if got := CommitMessage("  fix the thing  \n\nlong body"); got != "fix the thing" {
		t.Errorf("multiline commit = %q, want first line only", got)
	}

	long := strings.Repeat("a", 100)
	got := CommitMessage(long)
	if len([]rune(got)) != maxCommitMessage {
		t.Errorf("truncated length = %d, want %d", len([]rune(got)), maxCommitMessage)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("truncated message %q does not end in an ellipsis", got)
	}

	unicode := strings.Repeat("å", 100)
	unicodeGot := CommitMessage(unicode)
	if len([]rune(unicodeGot)) != maxCommitMessage {
		t.Errorf("unicode truncated length = %d, want %d", len([]rune(unicodeGot)), maxCommitMessage)
	}
}

func TestRepoIndex(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{1, "01"},
		{10, "10"},
		{255, "255"},
		{0, "00"},
		{-1, "00"},
	}
	for _, tt := range tests {
		if got := RepoIndex(tt.in); got != tt.want {
			t.Errorf("RepoIndex(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
