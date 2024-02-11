package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

func SendApprovalNotification(receiver []string) error {
	var subject string
	var role string
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	smtpServer := os.Getenv("MAIL_HOST")
	smtpPort, errParsing := strconv.ParseInt(os.Getenv("MAIL_PORT"), 10, 64)

	if errParsing != nil {
		fmt.Println("Error:", errParsing)
		return errParsing
	}

	htmlContent := fmt.Sprintf("<html><body><h4>Halo</h4><p>Dokumen P2H dengan nomor membutuhkan verifikasi dari. Anda dapat melakukan verifikasi dengan mengklik link dibawah.</p><a href='/dashboard/p2h/%s/edit'>Klik disini untuk verifikasi</a><p>Untuk dokumen tersebut mohon segera ditindaklanjuti, terima kasih.</p><sup>Jika Anda mengalami masalah saat mengeklik link Verifikasi, salin dan tempel URL di bawah ke web browser Anda: %s/dashboard/p2h/%s/edit</sup></body></html>", role, os.Getenv("APP_FE_URL"), os.Getenv("APP_FE_URL"))

	message := []byte(fmt.Sprintf("Subject: %s\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-type: text/html;\r\n\r\n"+
		"%s", subject, from, strings.Join(receiver, ", "), htmlContent))
	auth := smtp.PlainAuth("", from, password, smtpServer)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, receiver, message)

	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	return nil
}
