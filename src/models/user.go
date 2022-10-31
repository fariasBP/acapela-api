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

	Status:
		0: Normal (recibe notificaciones)
		1: Descanzo (no recibe notificaciones)
		2: Eliminado (ha eliminado su cuenta "supestamente")
*/
type (
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
		CreateDate         time.Time          `json:"create_date" bson:"create_date,omitempty"`
		UpdateDate         time.Time          `json:"update_date" bson:"update_date,omitempty"`
	}
)
