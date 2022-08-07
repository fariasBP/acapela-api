package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// enviar cualquier mensajes (solo si abre la comunicacion)
func SendAnyMessageText(to, msg string) error {
	// obtener las variables de entorno
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	// estructura del mensaje
	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "text",
		"text": {
			"body": "` + msg + `"
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return err
}

// mensaje de bienvenida
func SendWelcomeMessage(to, userName string) error {
	// obtener las variables de entorno
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "template",
		"template": {
			"name": "welcome_new_user",
			"language": {
				"code": "es",
			},
			"components" : [
				{
					"type": "body",
					"parameters": [
						{
							"type": "text",
							"text": "` + userName + `",
						}
					]
				}
			]
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	if resp.StatusCode != 200 {
		return err
	}

	return err
}

// mensaje de codigo
func SendCodeMessage(to, code string) error {
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	fmt.Println("to:", to)
	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "template",
		"template": {
			"name": "verification_code",
			"language": {
				"code": "es",
			},
			"components" : [
				{
					"type": "body",
					"parameters": [
						{
							"type": "text",
							"text": "` + code + `",
						}
					]
				}
			]
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return err
}

func SendNewProduct(codePhone, phone int, userName, kindName, forGender string) error {
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	codePhoneStr, phoneStr := strconv.Itoa(codePhone), strconv.Itoa(phone)

	to := codePhoneStr + phoneStr

	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "template",
		"template": {
			"name": "nuevos_abrigos",
			"language": {
				"code": "es",
			},
			"components" : [
				{
					"type": "body",
					"parameters": [
						{
							"type": "text",
							"text": "` + userName + `",
						},
						{
							"type": "text",
							"text": "` + forGender + `",
						},
					]
				}
			]
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	if resp.StatusCode != 200 {
		return err
	}

	return err
}

func SendResponseAnyMessage(codePhone, phone int, msg string) error {
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	codePhoneStr, phoneStr := strconv.Itoa(codePhone), strconv.Itoa(phone)

	to := codePhoneStr + phoneStr

	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "text",
		"text": {
			"body: "` + msg + `"
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return err

}

func SendImageByLink(codePhone, phone int, linkImg string) error {
	tokenMETA, _ := os.LookupEnv("META_BUSSINES_TOKEN")

	versionWP, _ := os.LookupEnv("WP_API_VERSION")
	phoneIdWP, _ := os.LookupEnv("WP_PHONE_ID")

	codePhoneStr, phoneStr := strconv.Itoa(codePhone), strconv.Itoa(phone)

	to := codePhoneStr + phoneStr

	jsonStr := []byte(`{
		"messaging_product": "whatsapp",
		"to": "` + to + `",
		"type": "image",
		"image": {
			"link: "` + linkImg + `"
		}
	}`)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v"+versionWP+"/"+phoneIdWP+"/messages", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenMETA)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	return err
}

// func SendFirstMessageWelcome(to, name string) error {
// 	// obtenemos el cliente de twilio
// 	client := clientTwilio()
// 	// esteblecemos los parametros
// 	msg := "Hola " + name + ".\nBienvenido a 'Acapela - Diseño y moda' la tienda de ropa que siempre te ayuda con vestirte bien.\nPuedes seguirnos en facebook: \nhttps://www.facebook.com/Acapela-Dise%C3%B1o-y-Moda-117114304308121 \nSi ya no quieres recibir notificaciones puedes darte de baja con un mensaje SALIR o inactivarte con INACTIVAR."
// 	params := paramsTwilioText(to, FROM, msg)
// 	// creando el mensaje
// 	resp, err := client.Api.CreateMessage(params)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Message Sid: " + *resp.Sid)
// 	return nil
// }
// func SendFirstNotification(name, to, types, gender string) error {
// 	// obtenemos el cliente de twilio
// 	client := clientTwilio()
// 	// esteblecemos los parametros
// 	msg := "Hola " + name + ".\nNotificarle que llego nuevos productos de " + types + " para " + gender + ". Puedes confirmar para que te enviemos fotos, pasar por la galeria o visitar nuestra pagina de facebook.\nSi ya no quieres recibir notificaciones puedes darte de baja con un mensaje SALIR o inactivarte por una semana con INACTIVAR."
// 	params := paramsTwilioText(to, FROM, msg)
// 	// creando el mensaje
// 	resp, err := client.Api.CreateMessage(params)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Message Sid: " + *resp.Sid)
// 	return nil
// }
// func clientTwilio() *twilio.RestClient {
// 	return twilio.NewRestClientWithParams(twilio.ClientParams{
// 		Username: "AC487e130568f8db39fc1687056035aacc",
// 		Password: "de1de83023e04012d86c9144f1f3a21c",
// 	})
// }
// func paramsTwilioText(to, from, msg string) (params *openapi.CreateMessageParams) {
// 	params = &openapi.CreateMessageParams{}
// 	params.SetTo("whatsapp:" + to)
// 	params.SetFrom("whatsapp:" + from)
// 	params.SetBody(msg)
// 	return
// }
