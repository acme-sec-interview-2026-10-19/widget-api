FROM golang:latest

WORKDIR /app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin

COPY . .

RUN go mod download && go build -o /app/widget-api .

EXPOSE 8080

CMD ["/app/widget-api"]
