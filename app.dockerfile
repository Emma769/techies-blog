FROM golang:1.22.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go get -d -v ./...

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api .

FROM alpine:latest AS prod

COPY --from=builder ./app/bin/api bin/api
COPY --from=builder ./app/migrations migrations

EXPOSE 8000

CMD ["./bin/api"]
