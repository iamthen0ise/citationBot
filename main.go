package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type MessageFrom struct {
	Id           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	Id               int64  `json:"id"`
	Title            string `json:"title"`
	Type             string `json:"type"`
	IsAdminOnlyGroup bool   `json:"all_members_are_administrators"`
}

type MessageBody struct {
	MessageID int64       `json:"message_id"`
	From      MessageFrom `json:"from"`
	Chat      Chat        `json:"chat"`
	Date      int64       `json:"date"`
	Text      string      `json:"text,omitempty"`
}

type GroupMessage struct {
	UpdateID int64       `json:"update_id"`
	Message  MessageBody `json:"message"`
}

type BotResponse struct {
	ChatId              int64  `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

func createChannelPost(targetId int64, text string, token string) bool {
	client := http.Client{}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload, err := json.Marshal(BotResponse{targetId, text, false})
	if err != nil {
		panic("Can't marshall response")
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic("Can't send request")
	}
	resp, err := client.Do(req)

	if err != nil {
		panic("Can't connect the internet")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		readBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(readBytes)
		if err != nil {
			log.Panic(err, bodyString)
		}
	}
	return true
}

func readStopWords(s3Bucket string, s3Key string) ([]string, error) {
	sess := session.Must(session.NewSession())

	s3Client := s3.New(sess)

	rawObject, err := s3Client.GetObject(
		&s3.GetObjectInput{
			Bucket: &s3Bucket,
			Key:    &s3Key,
		})

	if err != nil {
		panic("Can't read from s3")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rawObject.Body)

	var lines []string
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func findStopwordsAndSend(stopwords []string, groupMessage GroupMessage, repostChannelId int64, token string) bool {
	created := false

	for _, stopword := range stopwords {
		if strings.Contains(strings.ToLower(groupMessage.Message.Text), strings.ToLower(stopword)) {
			created = createChannelPost(int64(repostChannelId), groupMessage.Message.Text, token)
		}
	}
	return created
}

func main() {
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		s3Bucket, ok := os.LookupEnv("STOPWORDS_S3_BUCKET")
		if !ok {
			log.Panic("Can't read STOPWORDS_S3_BUCKET")
		}
		s3Key, ok := os.LookupEnv("STOPWORDS_S3_KEY")
		if !ok {
			log.Panic("Can't read STOPWORDS_S3_KEY")
		}
		stopwords, err := readStopWords(s3Bucket, s3Key)
		if err != nil {
			log.Panic("Can't read Stopwords")
		}

		token, ok := os.LookupEnv("BOT_TOKEN")
		if !ok {
			return events.APIGatewayProxyResponse{Body: "ok", StatusCode: 200}, nil
		}

		repostChannel, ok := os.LookupEnv("REPOST_CHANNEL_ID")
		if !ok {
			return events.APIGatewayProxyResponse{Body: "ok", StatusCode: 200}, nil
		}

		repostChannelId, err := strconv.ParseInt(repostChannel, 10, 64)
		if err != nil {
			panic("Can't convert Repost Channel ID to int")
		}

		var groupMessage GroupMessage
		json.Unmarshal([]byte(req.Body), &groupMessage)

		findStopwordsAndSend(stopwords, groupMessage, repostChannelId, token)

		return events.APIGatewayProxyResponse{Body: "ok", StatusCode: 200}, nil
	})
}
