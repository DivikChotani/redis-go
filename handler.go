// Package main implements a minimal Redis-like server with basic in-memory command support.
// This file specifically defines the command handlers and thread-safe data stores.

package main

import (
	"sync"
)

// -----------------------------------------------------------------------------
// Command Handler Registry
// -----------------------------------------------------------------------------

// Handlers maps Redis command names to their corresponding handler functions.
// Each handler receives a slice of Values as arguments and returns a Value as the result.
var Handlers = map[string]func([]Value) Value{
	"PING":    handlePing,    // Health check
	"SET":     handleSet,     // Set a string key
	"GET":     handleGet,     // Get a string key
	"HSET":    handleHSet,    // Set a field in a hash
	"HGET":    handleHGet,    // Get a field from a hash
	"HGETALL": handleHGetAll, // Get all fields/values in a hash
	"DEL":     handleDel,     // Delete one or more keys
}

// -----------------------------------------------------------------------------
// String Key-Value Store (SET/GET)
// -----------------------------------------------------------------------------

var (
	stringStore   = map[string]string{} // In-memory key-value store for strings
	stringStoreMu = sync.RWMutex{}      // Read-write mutex for stringStore
)

// handleSet implements the SET command.
// Syntax: SET key value
func handleSet(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'SET' command"}
	}

	key := args[0].bulk
	val := args[1].bulk

	stringStoreMu.Lock()
	stringStore[key] = val
	stringStoreMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// handleGet implements the GET command.
// Syntax: GET key
func handleGet(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'GET' command"}
	}

	key := args[0].bulk

	stringStoreMu.RLock()
	val, exists := stringStore[key]
	stringStoreMu.RUnlock()

	if !exists {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}
