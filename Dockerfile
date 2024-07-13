# Use amd64 version of Golang as the builder
FROM --platform=linux/amd64 golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN GOARCH=amd64 go mod download
RUN GOARCH=amd64 go build -o main .

# Use amd64 version of Ubuntu as the base image for the final container
FROM --platform=linux/amd64 ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/main .
COPY global-bundle.pem /app/global-bundle.pem
EXPOSE 8080
CMD ["./main"]
