FROM golang:1.24.1

WORKDIR /app

COPY src/ /app/src
COPY db/ /app/db

WORKDIR /app/src

RUN go mod tidy
EXPOSE 8080
RUN go build -o backend .

CMD ["./backend"]
