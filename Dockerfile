FROM golang:1.23.2
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /exporter ./cmd/main.go
CMD ["/exporter"]