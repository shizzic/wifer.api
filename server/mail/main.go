package mail

import (
	"fmt"
	"time"
	"wifer/server/structs"

	"github.com/google/uuid"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Email = structs.Email
type Props = structs.Props

// Отправить код подтверждения пользователю в виде ссылки для валидации какого то действия
func SendCode(props *Props, to, code, id string) error {
	server := mail.NewSMTPClient()
	server.Host = props.Conf.EMAIL.HOST
	server.Port = props.Conf.EMAIL.PORT
	server.Username = props.Conf.EMAIL.USERNAME
	server.Password = props.Conf.EMAIL.PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.TLSConfig = &tls.props.Conf{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(props.Conf.PRODUCT_NAME + " <" + props.Conf.EMAIL.USERNAME + ">").
		AddTo(to).
		SetSubject("Confirm registration")

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+props.Conf.SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	email.SetBody(mail.TextHTML, "<p><h1>Here is a link to sign into "+props.Conf.PRODUCT_NAME+" :)</h1></p><p><a href=\""+props.Conf.CLIENT_DOMAIN+"/auth/"+id+"/"+code+"\">Enjoy</a></p>")
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}

// Отправить сообщение юзера с рабочей почты на обычную
func ContactMe(props *Props, data *structs.EmailMessage) error {
	server := mail.NewSMTPClient()
	server.Host = props.Conf.EMAIL.HOST
	server.Port = props.Conf.EMAIL.PORT
	server.Username = props.Conf.EMAIL.USERNAME
	server.Password = props.Conf.EMAIL.PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.TLSConfig = &tls.props.Conf{InsecureSkipVerify: true}
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(props.Conf.PRODUCT_NAME + " <" + props.Conf.EMAIL.USERNAME + ">").
		AddTo(props.Conf.ADMIN_EMAIL).
		SetSubject(data.Subject)

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@"+props.Conf.SELF_DOMAIN_NAME+">", msgUUID.String())
	email.AddHeader("Message-ID", msgID)

	email.SetBody(mail.TextHTML, "<p>name: "+data.Name+" - email: "+data.Sender+"</p>"+data.Message)
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}
