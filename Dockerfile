FROM golang:latest AS builder

WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 go build -o reporter ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/reporter /app/reporter
EXPOSE 8080
CMD ["/app/reporter"]