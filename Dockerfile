FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o ./stress ./main.go && chmod +x ./stress

ENTRYPOINT ["./stress"]