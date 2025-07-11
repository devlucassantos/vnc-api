package filters

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Event struct {
	StartDate               *time.Time
	EndDate                 *time.Time
	SituationId             *uuid.UUID
	LegislativeBodyId       *uuid.UUID
	RapporteurId            *uuid.UUID
	RemoveEventsInTheFuture *bool
}

func (instance Event) IsZero() bool {
	return reflect.DeepEqual(instance, Event{})
}
