package dao

import (
	bizConfig "isc-envoy-control-service/config"
	"isc-envoy-control-service/pojo/domain"
)

func GetServiceVersion(serviceId string) uint32 {
	version := domain.EnvoyControlVersion{ServiceId: serviceId}
	bizConfig.Db.Find(&version, "service_id=?", serviceId)
	return version.Version
}

func UpdateServiceVersion(serviceId string, version uint32) {
	result := bizConfig.Db.Find(&domain.EnvoyControlVersion{}, "service_id", serviceId)
	if result.RowsAffected == 0 {
		bizConfig.Db.Create(&domain.EnvoyControlVersion{
			Version:   version,
			ServiceId: serviceId,
		})
	} else {
		bizConfig.Db.Model(&domain.EnvoyControlVersion{}).Where("service_id", serviceId).Update("version", version)
	}
}
