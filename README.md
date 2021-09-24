# Citation Bot

Telegram bot which repost citations to another chat/channel by stopwords matching

## Build
```shell
make build
```

## Deploy (Using Serverless.js and AWS CloudFormation)
Install serverless.js first
```shell
npm install -g serverless
```
Then deploy project

```shell
make deploy
```

## Config
Create .env file and set variables:
```shell
BOT_TOKEN=<Telegram bot token>
REPOST_CHANNEL_ID=<target channel or group>
STOPWORDS_S3_BUCKET=<S3 bucket for stopwords>
STOPWORDS_S3_KEY=<S3 key for stopwords>
```