package models

type ProductPattern struct {
	Talla           string `json:"talla" bson:"talla,omitempty"`
	LargoTorso      uint8  `json:"largotorso" bson:"largotorso,omitempty"`
	ContornoPecho   uint8  `json:"contornopecho" bson:"contornopecho,omitempty"`
	ContornoCintura uint8  `json:"contornocintura" bson:"contornocintura,omitempty"`
	ContornoCadera  uint8  `json:"contornocadera" bson:"contornocadera,omitempty"`
	ContornoSisa    uint8  `json:"contornosisa" bson:"contornosisa,omitempty"`
	LargoHombro     uint8  `json:"largohombro" bson:"largohombro,omitempty"`
	LargoManga      uint8  `json:"largomanga" bson:"largomanga,omitempty"`
}
