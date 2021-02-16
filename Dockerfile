# Build the Go API
FROM golang:1.15 AS builder
ADD . /app
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /healthCheck ./cmd/healthCheck.go
COPY config/app.env /app.env

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main /app.env ./
COPY --from=builder /healthCheck ./
RUN chmod +x ./main
EXPOSE 8080
HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/healthCheck" ]
CMD ./main
