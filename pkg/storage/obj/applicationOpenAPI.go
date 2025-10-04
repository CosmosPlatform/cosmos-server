package obj

type ApplicationOpenAPI struct {
	CosmosObj
	ApplicationID int
	Application   *Application `gorm:"foreignKey:ApplicationID"`
	OpenAPI       string       `gorm:"type:jsonb"`
}
