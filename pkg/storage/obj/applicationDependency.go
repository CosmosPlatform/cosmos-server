package obj

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type ApplicationDependency struct {
	CosmosObj
	ConsumerID int
	ProviderID int
	Consumer   *Application   `gorm:"foreignKey:ConsumerID"`
	Provider   *Application   `gorm:"foreignKey:ProviderID"`
	Reasons    pq.StringArray `gorm:"type:text[]"`
	Endpoints  Endpoints      `gorm:"type:jsonb"`
}

type PendingApplicationDependency struct {
	CosmosObj
	ConsumerID   int
	Consumer     *Application `gorm:"foreignKey:ConsumerID"`
	ProviderName string
	Reasons      pq.StringArray `gorm:"type:text[]"`
	Endpoints    Endpoints      `gorm:"type:jsonb"`
}

type Endpoints map[string]EndpointMethods

type EndpointMethods map[string]EndpointDetails

type EndpointDetails struct {
	Reasons []string `json:"reasons,omitempty"`
}

func (e Endpoints) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Endpoints) Scan(value any) error {
	if value == nil {
		*e = make(Endpoints)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Endpoints", value)
	}

	return json.Unmarshal(bytes, e)
}
