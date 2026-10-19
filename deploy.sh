#!/usr/bin/env bash
# Install/refresh the 0xViva.dev Quadlet and restart it.
#
# Runs on the VPS as the `deploy` user (the GitHub Actions deploy workflow does
# this over SSH). The runtime configuration arrives as environment variables and
# is written to .env, which the Quadlet passes to the container.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UNIT="0xviva"
CONTAINER="0xViva.dev"
QUADLET_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}/containers/systemd"
CURRENT_USER="$(id -un)"

: "${IMAGE_TAG:?IMAGE_TAG must be set}"
: "${INTERNAL_PORT:?INTERNAL_PORT must be set}"
: "${EXTERNAL_PORT:?EXTERNAL_PORT must be set}"
: "${GITHUB_TOKEN:?GITHUB_TOKEN must be set}"
: "${CSS_VERSION:=dev}"

# Non-login SSH sessions do not always get these from PAM, but systemctl --user
# and the rootless Podman socket need them.
export XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"
if [[ -z "${DBUS_SESSION_BUS_ADDRESS:-}" && -S "${XDG_RUNTIME_DIR}/bus" ]]; then
    export DBUS_SESSION_BUS_ADDRESS="unix:path=${XDG_RUNTIME_DIR}/bus"
fi

# Without lingering the user manager (and therefore every Quadlet) is stopped
# when the last session ends, and the container will not come back after a
# provider reboot. This is the same marker the config playbook ensures. Fail
# instead of deploying something that looks healthy now but dies on reboot.
if [[ ! -e "/var/lib/systemd/linger/${CURRENT_USER}" ]]; then
    echo "ERROR: lingering is disabled for ${CURRENT_USER}; Quadlet units will not start after a reboot." >&2
    echo "Enable it once as an administrator: sudo loginctl enable-linger ${CURRENT_USER}" >&2
    exit 1
fi

# Fail here, with a clear message, rather than having the Quadlet start fail
# because `podman load` stored the image under a different reference.
if ! podman image exists "localhost/0xviva/0xviva.dev:${IMAGE_TAG}"; then
    echo "ERROR: image localhost/0xviva/0xviva.dev:${IMAGE_TAG} is not loaded" >&2
    exit 1
fi

# GITHUB_TOKEN is a secret (PAT) used by /browse-repos at runtime.
(
    umask 077
    cat > "${SCRIPT_DIR}/.env" <<EOF
PORT=${INTERNAL_PORT}
GITHUB_TOKEN=${GITHUB_TOKEN}
CSS_VERSION=${CSS_VERSION}
EOF
)

install -d -m 0755 "${QUADLET_DIR}"
sed \
    -e "s|@IMAGE_TAG@|${IMAGE_TAG}|g" \
    -e "s|@EXTERNAL_PORT@|${EXTERNAL_PORT}|g" \
    -e "s|@INTERNAL_PORT@|${INTERNAL_PORT}|g" \
    "${SCRIPT_DIR}/quadlet/${UNIT}.container.in" \
    > "${QUADLET_DIR}/${UNIT}.container"
chmod 0644 "${QUADLET_DIR}/${UNIT}.container"

# The compose deployment is gone; remove the file CI used to place so nobody
# accidentally runs `podman compose up` and fights the Quadlet for the same
# container name.
rm -f "${SCRIPT_DIR}/docker-compose.yml"

# Stop the Quadlet if it is already running, then remove a container left over
# from the compose-based deployment.
systemctl --user stop "${UNIT}.service" 2>/dev/null || true
podman stop --time 30 "${CONTAINER}" >/dev/null 2>&1 || true
podman rm -f "${CONTAINER}" >/dev/null 2>&1 || true

systemctl --user daemon-reload
systemctl --user start "${UNIT}.service"
systemctl --user --no-pager --lines=0 status "${UNIT}.service"
