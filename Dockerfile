FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun .

FROM alpine:latest
WORKDIR /app
# Instalar certificados SSL
RUN apk add --no-cache ca-certificates
COPY --from=build /app/cloudrun .
COPY --from=build /app/main_test.go .
ENTRYPOINT ["./cloudrun"]
