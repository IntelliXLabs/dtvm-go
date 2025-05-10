package dtvm

import (
	"fmt"
	"testing"
)

func TestRuntimeCreation(t *testing.T) {
	config := NewRuntimeConfig(0) // 0 means interpreter mode
	defer config.Delete()
}

func TestWasmFib(t *testing.T) {
	// Create runtime config and runtime
	config := NewRuntimeConfig(0) // Interpreter mode
	defer config.Delete()

	runtime := NewRuntime(config)
	defer runtime.Delete()

	// Load WASM module
	wasmPath := "example/fib.0.wasm"
	wasmMod, err := runtime.LoadModuleFromFile(wasmPath)
	if err != nil {
		t.Fatalf("Failed to load WASM module: %v", err)
		return
	}
	defer wasmMod.Delete()

	// Create isolation
	isolation := runtime.CreateIsolation()
	defer isolation.Delete()

	// Create WASM instance
	gasLimit := uint64(100000000)
	instance, err := isolation.CreateInstanceWithGas(wasmMod, gasLimit)
	if err != nil {
		t.Fatalf("Failed to create WASM instance: %v", err)
		return
	}
	defer instance.Delete()

	// Call WASM function
	args := []string{"5"} // Pass argument 5
	results, err := instance.CallWasmFuncByName(runtime, "fib", args)
	if err != nil {
		t.Fatalf("Failed to call WASM function: %v", err)
		return
	}

	// Verify we have results
	if len(results) == 0 {
		t.Fatal("No results returned")
	}

	// Verify result type is valid
	result := results[0]
	switch result.Value.(type) {
	case int32, int64:
		// Valid types
	default:
		t.Fatalf("Result is not int32 or int64, but %T", result.Value)
	}

	// Call the function again to verify consistency
	repeatResults, err := instance.CallWasmFuncByName(runtime, "fib", args)
	if err != nil {
		t.Fatalf("Failed to call WASM function again: %v", err)
	}

	if len(repeatResults) == 0 {
		t.Fatal("No results returned on second call")
	}

	// Verify the result is the same on second call
	if fmt.Sprintf("%v", result.Value) != fmt.Sprintf("%v", repeatResults[0].Value) {
		t.Errorf("Inconsistent results: first=%v, second=%v",
			result.Value, repeatResults[0].Value)
	}
}
