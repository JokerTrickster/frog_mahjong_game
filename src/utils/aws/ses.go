package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type emailType string

const (
	emailTypePassword = emailType("password")
	emailTypeReport   = emailType("report")
)

type sesMailData struct {
	email        string
	validateCode string
	mailType     emailType
	failCount    uint8
	templateData string
}

func EmailSendPassword(email string, validateCode string) {

	emailSend(email, validateCode, emailTypePassword, validateCode)
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
						TemplateData: aws.String(mailReq.templateData),
						TemplateName: aws.String("password"),
					},
				},
				Destination: &types.Destination{
					ToAddresses: []string{mailReq.email},
				},
				EmailTags: []types.MessageTag{{
					Name:  aws.String("type"),
					Value: aws.String(string(mailReq.mailType)),
				}},
				FromEmailAddress: aws.String("root@jokertrickster.com"),
			})
			if err != nil {
				if mailReq.failCount < 3 {
					fmt.Println("Error sending email:", err)
					mailReq.failCount += 1
					sesMailReqChan <- mailReq
				}
			}
		}
	}()
	return nil
}
