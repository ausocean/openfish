# Build stage.
FROM golang:1.23-alpine as build-stage
WORKDIR /src

COPY cmd/openfish ./cmd/openfish
COPY datastore ./datastore
COPY go.mod go.sum ./

RUN go build ./cmd/openfish

# Production container.
FROM alpine:3.21 as production-stage
WORKDIR /app
COPY --from=build-stage /src/openfish ./
EXPOSE 8080
ENTRYPOINT [ "/app/openfish" ] 
