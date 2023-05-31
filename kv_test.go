package cache_bench

import (
	"fmt"
	"testing"

	"github.com/c-pro/geche"
	"github.com/google/go-cmp/cmp"
)

func ExampleNewKV() {
	kv := NewKVCacheMap[string]()

	kv.Set("foo", "bar")
	kv.Set("foo2", "bar2")
	kv.Set("foo3", "bar3")
	kv.Set("foo1", "bar1")

	res, _ := kv.ListByPrefix("foo")
	fmt.Println(res)
	// Output: [bar bar1 bar2 bar3]
}

func compareSlice(t *testing.T, exp, got []string) {
	t.Helper()

	t.Log(got)
	if len(exp) != len(got) {
		t.Fatalf("expected length %d, got %d", len(exp), len(got))
	}

	for i := 0; i < len(exp); i++ {
		if exp[i] != got[i] {
			t.Errorf("expected %q, got %q", exp[i], got[i])
		}
	}
}

func TestKV(t *testing.T) {
	kv := NewKVCacheMap[string]()

	for i := 999; i >= 0; i-- {
		key := fmt.Sprintf("%03d", i)
		kv.Set(key, key)
	}

	if kv.Len() != 1000 {
		t.Fatalf("expected length %d, got %d", 1000, kv.Len())
	}

	expected := []string{
		"000", "001", "002", "003", "004", "005", "006", "007", "008", "009",
	}

	got, err := kv.ListByPrefix("00")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}
	compareSlice(t, expected, got)

	expected = []string{
		"120", "121", "122", "123", "124", "125", "126", "127", "128", "129",
	}

	got, err = kv.ListByPrefix("12")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}
	compareSlice(t, expected, got)

	expected = []string{"888"}

	got, err = kv.ListByPrefix("888")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}
	compareSlice(t, expected, got)

	_ = kv.Del("777")
	_ = kv.Del("779")

	if kv.Len() != 998 {
		t.Fatalf("expected length %d, got %d", 998, kv.Len())
	}

	if _, err := kv.Get("777"); err != geche.ErrNotFound {
		t.Fatalf("expected error %v, got %v", geche.ErrNotFound, err)
	}

	expected = []string{
		"770", "771", "772", "773", "774", "775", "776", "778",
	}

	got, err = kv.ListByPrefix("77")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}

	compareSlice(t, expected, got)

	kv.Set("777", "777")
	kv.Set("779", "779")

	if kv.Len() != 1000 {
		t.Fatalf("expected length %d, got %d", 1000, kv.Len())
	}

	expected = []string{
		"770", "771", "772", "773", "774", "775", "776", "777", "778", "779",
	}

	got, err = kv.ListByPrefix("77")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}

	compareSlice(t, expected, got)

	kv.Set("77", "77")

	if kv.Len() != 1001 {
		t.Fatalf("expected length %d, got %d", 1001, kv.Len())
	}

	expected = []string{
		"77", "770", "771", "772", "773", "774", "775", "776", "777", "778", "779",
	}

	got, err = kv.ListByPrefix("77")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}

	compareSlice(t, expected, got)
}

func TestKVEmptyPrefix(t *testing.T) {
	kv := NewKVCacheMap[string]()

	expected := []string{}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%02d", i)
		expected = append(expected, key)
		kv.Set(key, key)
	}

	got, err := kv.ListByPrefix("")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}

	compareSlice(t, expected, got)
}

func TestKVNonexist(t *testing.T) {
	kv := NewKVCacheMap[string]()

	kv.Set("test", "best")

	got, err := kv.ListByPrefix("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error in ListByPrefix: %v", err)
	}

	if len(got) > 0 {
		t.Errorf("unexpected len %d", len(got))
	}
}

func TestKVSnapshot(t *testing.T) {
	t.Skip()
	kv := NewKVCacheMap[string]()

	kv.Set("test", "best")
	kv.Set("test2", "best2")

	snap := kv.Snapshot()

	kv.Set("test3", "best3")

	if kv.Len() != 3 {
		t.Fatalf("expected length %d, got %d", 3, kv.Len())
	}

	if len(snap) != 2 {
		t.Fatalf("expected length %d, got %d", 2, len(snap))
	}

	expected := map[string]string{"test": "best", "test2": "best2"}

	if diff := cmp.Diff(expected, snap); diff != "" {
		t.Errorf("unexpected diff: %s", diff)
	}

	snap = kv.Snapshot()
	expected = map[string]string{"test": "best", "test2": "best2", "test3": "best3"}

	if diff := cmp.Diff(expected, snap); diff != "" {
		t.Errorf("unexpected diff: %s", diff)
	}
}
