package service

import (
	"context"
	"encoding/json"
	"fmt"

	// "strconv"

	// "math/rand"
	pb "message-service/api/helloworld/v1"
	"message-service/internal/biz"

	"github.com/nats-io/nats.go"
	gomail "gopkg.in/mail.v2"
)


type NotificationService struct {
	pb.UnimplementedNotificationServer

	uc *biz.NotificationUsecase
}

type Message1 struct {
    CustomerNumber string `json:"CustomerNumber"`
    FirstName      string `json:"FirstName"`
    Email          string `json:"Email"`
    Status         string `json:"Status"`
}

// type PaymentStatus struct {
// 	FromEmail     string
// 	ToEmail       string
// 	FromBalance   int64
// 	ToBalance     int64
// 	FromFirstName string
// 	ToFirstName   string
// 	Amount        int64
// 	AccountNumber int64
// }

// func (p PaymentStatus) int64(Amount int) string {
// 	panic("unimplemented")
// }


func NewGreeterService(uc *biz.NotificationUsecase) *NotificationService {
	return &NotificationService{uc: uc}
}
func connectToNATS() (*nats.Conn, error) {
	fmt.Println("Nats is connected successfully")
	return nats.Connect(nats.DefaultURL)
}

func (s *NotificationService) CreateNotification(ctx context.Context, req *pb.RequestNotification) (*pb.ReplyNotification, error) {
	g, err := s.uc.CreateNotification(ctx, &biz.Notification{
	CustomerNumber: req.CustomerNumber,
    FirstName:      req.FirstName,
    Email:          req.Email,
    Status:         req.Status,
	})
	fmt.Println("EMail Id For NATS:",g.Email)
	if err != nil {
		return nil, err
	}
	fmt.Println(g)
	// otp:=generate()
	// val2:=strconv.Itoa(otp)

	nc, _ := connectToNATS()
	fmt.Println("Msg service NAts Connected......")
	nc.Subscribe("Create", func(msg *nats.Msg) {
		var data Message1
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			fmt.Println("Failed to unmarshal NATS message:", err)
			return
		}
fmt.Println("Received NATS message:", string(msg.Data))
fmt.Println("Customer email from NATS message:", data.Email)

		fmt.Println("Received customer data for:", data.Email)

		message := gomail.NewMessage()
		message.SetHeader("From", "tksibi0@gmail.com")
		message.SetHeader("To", data.Email)
		
		message.SetHeader("Subject", "Started Creating account ")
		// message.SetBody("text/plain", "Your OTP is: "+val2)
		message.SetBody("text/plain", "Dear "+data.FirstName+"\n as per your Request We have created  Your customer id for You. Complete the Account Creation As soon as Possible .\n you can enjoy our bank service.")

		dialer := gomail.NewDialer("smtp.gmail.com", 587, "tksibi0@gmail.com", "bxsh fvld zhem moou")
//bxsh fvld zhem moou
		if err := dialer.DialAndSend(message); err != nil {
			fmt.Println("Failed to send email:", err)
			return
		}
		fmt.Println("Email sent to", data.Email)
	})
	return &pb.ReplyNotification{
			// Otp:int64(otp),
			Messages: "u have created the customer Id",
		},
		nil
}



// func (s *NotificationService) PaymentNotification(ctx context.Context, req *pb.RequestPaymentDetails) (*pb.ReplyPaymentDetails, error) {
// 	Amount := 0
// 	payment := &biz.Payment{
// 		PaymentId:         req.PaymentId,
// 		FromAccountId:     req.FromAccountId,
// 		ToAccountId:       req.ToAccountId,
// 		PaymentType:       req.PaymentType,
// 		Amount:            req.Amount,
// 		Currency:          req.Currency,
// 		Status:            req.Status,
// 		PaymentMethod:     req.PaymentMethod,
// 		ReferenceNumber:   req.ReferenceNumber,
// 		ExternalReference: req.ExternalReference,
// 	}

// 	g, err := s.uc.PaymentNotification(ctx, payment)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nc, _ := connectToNATS()
// 	nc.Subscribe("Create", func(msg *nats.Msg) {
// 		var data PaymentStatus
// 		if err := json.Unmarshal(msg.Data, &data); err != nil {
// 			fmt.Println("Failed to unmarshal NATS message:", err)
// 			return
// 		}
// 		fmt.Println("The nats is live")

// 		fmt.Println("Received Amount for:", data.Amount)

// 		message := gomail.NewMessage()
// 		message.SetHeader("From", "tksibi0@gmail.com")
// 		message.SetHeader("To", data.FromEmail)
// 		fmt.Println(data.FromEmail + "---->email")
// 		message.SetHeader("Subject", "Payment From Your Account")
// 		message.SetBody("text/plain", "Amount "+data.int64(Amount)+"is Debited from your account ")

// 		message.SetHeader("Subject", "Started Creating Payment ")
// 		dialer := gomail.NewDialer("smtp.gmail.com", 587, "tksibi0@gmail.com", "bxsh fvld zhem moou")

// 		if err := dialer.DialAndSend(message); err != nil {
// 			fmt.Println("Failed to send email:", err)
// 			return
// 		}
// 		fmt.Println("Email sent to", data.ToEmail)
// 	})
// 	return &pb.ReplyPaymentDetails{
// 		PaymentId:         g.PaymentId,
// 		FromAccountId:     g.FromAccountId,
// 		ToAccountId:       g.ToAccountId,
// 		PaymentType:       g.PaymentType,
// 		Amount:            g.Amount,
// 		Currency:          g.Currency,
// 		Status:            g.Status,
// 		PaymentMethod:     g.PaymentMethod,
// 		ReferenceNumber:   g.ReferenceNumber,
// 		ExternalReference: g.ExternalReference,
// 	}, nil
// }
