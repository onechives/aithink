FROM golang:1.25.5-bookworm AS builder

WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /out/server /app/server
COPY backend/conf /app/conf
COPY backend/logfile /app/logfile
COPY backend/uploads /app/uploads

EXPOSE 8080
ENTRYPOINT ["/app/server"]
CMD ["/app/conf/config.docker.yaml"]
