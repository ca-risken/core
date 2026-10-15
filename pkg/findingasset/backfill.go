package findingasset

import (
	"context"
	"database/sql"
)

// BackfillChunk locks the latest rows on the primary. Scan writes wait for this
// transaction; whichever write happens last classifies its own JSON. Only asset
// columns change, including an explicit updated_at self-assignment.
func BackfillChunk(ctx context.Context, db *sql.DB, after uint64, limit int) (uint64, int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return after, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, "SELECT finding_id, data_source, resource_name, COALESCE(data,'{}') FROM finding WHERE finding_id > ? AND data_source = 'google:scc' ORDER BY finding_id LIMIT ? FOR UPDATE", after, limit)
	if err != nil {
		return after, 0, err
	}
	type entry struct {
		id                     uint64
		source, resource, data string
	}
	var entries []entry
	for rows.Next() {
		var e entry
		if err = rows.Scan(&e.id, &e.source, &e.resource, &e.data); err != nil {
			rows.Close()
			return after, 0, err
		}
		entries = append(entries, e)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return after, 0, err
	}
	for _, e := range entries {
		a := Classify(e.source, e.resource, e.data)
		_, err = tx.ExecContext(ctx, "UPDATE finding SET asset_type=NULLIF(?,''), asset_key=NULLIF(?,''), asset_name=NULLIF(?,''), updated_at=updated_at WHERE finding_id=?", a.Type, a.Key, a.Name, e.id)
		if err != nil {
			return after, 0, err
		}
		after = e.id
	}
	if err = tx.Commit(); err != nil {
		return after, 0, err
	}
	return after, len(entries), nil
}
