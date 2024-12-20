# Build stage.
FROM golang:1.23-alpine as build-stage
WORKDIR /src
COPY cmd/ ./cmd
COPY datastore ./datastore
COPY storage ./storage
COPY jobrunner ./jobrunner
COPY go.mod go.sum ./
RUN go build -o ./yt-dl-task ./cmd/tasks/youtube-download

# Production container.
FROM alpine:3.20 as production-stage
WORKDIR /app
RUN apk -U add yt-dlp
COPY --from=build-stage /src/yt-dl-task ./
ENTRYPOINT [ "/app/yt-dl-task" ] 
