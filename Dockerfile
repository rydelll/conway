FROM golang:1.22.2-alpine3.19 AS build
WORKDIR /app

ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GOCACHE=/go/cache \
    GOMODCACHE=/go/modcache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/modcache \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/cache \
    go build -o conway cmd/conway/main.go

FROM scratch
USER 65535:65535
WORKDIR /app

COPY --from=build /app/conway /app/conway

EXPOSE 8080
ENTRYPOINT ["/app/conway"]
