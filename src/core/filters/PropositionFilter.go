package filters

import (
	"github.com/google/uuid"
	"reflect"
)

type Proposition struct {
	DeputyId         *uuid.UUID
	PartyId          *uuid.UUID
	ExternalAuthorId *uuid.UUID
}

func (instance Proposition) IsZero() bool {
	return reflect.DeepEqual(instance, Proposition{})
}
