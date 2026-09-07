package main

import (
	"github.com/0xViva/webpage/components"
	"github.com/0xViva/webpage/github"
	"github.com/0xViva/webpage/views"
	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
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
	title := "augustg.dev | AJG's Home"
	return render(c, views.Home(title, name))

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
