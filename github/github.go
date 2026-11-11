package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/0xViva/0xViva.dev/components"
)

type apiRepo struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	UpdatedAt       time.Time `json:"updated_at"`
	Visibility      string    `json:"visibility"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Language        string    `json:"language"`
	Owner           struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type repoJob struct {
	repo  apiRepo
	token string
}

func GetLatestRepos(tokens []string) ([]components.GitHubRepo, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	author := "0xViva"

	var (
		jobs      []repoJob
		seenRepos = make(map[string]struct{})
		lastErr   error
	)
	for i, token := range tokens {
		repos, err := fetchRepos(client, token)
		if err != nil {
			log.Printf("failed to fetch repos for token %d: %v", i+1, err)
			lastErr = err
			continue
		}
		for _, repo := range repos {
			// The profile repo only holds the README shown on the GitHub profile,
			// so its commits do not belong in the feed.
			if strings.EqualFold(repo.Name, repo.Owner.Login) {
				continue
			}

			key := repo.FullName
			if key == "" {
				key = repo.Owner.Login + "/" + repo.Name
			}
			if _, ok := seenRepos[key]; ok {
				continue
			}
			seenRepos[key] = struct{}{}
			jobs = append(jobs, repoJob{repo: repo, token: token})
		}
	}
	if len(jobs) == 0 {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, fmt.Errorf("no repositories found")
	}

	var (
		repoWg sync.WaitGroup
		mu     sync.Mutex
		result = make([]components.GitHubRepo, 0, len(jobs))
	)

	for _, job := range jobs {
		job := job
		repo := job.repo
		repoToken := job.token
		repoWg.Add(1)
		go func() {
			defer repoWg.Done()

			languagesReq, _ := http.NewRequest(
				"GET",
				fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", repo.Owner.Login, repo.Name),
				nil,
			)
			languagesReq.Header.Set("Authorization", "Bearer "+repoToken)

			languagesResp, err := client.Do(languagesReq)
			if err != nil {
				log.Printf("failed to fetch languages for %s: %v", repo.Name, err)
				return
			}
			defer languagesResp.Body.Close()

			languages := make(map[string]int)
			if err := json.NewDecoder(languagesResp.Body).Decode(&languages); err != nil {
				log.Printf("failed to decode languages for %s: %v", repo.Name, err)
				return
			}

			branchesReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/branches", repo.Owner.Login, repo.Name), nil)
			branchesReq.Header.Set("Authorization", "Bearer "+repoToken)

			branchesResp, err := client.Do(branchesReq)
			if err != nil {
				log.Printf("failed to fetch branches for %s: %v", repo.Name, err)
				return
			}
			defer branchesResp.Body.Close()

			var branches []struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(branchesResp.Body).Decode(&branches); err != nil {
				log.Printf("failed to decode branches for %s: %v", repo.Name, err)
				return
			}

			branchNames := make([]string, 0, len(branches))
			for _, b := range branches {
				branchNames = append(branchNames, b.Name)
			}

			var (
				branchWg        sync.WaitGroup
				commitsMu       sync.Mutex
				enrichedCommits []components.GitHubCommit
				commitsSeen     = make(map[string]struct{})
			)

			for _, branch := range branches {
				branch := branch
				branchWg.Add(1)
				go func() {
					defer branchWg.Done()

					commitReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?sha=%s&per_page=6&author=%s", repo.Owner.Login, repo.Name, branch.Name, author), nil)
					commitReq.Header.Set("Authorization", "Bearer "+repoToken)

					commitResp, err := client.Do(commitReq)
					if err != nil {
						log.Printf("failed to fetch commits for %s on %s: %v", repo.Name, branch.Name, err)
						return
					}
					defer commitResp.Body.Close()

					var baseCommits []struct {
						SHA     string `json:"sha"`
						HTMLURL string `json:"html_url"`
						Commit  struct {
							Message string `json:"message"`
							Author  struct {
								Name string    `json:"name"`
								Date time.Time `json:"date"`
							} `json:"author"`
						} `json:"commit"`
					}
					if err := json.NewDecoder(commitResp.Body).Decode(&baseCommits); err != nil {
						log.Printf("failed to decode commits for %s on %s: %v", repo.Name, branch.Name, err)
						return
					}

					for _, base := range baseCommits {
						// The API query already filters by author login. A commit can
						// carry a different display name, so do not filter again here.
						detailsReq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", repo.Owner.Login, repo.Name, base.SHA), nil)
						detailsReq.Header.Set("Authorization", "Bearer "+repoToken)

						detailsResp, err := client.Do(detailsReq)
						if err != nil {
							log.Printf("failed to fetch commit detail for %s: %v", base.SHA, err)
							continue
						}
						defer detailsResp.Body.Close()

						var detail struct {
							Stats struct {
								Additions int `json:"additions"`
								Deletions int `json:"deletions"`
							} `json:"stats"`
						}
						if err := json.NewDecoder(detailsResp.Body).Decode(&detail); err != nil {
							log.Printf("failed to decode commit detail for %s: %v", base.SHA, err)
							continue
						}

						commit := components.GitHubCommit{
							SHA:       base.SHA,
							HTMLURL:   base.HTMLURL,
							Message:   base.Commit.Message,
							Additions: detail.Stats.Additions,
							Deletions: detail.Stats.Deletions,
							Author:    base.Commit.Author,
						}

						commitsMu.Lock()
						if _, exists := commitsSeen[base.SHA]; !exists {
							enrichedCommits = append(enrichedCommits, commit)
							commitsSeen[base.SHA] = struct{}{}
						}
						commitsMu.Unlock()
					}
				}()
			}

			branchWg.Wait()

			if len(enrichedCommits) == 0 {
				return
			}

			for _, c := range enrichedCommits {
				if c.Author.Date.After(repo.UpdatedAt) {
					repo.UpdatedAt = c.Author.Date
				}
			}

			sort.Slice(enrichedCommits, func(i, j int) bool {
				return enrichedCommits[i].Author.Date.After(enrichedCommits[j].Author.Date)
			})

			if len(enrichedCommits) > 6 {
				enrichedCommits = enrichedCommits[:6]
			}

			mu.Lock()
			result = append(result, components.GitHubRepo{
				Name:        repo.Name,
				HTMLURL:     repo.HTMLURL,
				Description: repo.Description,
				UpdatedAt:   repo.UpdatedAt,
				Visibility:  repo.Visibility,
				Owner:       repo.Owner,
				Commits:     enrichedCommits,
				Branches:    branchNames,
				Languages:   languages,
				Stars:       repo.StargazersCount,
				Forks:       repo.ForksCount,
				Language:    repo.Language,
			})
			mu.Unlock()
		}()
	}

	repoWg.Wait()
	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})
	if len(result) > 3 {
		result = result[:3]
	}
	return result, nil
}

func fetchRepos(client *http.Client, token string) ([]apiRepo, error) {
	req, err := http.NewRequest("GET",
		"https://api.github.com/user/repos?sort=pushed&direction=desc&per_page=10&affiliation=owner,collaborator,organization_member", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create repos request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "personal-webpage")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api returned status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var repos []apiRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode repos response: %w", err)
	}
	return repos, nil
}
