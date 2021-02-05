# Build the Go API
FROM golang:1.15 AS builder
ADD . /app
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main ./cmd
COPY config/app.env /app.env

## Build the React app
#FROM node:alpine3.10 AS node_builder
#COPY --from=builder /app/frontend ./
## node-sass needs python & others :(
#RUN apk --no-cache --virtual build-dependencies add \
#        python \
#        make \
#        g++ \
#        bash
#RUN yarn install
#RUN yarn build
#RUN apk del build-dependencies

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main /app.env ./
#COPY --from=node_builder /build ./frontend/build
RUN chmod +x ./main
EXPOSE 8080
CMD ./main
