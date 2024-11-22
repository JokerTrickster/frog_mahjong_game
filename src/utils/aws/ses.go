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
	emailTypePassword         = emailType("password")
	emailTypeReport           = emailType("report")
	emailTypeSignup           = emailType("signup")
	emailTypeCardInfo         = emailType("cardInfo")
	emailTypeCardUploadReport = emailType("cardUploadReport")
)

const (
	ReportSpamAbuse      = 1 + iota // 도배 및 불건전한 언어 사용
	ReportIllegalProgram            // 불법 프로그램 사용
	ReportBadManners                // 비매너 행위
	ReportETC                       // 기타
)

var ReportMap = map[string]int{
	"도배 및 불건전한 언어 사용": ReportSpamAbuse,
	"불법 프로그램 사용":      ReportIllegalProgram,
	"비매너 행위":          ReportBadManners,
	"기타":              ReportETC,
}
var ReportReverseMap = make(map[int]string)

type sesMailData struct {
	email        []string
	mailType     emailType
	failCount    uint8
	templateData string
	templateName string
}
type ReqReportSES struct {
	UserID       string
	TargetUserID string
	CategoryID   string
	Reason       string
}

func EmailSendCardUploadReport(email []string, successList []string, failedList []string) {
	templateDataMap := map[string][]string{
		"successList": successList,
		"failList":    failedList,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend(email, emailTypeCardUploadReport, string(templateDataJson), "cardUploadReport")
}

func EmailSendCardInfo(email []string, cardList []string) {
	templateDataMap := map[string][]string{
		"cardList": cardList,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend(email, emailTypeCardInfo, string(templateDataJson), "cardInfo")
}

func EmailSendPassword(email string, validateCode string) {
	templateDataMap := map[string]string{
		"randomValue": validateCode,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{email}, emailTypePassword, string(templateDataJson), "password")
}

func EmailSendSignup(email string, validateCode string) {
	templateDataMap := map[string]string{
		"randomValue": validateCode,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{email}, emailTypePassword, string(templateDataJson), "signup")
}

func EmailSendReport(email []string, req *ReqReportSES) {
	templateDataMap := map[string]string{
		"userID":       req.UserID,
		"targetUserID": req.TargetUserID,
		"categoryID":   req.CategoryID,
		"reason":       req.Reason,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}
	emailSend(email, emailTypeReport, string(templateDataJson), "report")
}

func emailSend(email []string, mailType emailType, templateDataJson, templateName string) {

	mailData := sesMailData{
		email:        email,
		mailType:     mailType,
		failCount:    0,
		templateData: templateDataJson,
		templateName: templateName,
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
	//메타 정보 저장
	InitMeta()

	sesMailReqChan = make(chan sesMailData, 100)
	go func() {
		for {
			mailReq := <-sesMailReqChan
			_, err := awsClientSes.SendEmail(context.TODO(), &sesv2.SendEmailInput{
				Content: &types.EmailContent{
					Template: &types.Template{
						TemplateData: aws.String(mailReq.templateData),
						TemplateName: aws.String(mailReq.templateName),
					},
				},
				Destination: &types.Destination{
					ToAddresses: mailReq.email,
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

func GetReportID(name string) (int, error) {
	id, exists := ReportMap[name]
	if !exists {
		return 0, fmt.Errorf("시나리오 이름을 찾을 수 없습니다: %s", name)
	}
	return id, nil
}

// 맵에서 키에 해당하는 값을 가져오는 함수
func GetReportKey(val int) (string, bool) {
	key, ok := ReportReverseMap[val]
	return key, ok
}

func InitMeta() {
	for k, v := range ReportMap {
		ReportReverseMap[v] = k
	}

}
