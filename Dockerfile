FROM golang:1.23

WORKDIR /app
COPY . .

RUN go build -o analytics ./cmd/analytics

EXPOSE 44044

CMD ["./analytics"]