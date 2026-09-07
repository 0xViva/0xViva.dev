package main

import (
	"github.com/0xViva/0xViva.dev/components"
	"github.com/0xViva/0xViva.dev/github"
	"github.com/0xViva/0xViva.dev/views"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"strings"
)

var (
	githubToken string
)

func main() {
	e := echo.New()

	godotenv.Load()

	githubToken = os.Getenv("GITHUB_TOKEN")

	e.Use(middleware.Logger())
	e.Static("/style", "style")
	e.Static("/assets", "assets")

	e.GET("/", homeView)
	e.GET("/browse-repos", browseRepos)
	e.Logger.Fatal(e.Start(":8080"))
}

func homeView(c echo.Context) error {
	name := "August Justinus Gran"
	title := titleForHost(requestHost(c))
	return render(c, views.Home(title, name))

}

// requestHost prefers X-Forwarded-Host (when behind a proxy) and falls back to Host.
func requestHost(c echo.Context) string {
	if fwd := c.Request().Header.Get("X-Forwarded-Host"); fwd != "" {
		// May contain a comma-separated list; the first entry is the original host.
		if h, _, ok := strings.Cut(fwd, ","); ok {
			return strings.TrimSpace(h)
		}
		return strings.TrimSpace(fwd)
	}
	return c.Request().Host
}

// titleForHost returns just the address the user visited from.
func titleForHost(host string) string {
	h := strings.ToLower(strings.TrimSpace(host))
	// Strip port if present (e.g. "augustg.dev:8080" -> "augustg.dev").
	if i := strings.LastIndex(h, ":"); i != -1 && !strings.HasSuffix(h, "]") {
		h = h[:i]
	}
	h = strings.Trim(h, "[]")
	switch {
	case strings.Contains(h, "0xviva.dev"):
		return "0xviva.dev"
	case strings.Contains(h, "augustg.dev"):
		return "augustg.dev"
	case h != "":
		return h
	default:
		return "augustg.dev"
	}
}
func browseRepos(c echo.Context) error {
	repos, err := github.GetLatestRepos(githubToken)
	if err != nil {
		c.Logger().Errorf("Failed to fetch GitHub repos: %v", err)
		return render(c, components.RepoContainer(nil))
	}
	return render(c, components.RepoContainer(repos))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}
