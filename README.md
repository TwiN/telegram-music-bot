# telegram-music-bot

[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/telegram-music-bot)](https://cloud.docker.com/repository/docker/twinproduction/telegram-music-bot)

This is a minimal music bot for Telegram.

It uses `youtube-dl` to search and download the video as well as `ffmpeg` to extract the audio.

This bot is very similar to [TwinProduction/discord-music-bot](https://github.com/TwinProduction/discord-music-bot), 
the main difference being that this one does not require streaming, but only uploading the music file to Telegram, 
making it much less complex.


## Usage

| Environment variable | Description | Required | Default |
| --- | --- | --- | --- |
| TELEGRAM_BOT_TOKEN | Discord bot token | yes | `""` |
| MAXIMUM_AUDIO_DURATION_IN_SECONDS | Maximum duration of audio clips in second | no | `480` |


## Prerequisites

If you want to run it locally, you'll need the following applications:
- youtube-dl
- ffmpeg


## Docker

### Pulling from Docker Hub

```
docker pull twinproduction/telegram-music-bot
```


### Building image locally

Building the Docker image is done as following:

```
docker build . -t telegram-music-bot
```

You can then run the container with the following command:

```
docker run -e TELEGRAM_BOT_TOKEN=secret --name telegram-music-bot telegram-music-bot
```
