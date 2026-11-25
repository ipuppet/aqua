FROM golang:1.27-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE=mono
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./services/${SERVICE}/cmd

FROM alpine:3.22
RUN addgroup -S app && adduser -S app -G app \
    && mkdir -p /app/server/storage /app/server/config \
    && chown -R app:app /app
WORKDIR /app/server
COPY --from=build /out/server ./server
COPY scripts/entrypoint.sh ./scripts/entrypoint.sh
ENV SERVICE_BIN=/app/server/server
USER app
ENTRYPOINT ["sh", "/app/server/scripts/entrypoint.sh"]
