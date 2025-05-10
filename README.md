# dtvm-go

A Go FFI wrapper for DTVM-ZetaEngine. This library allows you to run WebAssembly (WASM) modules from Go applications.

## Overview

dtvm-go provides a high-performance, safe interface to the DTVM-ZetaEngine WebAssembly runtime. It enables Go developers to execute WASM modules with full control over execution environment, memory, and resource usage.

## Features

- Load and execute WebAssembly modules
- Control execution with gas limits
- Isolate WASM execution environments
- Support for various runtime modes (interpreter, single-pass compilation, multi-pass optimization)
- Memory-safe bindings with proper resource cleanup

## Installation

```bash
go get github.com/IntelliXLabs/dtvm-go
```

## Prerequisites

This library requires the DTVM-ZetaEngine shared libraries to be available on your system. The required libraries include:
- libzetaengine
- libutils_lib
- libasmjit
- libspdlogd

## Usage Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/IntelliXLabs/dtvm-go"
)

func main() {
	// Create runtime with interpreter mode
	config := dtvm.NewRuntimeConfig(0)
	defer config.Delete()

	runtime := dtvm.NewRuntime(config)
	defer runtime.Delete()

	// Load WASM module from file
	wasmPath := "example/hello.wasm"
	module, err := runtime.LoadModuleFromFile(wasmPath)
	if err != nil {
		log.Fatalf("Failed to load module: %v", err)
	}
	defer module.Delete()

	// Create isolation environment
	isolation := runtime.CreateIsolation()
	defer isolation.Delete()

	// Create instance with gas limit
	instance, err := isolation.CreateInstanceWithGas(module, 10000000)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Delete()

	// Call WASM function
	args := []string{"world"}
	results, err := instance.CallWasmFuncByName(runtime, "hello", args)
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}

	// Process results
	if len(results) > 0 {
		fmt.Printf("Result: %v\n", results[0].Value)
	}
}
```

## Memory Management

dtvm-go uses Go finalizers to automatically clean up C resources when they're no longer needed. However, for deterministic resource cleanup, it