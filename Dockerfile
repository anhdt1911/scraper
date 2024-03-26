# syntax=docker/dockerfile:1

FROM golang:1.21.6 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd
COPY internal/ internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /scraper ./cmd/scraper

# Remove redundency after build binary
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /scraper /scraper

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT [ "/scraper" ]