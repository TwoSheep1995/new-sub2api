package service

import (
	"context"
	"fmt"
)

// ListLatestUsableModels returns models whose latest enabled monitor check is operational.
func (s *ChannelMonitorService) ListLatestUsableModels(ctx context.Context) ([]LatestUsableMonitorModel, error) {
	monitors, err := s.repo.ListEnabled(ctx)
	if err != nil {
		return nil, fmt.Errorf("list enabled monitors: %w", err)
	}
	if len(monitors) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(monitors))
	for _, m := range monitors {
		ids = append(ids, m.ID)
	}

	latestMap, err := s.repo.ListLatestForMonitorIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list latest monitor results: %w", err)
	}

	out := make([]LatestUsableMonitorModel, 0)
	for _, m := range monitors {
		configured := configuredMonitorModelSet(m)
		for _, latest := range latestMap[m.ID] {
			if latest == nil || !isLatestMonitorStatusUsable(latest.Status) {
				continue
			}
			if _, ok := configured[latest.Model]; !ok {
				continue
			}
			out = append(out, LatestUsableMonitorModel{
				Provider: m.Provider,
				Model:    latest.Model,
			})
		}
	}
	return out, nil
}

func isLatestMonitorStatusUsable(status string) bool {
	return status == MonitorStatusOperational
}

func configuredMonitorModelSet(m *ChannelMonitor) map[string]struct{} {
	out := make(map[string]struct{})
	if m == nil {
		return out
	}
	if m.PrimaryModel != "" {
		out[m.PrimaryModel] = struct{}{}
	}
	for _, model := range m.ExtraModels {
		if model != "" {
			out[model] = struct{}{}
		}
	}
	return out
}
