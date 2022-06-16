package domain

type EnvoyControlVersion struct {
	Id        uint64
	Version   uint32
	ServiceId string
}

func (EnvoyControlVersion) TableName() string {
	return "control_version"
}
