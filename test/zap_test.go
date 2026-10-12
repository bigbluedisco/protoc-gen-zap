package test

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
	"google.golang.org/genproto/googleapis/type/timeofday"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bigbluedisco/protoc-gen-zap/test/output"
)

func TestOrderMarshalLogObject(t *testing.T) {
	order := &output.Order{
		ItemsOne: []*output.OrderItem{{Sku: "sku-1", Quantity: 2}},
		Of:       &output.Order_Foo{Foo: &output.Foo{S: "bar"}},
	}

	enc := zapcore.NewMapObjectEncoder()
	if err := order.MarshalLogObject(enc); err != nil {
		t.Fatalf("MarshalLogObject: %v", err)
	}

	items, ok := enc.Fields["items_one"].([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("items_one = %#v, want one object", enc.Fields["items_one"])
	}
	if got := items[0].(map[string]interface{})["sku"]; got != "sku-1" {
		t.Errorf("items_one[0].sku = %v, want sku-1", got)
	}
	if got := enc.Fields["of"].(map[string]interface{})["s"]; got != "bar" {
		t.Errorf("of.s = %v, want bar", got)
	}
}

func TestScheduleMarshalLogObject(t *testing.T) {
	cutoff := &timeofday.TimeOfDay{Hours: 14, Minutes: 30}
	slot := &timeofday.TimeOfDay{Hours: 9}
	start := &timeofday.TimeOfDay{Hours: 8, Minutes: 15}
	createdAt := timestamppb.New(time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC))

	schedule := &output.Schedule{
		Cutoff:         cutoff,
		Slots:          []*timeofday.TimeOfDay{slot},
		SlotsByCarrier: map[string]*timeofday.TimeOfDay{"ups": slot},
		Window:         &output.Schedule_WindowStart{WindowStart: start},
		CreatedAt:      createdAt,
		UpdatedAt:      []*timestamppb.Timestamp{createdAt},
	}

	enc := zapcore.NewMapObjectEncoder()
	if err := schedule.MarshalLogObject(enc); err != nil {
		t.Fatalf("MarshalLogObject: %v", err)
	}

	for key, want := range map[string]interface{}{
		"cutoff":     cutoff.String(),
		"window":     start.String(),
		"created_at": createdAt.String(),
		"slots":      []interface{}{slot.String()},
		"updated_at": []interface{}{createdAt.String()},
	} {
		if got := enc.Fields[key]; !equal(got, want) {
			t.Errorf("%s = %#v, want %#v", key, got, want)
		}
	}
	if _, ok := enc.Fields["slots_by_carrier"]; !ok {
		t.Errorf("slots_by_carrier missing")
	}

	empty := zapcore.NewMapObjectEncoder()
	if err := (&output.Schedule{}).MarshalLogObject(empty); err != nil {
		t.Fatalf("MarshalLogObject on empty schedule: %v", err)
	}
	if _, ok := empty.Fields["cutoff"]; ok {
		t.Errorf("cutoff should be omitted when unset")
	}
}

func equal(got, want interface{}) bool {
	gs, ok := got.([]interface{})
	if !ok {
		return got == want
	}
	ws := want.([]interface{})
	if len(gs) != len(ws) {
		return false
	}
	for i := range gs {
		if gs[i] != ws[i] {
			return false
		}
	}
	return true
}
