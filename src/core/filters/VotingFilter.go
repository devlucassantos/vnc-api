package filters

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Voting struct {
	StartDate         *time.Time
	EndDate           *time.Time
	Result            string
	LegislativeBodyId *uuid.UUID
}

func (instance Voting) IsZero() bool {
	return reflect.DeepEqual(instance, Voting{})
}
