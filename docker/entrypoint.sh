#!/bin/sh
# Entry point script for backend container
# Handles docker socket permissions for nonroot user

# Add nonroot user to docker group if docker socket exists
if [ -S /var/run/docker.sock ]; then
    addgroup nonroot docker 2>/dev/null || true
fi

# Execute the main application
exec /app/main "$@"
