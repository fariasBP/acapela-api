package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Rol:

	1: SuperAdmin (Dueño, Administrado) (Todo)
	2: AdminEmploye (Empleado Mayor) (crea, actualiza y lee)
	3: Employe (Empleado Menor) (actualiza y lee)
	4: Client (Cliente) (actualiza y lee)
*/
type (
	/*
		ID: id del usuario (mongo establece este dato)
		Name: nombre del usuario {REQUERIDO}
		Lastname: apellido del usuario
		Photo: url de la foto del usuario (preferencia cloudinary)
		Email: correo electronico del usuario
		Code: codigo de ingreso del usuario (parecido a la contraseña que se genera automaticamente atravez de whatsapp o email)
		Phone: numero de telefono del usuario {REQUERIDO}
		Sleep: silenciar notificaciones
			
		NotReaded: Contador que indica cuantas veces el usuario no ha "recibido" sus mensajes recibidos
			- 0: indica que el usuario ha recibido mensajes (no necesariamente los ha leido)
			- >0 (mayor a cero): indica que el usuario no a recibido mensajes
		

	*/
	User struct {
		ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name               string             `json:"name" bson:"name,omitempty"`
		Lastname           string             `json:"lastname" bson:"lastname,omitempty"`
		Photo              string             `json:"photo" bson:"photo,omitempty"`
		Email              string             `json:"email" bson:"email,omitempty"`
		Code               string             `json:"code" bson:"code,omitempty"`
		Rol                int                `json:"rol" bson:"rol,omitempty"`
		Phone              int                `json:"phone" bson:"phone"`
		Sleep              uint8              `json:"sleep" bson:"sleep"`
		NotReaded          int                `json:"not_readed" bson:"not_readed,omitempty"`
		CodeDate           time.Time          `json:"code_date" bson:"code_date,omitempty"`
		SleepDate          time.Time          `json:"sleep_date" bson:"sleep_date,omitempty"`
		WpRegistrationDate time.Time          `json:"wp_registration_date" bson:"wp_registration_date,omitempty"`
		Mailbox            string             `json:"mailbox" bson:"mailbox,omitempty"`
		CountMailbox       int                `json:"count_mailbox" bson:"count_mailbox,omitempty"`
		Suscriptions       []Suscription      `json:"suscriptions" bson:"suscriptions,omitempty"`
		CreateDate         time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate         time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
	/* suscription
	IdShop: id de la tienda a la que se ha suscrito
	LikesKind: tipos de prenda de preferencia del usuario que quiere que se le notifique (ej. abrigo y sacos)
	LikesGender: preferencia de dama o varon
	LikesSize: preferencias de talla
	*/
	Suscription struct {
		IdShop      string `json:"id_shop" bson:"id_shop,omitempty"`
		LikesKind   string `json:"likes_kind" bson:"likes_kind,omitempty"`
		LikesGender int8   `json:"likes_gender" bson:"likes_gender,omitempty"`
		LikesSize   string `json:"likes_size" bson:"likes_size,omitempty"`
	}
)
