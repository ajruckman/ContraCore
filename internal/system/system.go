// Package system stores common variables, loaded config values, and state info
// about ContraCore.
package system

import (
	"github.com/coredns/coredns/plugin/pkg/log"
	"go.uber.org/atomic"
)

var (
	// Whether any ContraCore CoreDNS zone has been initialized.
	// Used to prevent setup functions from running more than once
	// (for example multiple ContraWeb query log event servers).
	HasInitialized atomic.Bool

	// The common ContraCore instance of a CoreDNS logger.
	Console = log.NewWithPlugin("ContraCore")

	// Whether the system is currently known to be connected to ContraDB.
	ContraDBOnline atomic.Bool

	// Whether the system is currently known to be connected to ContraLog.
	ContraLogOnline atomic.Bool

	// Log time spent looking up whitelist and blacklist rules if true.
	LogRuleLookupDurations = true
)
