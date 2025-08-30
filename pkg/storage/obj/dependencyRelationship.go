package obj

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type DependencyRelationship struct {
	CosmosObj
	ConsumerID int
	ProviderID int
	Consumer   *Application   `gorm:"foreignKey:ConsumerID"`
	Provider   *Application   `gorm:"foreignKey:ProviderID"`
	Reasons    pq.StringArray `gorm:"type:text[]"`
	Endpoints  Endpoints      `gorm:"type:jsonb"`
}

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path    string   `json:"path"`
	Method  string   `json:"method"`
	Reasons []string `json:"reasons"`
}

func (e Endpoints) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Endpoints) Scan(value any) error {
	if value == nil {
		*e = Endpoints{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Endpoints", value)
	}

	return json.Unmarshal(bytes, e)
}
