# Citation Bot

Telegram bot which repost citations to another chat/channel by stopwords matching

## Build
```shell
make build
```

## Deploy (Using Serverless.js and AWS CloudFormation)
### Dependencies
Install serverless.js first
```shell
npm install -g serverless
```
Then deploy project

### First Deployment
Populate a text file with stop words, and deploy it onto S3.

Create .env file and set variables:
```shell
BOT_TOKEN=<Telegram bot token>
REPOST_CHANNEL_ID=<target channel or group>
STOPWORDS_S3_BUCKET=<S3 bucket for stopwords>
STOPWORDS_S3_KEY=<S3 key for stopwords>
```

Run deployment script
```shell
make deploy
```

