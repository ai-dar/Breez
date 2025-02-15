package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailAttachment struct {
	Filename string
	Content  []byte
}

// Конфигурация для токенов
const secretKey = "tVRp+9OE2lRCeU/NhR65afmhE8XbArx/Bz2TKRd4l6Q="

// SendEmailWithAttachments отправляет письмо с вложениями
func SendEmailWithAttachments(to, subject, body string, attachments []EmailAttachment) error {
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

	for _, attachment := range attachments {
		mail.Attach(attachment.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	// Настройка подключения к SMTP серверу
	d := gomail.NewDialer(host, port, from, password)

	// Отправка сообщения
	return d.DialAndSend(mail)
}

// GenerateToken создаёт токен с таймстампом
func GenerateToken(email string) string {
	timestamp := time.Now().Add(24 * time.Hour).Unix() // Устанавливаем срок действия токена (24 часа)
	data := fmt.Sprintf("%s|%d", email, timestamp)

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s|%s", data, signature)
}

// VerifyToken проверяет токен и возвращает email, если токен валиден
func VerifyToken(token string) string {
	parts := strings.Split(token, "|")
	if len(parts) != 3 {
		return "" // Неверный формат токена
	}

	email := parts[0]
	timestampStr := parts[1]
	signature := parts[2]

	// Проверка истечения срока действия
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil || time.Now().Unix() > timestamp {
		return "" // Токен истёк
	}

	// Генерация подписи для проверки
	data := fmt.Sprintf("%s|%s", email, timestampStr)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "" // Подпись не совпадает
	}

	return email // Если всё проверено, возвращаем email
}
