package cms

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestKeyValueSetGet(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{"string", "hello", "hello"},
		{"int", 42, int64(42)},
		{"int64", int64(-7), int64(-7)},
		{"bool true", true, true},
		{"bool false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv, err := NewKeyValue("k", tt.input)
			if err != nil {
				t.Fatalf("NewKeyValue(%v): %v", tt.input, err)
			}

			got, err := kv.Get()
			if err != nil {
				t.Fatalf("Get(): %v", err)
			}
			if got != tt.want {
				t.Fatalf("Get() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestKeyValueSetUnsupported(t *testing.T) {
	kv := &KeyValue{}
	if err := kv.Set(3.14); err == nil {
		t.Fatal("Set(3.14) = nil error, want unsupported type error")
	}
}

func TestKeyValueGetUnknownType(t *testing.T) {
	kv := &KeyValue{Type: "float"}
	if _, err := kv.Get(); err == nil {
		t.Fatal("Get() = nil error, want unknown type error")
	}
}

func TestKeyValueRoundTrip(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(&KeyValue{}).Error; err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	kv, err := NewKeyValue("answer", 42)
	if err != nil {
		t.Fatalf("NewKeyValue: %v", err)
	}

	if err := db.Create(kv).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	var got KeyValue
	if err := db.Where("key = ?", "answer").First(&got).Error; err != nil {
		t.Fatalf("find: %v", err)
	}

	if got.Type != ValueInt {
		t.Fatalf("Type = %q, want %q", got.Type, ValueInt)
	}
	n, err := got.Int()
	if err != nil {
		t.Fatalf("Int(): %v", err)
	}
	if n != 42 {
		t.Fatalf("Int() = %d, want 42", n)
	}
}
