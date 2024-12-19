//go:build !ci

package sqlstorage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/database"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	conf := database.Config{
		Username:          "calendar",
		Password:          "calendar-pwd",
		Host:              "localhost",
		Port:              5432,
		Database:          "calendar",
		SSLMode:           "disable",
		Timezone:          "Europe/Moscow",
		TargetSessionAttr: "read-write",
	}

	pool, err := database.Connect(ctx, conf)
	if err != nil {
		t.Fatalf("Unable to Connect to database: %v", err)
	}

	return pool
}

func TestStorage_Add(t *testing.T) {
	pool := setupTestDB(t)

	e := model.Event{
		ID:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartDt:     time.Now(),
		EndDt:       time.Now().Add(1 * time.Hour),
		UserID:      uuid.Nil,
		NotifyAt:    time.Minute * 10,
	}

	type fields struct {
		PgxPool *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		event model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Event
		wantErr bool
	}{
		{
			name: "Add event",
			fields: fields{
				PgxPool: pool,
			},
			args: args{
				ctx:   context.Background(),
				event: e,
			},
			want:    e,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				conn: tt.fields.PgxPool,
			}
			got, err := s.Add(tt.args.ctx, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}
