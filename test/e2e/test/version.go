// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package test

import (
	"fmt"
	"testing"

	"github.com/elastic/cloud-on-k8s/pkg/controller/common/version"
)

// Elastic Stack versions used in the E2E tests
const (
	// MinVersion68x minimum version for 6.8.x tested with the operator
	MinVersion68x = "6.8.20"
	// LatestVersion7x current latest version for 7.x
	LatestVersion7x = "7.16.2" // version to synchronize with the latest release of the Elastic Stack
)

// SkipInvalidUpgrade skips a test that would do an invalid upgrade.
func SkipInvalidUpgrade(t *testing.T, srcVersion string, dstVersion string) {
	t.Helper()
	isValid, err := isValidUpgrade(srcVersion, dstVersion)
	if err != nil {
		t.Fatalf("Failed to determine the validity of the upgrade path: %v", err)
	}
	if !isValid {
		t.SkipNow()
	}
}

// isValidUpgrade reports whether an upgrade from one version to another version is valid.
func isValidUpgrade(from string, to string) (bool, error) {
	srcVer, err := version.Parse(from)
	if err != nil {
		return false, fmt.Errorf("failed to parse version '%s': %w", from, err)
	}
	dstVer, err := version.Parse(to)
	if srcVer.Pre != nil && dstVer.Pre == nil {
		// an upgrade from a pre-release version to a released version must not be tested (mainly due to incompatible licensing)
		// but an upgrade from a released version to a pre-release version is to be tested (to catch any new issues before the release)
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to parse version '%s': %w", to, err)
	}
	// major digits must be equal or differ by only 1
	validMajorDigit := dstVer.Major == srcVer.Major || dstVer.Major == srcVer.Major+1
	return validMajorDigit && !srcVer.GTE(dstVer), nil
}
