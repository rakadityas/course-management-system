# Build Stage
FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o course-management-system ./cmd

# Runtime Stage
FROM alpine:latest  
WORKDIR /root/

COPY --from=build /app/course-management-system .

RUN chmod +x /root/course-management-system

CMD ["./course-management-system"]
