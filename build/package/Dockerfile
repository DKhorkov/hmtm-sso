# https://habr.com/ru/companies/otus/articles/660301/
FROM golang:alpine AS builder
LABEL authors='dkhorkov'

WORKDIR /build

COPY . .

RUN go build -o hmtmsso ./cmd/hmtmsso/main.go

FROM alpine AS runner

WORKDIR /app

COPY .env .
COPY --from=builder /build/internal/database/migrations/ /app/migrations/
COPY --from=builder /build/hmtmsso /app/hmtmsso
RUN go get

CMD ["./hmtmsso"]
