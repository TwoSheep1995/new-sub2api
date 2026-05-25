package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (r *channelRepository) ListModelSquareEntries(ctx context.Context) ([]service.ModelSquareEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, channel_id, group_id, platform, model_name, enabled, sort_order, created_at, updated_at
		FROM model_square_entries
		ORDER BY sort_order, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list model square entries: %w", err)
	}
	defer func() { _ = rows.Close() }()

	entries := make([]service.ModelSquareEntry, 0)
	for rows.Next() {
		var entry service.ModelSquareEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.ChannelID,
			&entry.GroupID,
			&entry.Platform,
			&entry.ModelName,
			&entry.Enabled,
			&entry.SortOrder,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan model square entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate model square entries: %w", err)
	}
	return entries, nil
}

func (r *channelRepository) ReplaceModelSquareEntries(ctx context.Context, entries []service.ModelSquareEntry) error {
	return r.runInTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM model_square_entries`); err != nil {
			return fmt.Errorf("clear model square entries: %w", err)
		}
		for i := range entries {
			entry := entries[i]
			if entry.SortOrder == 0 {
				entry.SortOrder = i + 1
			}
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO model_square_entries
					(channel_id, group_id, platform, model_name, enabled, sort_order)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, entry.ChannelID, entry.GroupID, entry.Platform, entry.ModelName, entry.Enabled, entry.SortOrder); err != nil {
				return fmt.Errorf("insert model square entry: %w", err)
			}
		}
		return nil
	})
}
