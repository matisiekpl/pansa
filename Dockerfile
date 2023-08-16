FROM golang AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

FROM alpine:3.8
RUN apk add pdftk
COPY --from=builder /app .
CMD ["./main"]