package email

import (
"fmt"
"net/smtp"
)

type Mailer struct {
Host string
Port string
User string
Pass string
From string
}

func (m *Mailer) Send(to []string, subject, body string) error {
auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
m.From, to[0], subject, body)
return smtp.SendMail(fmt.Sprintf("%s:%s", m.Host, m.Port), auth, m.From, to, []byte(msg))
}
