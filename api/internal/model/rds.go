package model

type RDSLogin struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	DdbInstanceIdentifier string `json:"dbInstanceIdentifier"`
}
