package models

import (
	"time"

	"github.com/google/uuid"
)

type Equipment struct {
	BaseModel
	Name            string
	Model           string
	ServicingPeriod time.Month
	HospitalID      uuid.UUID
	Hospital        Hospital
}
