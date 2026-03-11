FROM golang:1.23

# Install gcc for CGO (required by go-sqlite3)
RUN apt-get update && apt-get install -y gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o payment ./cmd

EXPOSE 4040

CMD ["./payment"]