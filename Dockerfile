# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /app
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o bin/telegram-music-bot .

FROM alpine:3.7
ENV TELEGRAM_BOT_TOKEN=""
ENV MAXIMUM_AUDIO_DURATION_IN_SECONDS=""
ENV APP_HOME=/app
WORKDIR ${APP_HOME}
RUN apk --update add --no-cache ca-certificates ffmpeg python
COPY --from=builder /app/bin/telegram-music-bot ./bin/telegram-music-bot
RUN wget --no-check-certificate https://yt-dl.org/downloads/latest/youtube-dl -O /usr/bin/youtube-dl
RUN chmod +x /usr/bin/youtube-dl
RUN youtube-dl --version
ENTRYPOINT ["/app/bin/telegram-music-bot"]