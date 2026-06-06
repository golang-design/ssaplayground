// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package route

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mkBuild creates a build folder named after a fresh UUID with the given
// modification time and returns its name.
func mkBuild(t *testing.T, dir string, modAge time.Duration, now time.Time) string {
	t.Helper()
	name := uuid.Must(uuid.NewUUID()).String()
	p := filepath.Join(dir, name)
	if err := os.Mkdir(p, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", p, err)
	}
	if err := os.WriteFile(filepath.Join(p, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	mtime := now.Add(-modAge)
	if err := os.Chtimes(p, mtime, mtime); err != nil {
		t.Fatalf("chtimes %s: %v", p, err)
	}
	return name
}

func TestCleanupBuildbox(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	maxAge := 30 * 24 * time.Hour

	// A non-UUID bookkeeping file must never be removed even if it is old.
	keep := filepath.Join(dir, "index.html")
	if err := os.WriteFile(keep, []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}
	old := now.Add(-365 * 24 * time.Hour)
	if err := os.Chtimes(keep, old, old); err != nil {
		t.Fatalf("chtimes index.html: %v", err)
	}

	stale1 := mkBuild(t, dir, 31*24*time.Hour, now)
	stale2 := mkBuild(t, dir, 90*24*time.Hour, now)
	fresh1 := mkBuild(t, dir, time.Hour, now)
	fresh2 := mkBuild(t, dir, 29*24*time.Hour, now)

	removed, err := cleanupBuildbox(dir, maxAge, now)
	if err != nil {
		t.Fatalf("cleanupBuildbox: %v", err)
	}

	sort.Strings(removed)
	want := []string{stale1, stale2}
	sort.Strings(want)
	if len(removed) != len(want) {
		t.Fatalf("removed = %v, want %v", removed, want)
	}
	for i := range want {
		if removed[i] != want[i] {
			t.Fatalf("removed = %v, want %v", removed, want)
		}
	}

	// Stale builds must be gone.
	for _, name := range []string{stale1, stale2} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Errorf("stale build %s still exists", name)
		}
	}
	// Fresh builds and the bookkeeping file must survive.
	for _, name := range []string{fresh1, fresh2, "index.html"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("expected %s to survive: %v", name, err)
		}
	}
}

func TestCleanupBuildboxGuards(t *testing.T) {
	now := time.Now()

	// Empty / root / dot paths must be refused outright.
	for _, dir := range []string{"", "/", "."} {
		if removed, err := cleanupBuildbox(dir, time.Hour, now); err != nil || removed != nil {
			t.Errorf("cleanupBuildbox(%q) = (%v, %v), want (nil, nil)", dir, removed, err)
		}
	}

	// A non-positive maxAge disables cleanup.
	dir := t.TempDir()
	mkBuild(t, dir, 365*24*time.Hour, now)
	if removed, err := cleanupBuildbox(dir, 0, now); err != nil || removed != nil {
		t.Errorf("cleanupBuildbox with maxAge=0 = (%v, %v), want (nil, nil)", removed, err)
	}

	// A missing buildbox is not an error.
	if removed, err := cleanupBuildbox(filepath.Join(dir, "does-not-exist"), time.Hour, now); err != nil || removed != nil {
		t.Errorf("cleanupBuildbox(missing) = (%v, %v), want (nil, nil)", removed, err)
	}
}
