//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLatestMonitorStatusUsable_OnlyOperational(t *testing.T) {
	require.True(t, isLatestMonitorStatusUsable(MonitorStatusOperational))
	require.False(t, isLatestMonitorStatusUsable(MonitorStatusDegraded))
	require.False(t, isLatestMonitorStatusUsable(MonitorStatusFailed))
	require.False(t, isLatestMonitorStatusUsable(MonitorStatusError))
	require.False(t, isLatestMonitorStatusUsable(""))
}
