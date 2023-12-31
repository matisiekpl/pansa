FROM golang AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main github.com/matisiekpl/pansa-plan/cmd

FROM ubuntu
RUN apt-get update && apt-get install pdftk wkhtmltopdf -y
COPY --from=builder /app .
CMD ["./main"]
