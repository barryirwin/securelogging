# Docker image for slog server
# Created by Franco Loyola - For Noroff Final Degree Project

# Small image as base
FROM alpine:latest

# Move required files
# Bin name can change (first one in the COPY), depending in the built type, check build-docker-image.sh
COPY ./slog-server-linux-arm64.bin /app/slog-server.bin
COPY ./slog-server.conf /app/
COPY ./keys /app/keys/

# Main process for the container
WORKDIR "/app"
ENTRYPOINT ["/app/slog-server.bin"]

# CMD can be replaced on the container startup, as flags are given
# See: https://docs.docker.com/engine/reference/builder/#understand-how-cmd-and-entrypoint-interact
# Eg: docker run slog-server -password="another-fancy-password"
CMD ["-password=my-fancy-password"]