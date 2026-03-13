package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

// EmailSender định nghĩa hành vi của một bộ gửi email
type EmailSender interface {
	SendReminderEmail(toEmail, userName, roomName, amount, dueDate string) error
	Send(toEmail, subject, body string) error
}

type smtpEmailSender struct {
	host     string
	port     int
	user     string
	password string
}

// NewSMTPEmailSender khởi tạo một bộ gửi email qua SMTP
func NewSMTPEmailSender(host string, port int, user, password string) EmailSender {
	return &smtpEmailSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}

// SendReminderEmail gửi email nhắc nhở thanh toán
func (s *smtpEmailSender) SendReminderEmail(toEmail, userName, roomName, amount, dueDate string) error {
	// 1. Khởi tạo xác thực SMTP
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	// 2. Viết nội dung Email (Template)
	// Lưu ý: Cần thêm header để email hiển thị đúng tiếng Việt và định dạng HTML
	subject := "Subject: Thông báo đóng tiền trọ tháng này\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	htmlTemplate := `
	<html>
		<body>
			<h2>Xin chào {{.UserName}},</h2>
			<p>Bạn đã nhận được thông báo thanh toán tiền trọ cho <b>{{.RoomName}}</b>.</p>
			<p>Số tiền cần thanh toán: <strong style="color:red;">{{.Amount}} VND</strong></p>
			<p>Hạn chót thanh toán: <b>{{.DueDate}}</b></p>
			<br/>
			<p>Vui lòng thanh toán đúng hạn để tránh phát sinh phí trễ hạn.</p>
			<p>Trân trọng,<br>Ban Quản Lý Tro-Go.</p>
		</body>
	</html>
	`

	// 3. Đưa dữ liệu vào Template
	t, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("lỗi parse email template: %w", err)
	}

	data := struct {
		UserName string
		RoomName string
		Amount   string
		DueDate  string
	}{
		UserName: userName,
		RoomName: roomName,
		Amount:   amount,
		DueDate:  dueDate,
	}

	var body bytes.Buffer
	body.Write([]byte(subject + mime))
	err = t.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("lỗi execute email template: %w", err)
	}

	// 4. Gửi email
	smtpAddr := fmt.Sprintf("%s:%d", s.host, s.port)
	err = smtp.SendMail(smtpAddr, auth, s.user, []string{toEmail}, body.Bytes())
	if err != nil {
		return fmt.Errorf("không thể gửi email: %w", err)
	}

	return nil
}

// Send gửi một email chung chung
func (s *smtpEmailSender) Send(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	fullSubject := fmt.Sprintf("Subject: %s\n", subject)

	msg := []byte(fullSubject + mime + body)

	smtpAddr := fmt.Sprintf("%s:%d", s.host, s.port)
	err := smtp.SendMail(smtpAddr, auth, s.user, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("không thể gửi email: %w", err)
	}

	return nil
}
