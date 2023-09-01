FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o main ./main.go

FROM alpine
WORKDIR /app
EXPOSE 8080
COPY --from=builder /app/main /app
COPY --from=builder /app/.env /app
CMD [ "./main" ]