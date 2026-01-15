package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNotificationService_Creation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	service := NewNotificationService(logger)
	assert.NotNil(t, service)
}
