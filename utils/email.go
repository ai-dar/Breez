package utils

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	// Конфигурация почтового сервера
	from := "admbreez@gmail.com"
	password := "meib hyss oepz azgj"
	host := "smtp.gmail.com"
	port := 587

	// Создание сообщения
	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	// Настройка подключения к SMTP серверу
	d := gomail.NewDialer(host, port, from, password)

	// Отправка сообщения
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
