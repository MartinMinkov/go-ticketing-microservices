FROM golang:1.20 as dev_builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build make dev

CMD ["./bin/payments"]

