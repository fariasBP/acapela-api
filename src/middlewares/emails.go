package middlewares

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

func SendEmailBody(nameTo, addressTo, bodyEmail string) error {
	// obtener variables de entorno
	nameApp, _ := os.LookupEnv("APP_NAME")
	emailApp, _ := os.LookupEnv("APP_EMAIL")
	pwdApp, _ := os.LookupEnv("PWD_EMAIL_APP")

	// estableciendo destinatario y remitente
	from := mail.Address{nameApp, emailApp}
	to := mail.Address{nameTo, addressTo}

	// estableciendo asunto y cuerpo del mensaje
	subj := "Envio del codigo de verificación"
	body := bodyEmail

	// estableciendo cabeceras
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// configurando mensaje
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// conectando a SMTP Server
	servername := "mail.privateemail.com:465"
	//host := "smtp.gmail.com"
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", emailApp, pwdApp, host)

	// configuracion TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To y From
	if err = c.Mail(from.Address); err != nil {
		return err
	}
	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// data
	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}

type varsTemplate struct { // variables del template
	Name string
	Code string
}

func SendEmailCodeTemplate(nameTo, addressTo, filePath, code string) error {
	// obtener variables de entorno
	nameApp, _ := os.LookupEnv("APP_NAME")
	emailApp, _ := os.LookupEnv("APP_EMAIL")
	pwdApp, _ := os.LookupEnv("PWD_EMAIL_APP")

	// estableciendo destinatario y remitente
	from := mail.Address{nameApp, emailApp}
	to := mail.Address{nameTo, addressTo}

	// estableciendo asunto y cuerpo del mensaje
	subj := "Envio del codigo de verificación"
	body := "Este es el cuerpo del mensaje"

	// estableciendo cabeceras
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// escribiendo las cabeceras al mensajemensaje
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// configurando y escribiendo el template en el mensaje
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return err
	}

	vars := varsTemplate{Name: to.Address, Code: code}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, vars)
	if err != nil {
		return err
	}
	message += buf.String()

	fmt.Println(emailApp)
	fmt.Println(pwdApp)
	// conectando a SMTP Server
	servername := "smtp.gmail.com:465"
	//host := "smtp.gmail.com"
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", emailApp, pwdApp, host)

	// configuracion TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To y From
	if err = c.Mail(from.Address); err != nil {
		return err
	}
	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// data
	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
