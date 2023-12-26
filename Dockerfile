FROM golang:latest

WORKDIR /app
COPY . .

RUN go build -o city_validation_s3

EXPOSE 8080

CMD ["./city_validation_s3"]
