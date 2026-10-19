# Personal Webpage

A modern, responsive personal webpage built with Go, Echo, Templ, GoMail and TailwindCSS.

## Tech Stack

- [Go](https://go.dev/) - Backend server
- [Echo](https://echo.labstack.com/) - Web framework
- [Templ](https://templ.guide/) - HTML templating
- [TailwindCSS](https://tailwindcss.com/) - Styling
- [Air](https://github.com/cosmtrek/air) - Live reloading
- [Task](https://taskfile.dev/) - Task runner
- [Docker](https://www.docker.com/) - Containerization

## Development

Start the development server with hot-reload:

```sh
task dev
```

Build the project:

```sh
task build
```

Build Docker container:

```sh
task container image_name=<name> image_tag=<tag>
```

## Project Structure

```
.
├── assets/         # Static assets and SVG icons
├── components/     # Reusable UI components
├── style/         # TailwindCSS config and styles
├── views/         # Page templates
├── main.go        # Entry point
└── Dockerfile     # Container definition
```

## Deployment

Production runs as a rootless Podman **Quadlet** (`0xviva.service`, container
`0xViva.dev`) on the VPS. `.github/workflows/deploy.yaml` is manual
(`workflow_dispatch`): it builds the image, rsyncs it together with `deploy.sh`
and `quadlet/`, then over SSH loads the image and runs `deploy.sh`, which
renders the unit and restarts it. The unit's
`[Install] WantedBy=default.target` plus lingering for `deploy` makes systemd
start it after a host reboot. `podman compose` is no longer used in production.

## Docker

Built with a multi-stage build:

- Go builder stage compiles binary
- Final stage uses distroless base image
- Exposes port 8080

## License

This project is open source and available under the MIT license.

## Want my badge? No problems:

add this file to your static assets:
`assets/js/augustg-dev-badge.js`

this one here:
https://github.com/0xViva/0xViva.dev/blob/master/assets/js/augustg-dev-badge.js

or get it from this url via <script> (at your own risk):
https://raw.githubusercontent.com/0xViva/0xViva.dev/refs/heads/master/assets/js/augustg-dev-badge.js

Badge is using this font:
https://github.com/0xViva/0xViva.dev/blob/master/assets/fonts/JetBrainsMonoNerdFontMono-Regular.ttf

then add it all together in your desired html location:
```
<augustg-dev-badge></augustg-dev-badge>
<script src="/whatever/location/you/use/augustg-dev-badge.js"></script>
```

