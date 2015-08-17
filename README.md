# BotFounder
Telegram Bot Founder, send messages easily to your bots

## How to run
#### Docker

```.sh

$ docker run --name postgres -e POSTGRES_USER=root -d postgres

$ docker build --no-cache=true -t izqui/botfounder .
$ docker run --link postgres:postgres  izqui/botfounder
```
