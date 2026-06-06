// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a GPLv3 license that can be found in the LICENSE file.

package route

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.design/x/ssaplayground/src/config"
)

const (
	// cleanInterval is how often the buildbox is swept for stale builds.
	cleanInterval = time.Hour
	// buildboxMaxAge is the maximum age a build folder is kept around.
	// Builds are addressed by their UUID path, so this is effectively the
	// lifetime of a shared link. A yearly horizon keeps the disk usage in
	// check while preserving recently created or bookmarked builds.
	// See https://github.com/golang-design/ssaplayground/issues/21
	buildboxMaxAge = 365 * 24 * time.Hour
)

// StartBuildboxCleanup periodically removes stale build folders so that the
// buildbox does not grow without bound and exhaust the disk.
func StartBuildboxCleanup() {
	dir := filepath.Join(config.Get().Static, "buildbox")
	tick := time.NewTicker(cleanInterval)
	defer tick.Stop()

	for ; ; <-tick.C {
		removed, err := cleanupBuildbox(dir, buildboxMaxAge, time.Now())
		if err != nil {
			log.Printf("buildbox cleanup error: %v", err)
		}
		if len(removed) > 0 {
			log.Printf("buildbox cleanup removed %d stale build(s): %v", len(removed), removed)
		}
	}
}

// cleanupBuildbox removes every build folder under dir whose last modification
// time is older than maxAge relative to now. It only touches directories whose
// name is a valid UUID so that bookkeeping files (e.g. index.html) are never
// deleted. It returns the names of the removed folders.
//
// The pure (dir, maxAge, now) signature keeps this testable without a clock.
func cleanupBuildbox(dir string, maxAge time.Duration, now time.Time) ([]string, error) {
	// Guard against an empty or root path so a misconfiguration can never
	// turn into a recursive delete of the filesystem.
	if dir == "" || dir == "/" || dir == "." {
		return nil, nil
	}
	if maxAge <= 0 {
		return nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		// A missing buildbox is not an error: it is created lazily on the
		// first build request.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var removed []string
	var firstErr error
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Only build folders are named after a UUID; skip anything else.
		if _, err := uuid.Parse(e.Name()); err != nil {
			continue
		}
		info, err := e.Info()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if now.Sub(info.ModTime()) <= maxAge {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		removed = append(removed, e.Name())
	}
	return removed, firstErr
}
