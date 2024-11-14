package sqlstorage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/db"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	conf := config.DatabaseConf{
		Username:          "calendar",
		Password:          "calendar-pwd",
		Host:              "localhost",
		Port:              5432,
		Database:          "calendar",
		SSLMode:           "disable",
		Timezone:          "Europe/Moscow",
		TargetSessionAttr: "read-write",
	}

	log := logger.New("debug", "stdout", "")

	pool, err := db.Connect(ctx, conf, log)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}

	return pool
}

func TestStorage_Add(t *testing.T) {
	t.Skip("Skipping TestStorage_Add")
	pool := setupTestDB(t)

	event := storage.Event{
		Id:          uuid.New(),
		Title:       "Test Event",
		Description: "This is a test event",
		StartDt:     time.Now(),
		EndDt:       time.Now().Add(1 * time.Hour),
		UserId:      uuid.Nil,
		Notify:      time.Minute * 10,
	}

	type fields struct {
		PgxPool *pgxpool.Pool
	}
	type args struct {
		ctx   context.Context
		event storage.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *storage.Event
		wantErr bool
	}{
		{
			name: "Add event",
			fields: fields{
				PgxPool: pool,
			},
			args: args{
				ctx:   context.Background(),
				event: event,
			},
			want:    &event,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				PgxPool: tt.fields.PgxPool,
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
