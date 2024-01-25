FROM golang:alpine AS build
ENV GOPROXY=https://proxy.golang.org
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN go build -o corsp

FROM alpine 
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=build /app/corsp /app/
ENTRYPOINT  ["/app/corsp"]