FROM golang:1.23

WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download

#EXPOSE 8080

CMD ["air", "-c", "/app/.devenv/sender/.air.toml", "--", "-config", "/app/.devenv/sender/config.toml"]