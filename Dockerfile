# 1 шаг
FROM golang:1.24-alpine AS build_stage
WORKDIR /voting-bot
COPY . .
RUN go mod tidy
RUN go build -o binary_bot cmd/voting-bot/main.go  

# 2 шаг
FROM alpine AS run_stage
WORKDIR /bot_binary
RUN apt update && apt install -y libssl-dev pkg-config

COPY --from=build_stage /voting-bot /bot_binary/

CMD [ "./binary_bot" ]