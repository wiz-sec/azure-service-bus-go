package servicebus

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBatchDispositionIterator(t *testing.T) {
	count := 20
	fetched := 0
	lockIDs := []*uuid.UUID{}

	for i := count; i > 0; i-- {
		lockIDs = append(lockIDs, &uuid.UUID{})
	}

	bdi := &BatchDispositionIterator{
		LockTokenIDs: lockIDs,
	}

	assert.Equal(t, 0, bdi.cursor)

	for !bdi.Done() {
		if uuid := bdi.Next(); uuid != nil {
			fetched++
		}
	}
	assert.Equal(t, count, fetched)
}

func TestBatchDispositionUnsupportedStatus(t *testing.T) {
	status := MessageStatus(suspendedDisposition)
	id := uuid.UUID{}
	bdi := BatchDispositionIterator{
		LockTokenIDs: []*uuid.UUID{
			&id, &id, &id,
		},
		Status: status,
	}

	subscription := Subscription{}
	err := subscription.SendBatchDisposition(context.Background(), bdi)
	be := err.(BatchDispositionError)
	assert.NotNil(t, be, fmt.Sprintf("Wrong error type %T", err))
	assert.EqualErrorf(t, err, fmt.Sprintf("Operation failed, %d error(s) reported.", len(be.Errors)), err.Error())

	for _, innerErr := range be.Errors {
		assert.NotNil(t, innerErr.UnWrap(), "Unwrapped error is nil")
		assert.EqualErrorf(t, innerErr, "unsupported bulk disposition status \"suspended\"", innerErr.Error())
	}
}
