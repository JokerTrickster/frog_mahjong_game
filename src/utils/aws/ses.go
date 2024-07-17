package aws

import (
	"context"
	"encoding/json"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"golang.org/x/exp/rand"
)

type emailType string

const (
	emailTypePassword = emailType("password")
)

type sesMailData struct {
	email        string
	validateCode string
	mailType     emailType
	failCount    uint8
	templateData string
}

var deepUrlLink string
var deepUrlRedirectSignup string
var deepUrlRedirectPassword string

func EmailSendPassword(email string, validateCode string) {

	// 랜덤 값 생성
	randomValue := generateRandomValue(6)

	fmt.Println("여기옴? ", randomValue)

	emailSend(email, validateCode, emailTypePassword, randomValue)
}

func emailSend(email string, validateCode string, mailType emailType, randomValue string) {
	templateDataMap := map[string]string{
		"randomValue": randomValue,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	mailData := sesMailData{
		email:        email,
		validateCode: validateCode,
		mailType:     mailType,
		failCount:    0,
		templateData: string(templateDataJson),
	}
	select {
	case sesMailReqChan <- mailData:
	default:
		<-sesMailReqChan
		sesMailReqChan <- mailData
	}
}

var sesMailReqChan chan sesMailData

func InitAwsSes() error {

	sesMailReqChan = make(chan sesMailData, 100)
	go func() {
		for {
			mailReq := <-sesMailReqChan
			_, err := awsClientSes.SendEmail(context.TODO(), &sesv2.SendEmailInput{
				Content: &types.EmailContent{
					Template: &types.Template{
						// TemplateData: aws.String(mailReq.templateData),
						TemplateName: aws.String("MyTemplate"),
					},
				},
				Destination: &types.Destination{
					ToAddresses: []string{mailReq.email},
				},
				EmailTags: []types.MessageTag{{
					Name:  aws.String("type"),
					Value: aws.String(string(mailReq.mailType)),
				}},
				FromEmailAddress: aws.String("pkjhj485@naver.com"),
			})
			if err != nil {
				if mailReq.failCount < 3 {
					mailReq.failCount += 1
					sesMailReqChan <- mailReq
				}
			}
		}
	}()
	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 랜덤 값 생성 함수
func generateRandomValue(length int) string {
	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
