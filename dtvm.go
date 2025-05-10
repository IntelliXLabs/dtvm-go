package dtvm

/*
#cgo CXXFLAGS: -std=c++11 -fPIE
#cgo LDFLAGS: -L${SRCDIR}/lib -lzetaengine -lutils_lib -lasmjit -lspdlogd -lstdc++ -lm -lc -lgcc_s -lgcc

#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

// C struct definitions, corresponding to Rust structures
typedef struct {
    int32_t _mode;  // 0: interp, 1: singlepass, 2: multipass
} ZenRuntimeConfigExtern;

typedef struct {
    int32_t _dummy;
} ZenRuntimeExtern;

typedef struct {
    int32_t _dummy;
} ZenModuleExtern;

typedef struct {
    int32_t _dummy;
} ZenIsolationExtern;

typedef struct {
    int32_t _dummy;
} ZenInstanceExtern;

typedef struct {
    const char* name;
    uint32_t num_args;
    const uint32_t* arg_types;  // 0: i32, 1: i64, 2: f32, 3: f64
    uint32_t num_returns;
    const uint32_t* ret_types;  // 0: i32, 1: i64, 2: f32, 3: f64
    const void* ptr;            // c function pointer
} ZenHostFuncDescExtern;

typedef struct {
    int32_t _dummy;
} ZenHostModuleDescExtern;

typedef struct {
    int32_t _dummy;
} ZenHostModuleExtern;

typedef struct {
    int value_type;  // enum ZenType, 0: i32, 1: i64, 2: f32, 3: f64
    int64_t value;   // union of i32/i64/f32/f64
} ZenValueExtern;

// C function declarations
ZenRuntimeConfigExtern* ZenCreateRuntimeConfig(int32_t mode);
void ZenDeleteRuntimeConfig(ZenRuntimeConfigExtern* config);

ZenRuntimeExtern* ZenCreateRuntime(const ZenRuntimeConfigExtern* config);
void ZenDeleteRuntime(ZenRuntimeExtern* rt);

ZenHostModuleDescExtern* ZenCreateHostModuleDesc(
    ZenRuntimeExtern* rt,
    const char* host_mod_name,
    ZenHostFuncDescExtern* host_func_descs,
    uint32_t num_host_funcs
);
void ZenDeleteHostModuleDesc(
    ZenRuntimeExtern* rt,
    ZenHostModuleDescExtern* host_module_desc
);

ZenHostModuleExtern* ZenLoadHostModule(
    ZenRuntimeExtern* rt,
    ZenHostModuleDescExtern* host_module_desc
);
uint8_t ZenMergeHostModule(
    ZenRuntimeExtern* rt,
    ZenHostModuleExtern* host_module,
    ZenHostModuleDescExtern* other_host_module_desc
);
uint8_t ZenDeleteHostModule(
    ZenRuntimeExtern* rt,
    ZenHostModuleExtern* host_module
);
uint8_t ZenFilterHostFunctions(
    ZenHostModuleExtern* host_module,
    char** func_names,
    uint32_t num_func_names
);

ZenModuleExtern* ZenLoadModuleFromFile(
    ZenRuntimeExtern* rt,
    const char* wasm_path,
    char* error_buf,
    uint32_t error_buf_size
);

ZenModuleExtern* ZenLoadModuleFromBuffer(
    ZenRuntimeExtern* rt,
    const char* module_name,
    const uint8_t* code,
    uint32_t code_size,
    char* error_buf,
    uint32_t error_buf_size
);

uint32_t ZenGetNumImportFunctions(ZenModuleExtern* module);

uint32_t ZenGetImportFuncName(
    ZenModuleExtern* module,
    uint32_t func_idx,
    char* host_module_name_out,
    uint32_t* host_module_name_size_out,
    uint32_t host_module_name_out_buf_size,
    char* func_name_out,
    uint32_t* func_name_size_out,
    uint32_t func_name_out_buf_size
);

void ZenDeleteModule(ZenRuntimeExtern* rt, ZenModuleExtern* module);

ZenIsolationExtern* ZenCreateIsolation(ZenRuntimeExtern* rt);
void ZenDeleteIsolation(ZenRuntimeExtern* rt, ZenIsolationExtern* isolation);

ZenInstanceExtern* ZenCreateInstance(
    ZenIsolationExtern* isolation,
    ZenModuleExtern* wasm_mod,
    char* error_buf,
    uint32_t error_buf_size
);
ZenInstanceExtern* ZenCreateInstanceWithGas(
    ZenIsolationExtern* isolation,
    ZenModuleExtern* wasm_mod,
    uint64_t gas_limit,
    char* error_buf,
    uint32_t error_buf_size
);
void ZenDeleteInstance(ZenIsolationExtern* isolation, ZenInstanceExtern* inst);

int8_t ZenCallWasmFuncByName(
    ZenRuntimeExtern* rt,
    ZenInstanceExtern* inst,
    const char* func_name,
    const char** in_args,
    uint32_t num_in_args,
    ZenValueExtern* out_results,
    uint32_t* out_num_results
);

int8_t ZenGetInstanceError(
    ZenInstanceExtern* inst,
    char* error_buf,
    uint32_t error_buf_size
);

int8_t ZenValidateHostMemAddr(
    ZenInstanceExtern* inst,
    const void* host_addr,
    uint32_t size
);
int8_t ZenValidateAppMemAddr(
    ZenInstanceExtern* inst,
    uint32_t offset,
    uint32_t size
);
void* ZenGetHostMemAddr(
    ZenInstanceExtern* inst,
    uint32_t offset
);

uint32_t ZenGetAppMemOffset(
    ZenInstanceExtern* inst,
    const void* host_addr
);

uint64_t ZenGetInstanceGasLeft(ZenInstanceExtern* inst);
void ZenSetInstanceGasLeft(ZenInstanceExtern* inst, uint64_t new_gas);

void ZenSetInstanceCustomData(ZenInstanceExtern* inst, const void* custom_data);
const void* ZenGetInstanceCustomData(ZenInstanceExtern* inst);

void ZenSetInstanceExceptionByHostapi(
    ZenInstanceExtern* inst,
    uint32_t error_code
);

uint32_t ZenGetErrCodeEnvAbort();
uint32_t ZenGetErrCodeGasLimitExceeded();
uint32_t ZenGetErrCodeOutOfBoundsMemory();

void ZenInstanceExit(ZenInstanceExtern* inst, int32_t exit_code);

int8_t ZenGetExportFunc(
    ZenModuleExtern* wasm_mod,
    const char* func_name,
    uint32_t* func_idx_out
);

void ZenEnableLogging();
void ZenDisableLogging();

void ZenInstanceProtectMemoryAgain(ZenInstanceExtern* inst);
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Constants
const (
	// Runtime modes
	ModeInterp     int32 = 0
	ModeSinglepass int32 = 1
	ModeMultipass  int32 = 2

	// Value types
	TypeI32 int = 0
	TypeI64 int = 1
	TypeF32 int = 2
	TypeF64 int = 3

	// Error buffer size
	ErrorBufSize = 1024
)

// RuntimeConfig wraps ZenRuntimeConfigExtern
type RuntimeConfig struct {
	ptr *C.ZenRuntimeConfigExtern
}

// Runtime wraps ZenRuntimeExtern
type Runtime struct {
	ptr *C.ZenRuntimeExtern
}

// Module wraps ZenModuleExtern
type Module struct {
	ptr     *C.ZenModuleExtern
	runtime *Runtime // Keep reference to prevent GC
}

// Isolation wraps ZenIsolationExtern
type Isolation struct {
	ptr     *C.ZenIsolationExtern
	runtime *Runtime // Keep reference to prevent GC
}

// Instance wraps ZenInstanceExtern
type Instance struct {
	ptr       *C.ZenInstanceExtern
	isolation *Isolation // Keep reference to prevent GC
}

// HostFuncDesc wraps ZenHostFuncDescExtern
type HostFuncDesc struct {
	Name     string
	ArgTypes []int
	RetTypes []int
	Ptr      unsafe.Pointer
}

// HostModuleDesc wraps ZenHostModuleDescExtern
type HostModuleDesc struct {
	ptr     *C.ZenHostModuleDescExtern
	runtime *Runtime // Keep reference to prevent GC
}

// HostModule wraps ZenHostModuleExtern
type HostModule struct {
	ptr     *C.ZenHostModuleExtern
	runtime *Runtime // Keep reference to prevent GC
}

// Value wraps ZenValueExtern
type Value struct {
	Type  int
	Value interface{} // Can be int32, int64, float32, float64
}

// NewRuntimeConfig creates a runtime configuration
func NewRuntimeConfig(mode int32) *RuntimeConfig {
	ptr := C.ZenCreateRuntimeConfig(C.int32_t(mode))
	config := &RuntimeConfig{ptr: ptr}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(config, (*RuntimeConfig).Delete)
	return config
}

// Delete releases RuntimeConfig resources
func (c *RuntimeConfig) Delete() {
	if c.ptr != nil {
		C.ZenDeleteRuntimeConfig(c.ptr)
		c.ptr = nil
		runtime.SetFinalizer(c, nil)
	}
}

// NewRuntime creates a runtime
func NewRuntime(config *RuntimeConfig) *Runtime {
	ptr := C.ZenCreateRuntime(config.ptr)
	rt := &Runtime{ptr: ptr}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(rt, (*Runtime).Delete)
	return rt
}

// Delete releases Runtime resources
func (r *Runtime) Delete() {
	if r.ptr != nil {
		C.ZenDeleteRuntime(r.ptr)
		r.ptr = nil
		runtime.SetFinalizer(r, nil)
	}
}

// CreateHostModuleDesc creates a host module description
func (r *Runtime) CreateHostModuleDesc(hostModuleName string, hostFuncDescs []HostFuncDesc) (*HostModuleDesc, error) {
	cModName := C.CString(hostModuleName)
	defer C.free(unsafe.Pointer(cModName))

	numHostFuncs := len(hostFuncDescs)
	if numHostFuncs == 0 {
		return nil, errors.New("no host functions provided")
	}

	// Allocate C struct array
	cHostFuncDescs := (*C.ZenHostFuncDescExtern)(C.malloc(C.size_t(numHostFuncs) * C.size_t(unsafe.Sizeof(C.ZenHostFuncDescExtern{}))))
	defer C.free(unsafe.Pointer(cHostFuncDescs))

	// For keeping references to names and types
	var nameRefs []*C.char
	var argTypeRefs []*C.uint32_t
	var retTypeRefs []*C.uint32_t

	// Fill C struct array
	for i, hostFuncDesc := range hostFuncDescs {
		// Convert function name
		cName := C.CString(hostFuncDesc.Name)
		nameRefs = append(nameRefs, cName)

		// Convert argument types
		var cArgTypes *C.uint32_t
		if len(hostFuncDesc.ArgTypes) > 0 {
			cArgTypes = (*C.uint32_t)(C.malloc(C.size_t(len(hostFuncDesc.ArgTypes)) * C.size_t(unsafe.Sizeof(C.uint32_t(0)))))
			argTypeRefs = append(argTypeRefs, cArgTypes)
			for j, argType := range hostFuncDesc.ArgTypes {
				*(*C.uint32_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cArgTypes)) + uintptr(j)*unsafe.Sizeof(C.uint32_t(0)))) = C.uint32_t(argType)
			}
		}

		// Convert return types
		var cRetTypes *C.uint32_t
		if len(hostFuncDesc.RetTypes) > 0 {
			cRetTypes = (*C.uint32_t)(C.malloc(C.size_t(len(hostFuncDesc.RetTypes)) * C.size_t(unsafe.Sizeof(C.uint32_t(0)))))
			retTypeRefs = append(retTypeRefs, cRetTypes)
			for j, retType := range hostFuncDesc.RetTypes {
				*(*C.uint32_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cRetTypes)) + uintptr(j)*unsafe.Sizeof(C.uint32_t(0)))) = C.uint32_t(retType)
			}
		}

		// Set C struct fields
		desc := (*C.ZenHostFuncDescExtern)(unsafe.Pointer(uintptr(unsafe.Pointer(cHostFuncDescs)) + uintptr(i)*unsafe.Sizeof(C.ZenHostFuncDescExtern{})))
		desc.name = cName
		desc.num_args = C.uint32_t(len(hostFuncDesc.ArgTypes))
		desc.arg_types = cArgTypes
		desc.num_returns = C.uint32_t(len(hostFuncDesc.RetTypes))
		desc.ret_types = cRetTypes
		desc.ptr = hostFuncDesc.Ptr
	}

	// Call C function
	ptr := C.ZenCreateHostModuleDesc(r.ptr, cModName, cHostFuncDescs, C.uint32_t(numHostFuncs))

	// Free temporary resources
	for _, cName := range nameRefs {
		C.free(unsafe.Pointer(cName))
	}
	for _, cArgTypes := range argTypeRefs {
		C.free(unsafe.Pointer(cArgTypes))
	}
	for _, cRetTypes := range retTypeRefs {
		C.free(unsafe.Pointer(cRetTypes))
	}

	if ptr == nil {
		return nil, errors.New("failed to create host module desc")
	}

	hostModuleDesc := &HostModuleDesc{
		ptr:     ptr,
		runtime: r,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(hostModuleDesc, (*HostModuleDesc).Delete)
	return hostModuleDesc, nil
}

// Delete releases HostModuleDesc resources
func (h *HostModuleDesc) Delete() {
	if h.ptr != nil && h.runtime != nil && h.runtime.ptr != nil {
		C.ZenDeleteHostModuleDesc(h.runtime.ptr, h.ptr)
		h.ptr = nil
		runtime.SetFinalizer(h, nil)
	}
}

// LoadModuleFromFile loads a module from a WASM file
func (r *Runtime) LoadModuleFromFile(wasmPath string) (*Module, error) {
	cWasmPath := C.CString(wasmPath)
	defer C.free(unsafe.Pointer(cWasmPath))

	errorBuf := (*C.char)(C.malloc(C.size_t(ErrorBufSize)))
	defer C.free(unsafe.Pointer(errorBuf))

	modulePtr := C.ZenLoadModuleFromFile(r.ptr, cWasmPath, errorBuf, C.uint32_t(ErrorBufSize))
	if modulePtr == nil {
		errMsg := C.GoString(errorBuf)
		return nil, fmt.Errorf("failed to load module: %s", errMsg)
	}

	module := &Module{
		ptr:     modulePtr,
		runtime: r,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(module, (*Module).Delete)
	return module, nil
}

// LoadModuleFromBuffer loads a WASM module from memory
func (r *Runtime) LoadModuleFromBuffer(moduleName string, code []byte) (*Module, error) {
	cModuleName := C.CString(moduleName)
	defer C.free(unsafe.Pointer(cModuleName))

	errorBuf := (*C.char)(C.malloc(C.size_t(ErrorBufSize)))
	defer C.free(unsafe.Pointer(errorBuf))

	var codePtr *C.uint8_t
	if len(code) > 0 {
		codePtr = (*C.uint8_t)(unsafe.Pointer(&code[0]))
	}

	modulePtr := C.ZenLoadModuleFromBuffer(
		r.ptr,
		cModuleName,
		codePtr,
		C.uint32_t(len(code)),
		errorBuf,
		C.uint32_t(ErrorBufSize),
	)

	if modulePtr == nil {
		errMsg := C.GoString(errorBuf)
		return nil, fmt.Errorf("failed to load module: %s", errMsg)
	}

	module := &Module{
		ptr:     modulePtr,
		runtime: r,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(module, (*Module).Delete)
	return module, nil
}

// Delete releases Module resources
func (m *Module) Delete() {
	if m.ptr != nil && m.runtime != nil && m.runtime.ptr != nil {
		C.ZenDeleteModule(m.runtime.ptr, m.ptr)
		m.ptr = nil
		runtime.SetFinalizer(m, nil)
	}
}

// GetNumImportFunctions gets the number of imported functions
func (m *Module) GetNumImportFunctions() uint32 {
	return uint32(C.ZenGetNumImportFunctions(m.ptr))
}

// GetImportFuncName gets the import function name
func (m *Module) GetImportFuncName(funcIdx uint32) (hostModuleName string, funcName string, err error) {
	const bufSize = 255
	hostModuleNameBuf := (*C.char)(C.malloc(bufSize))
	defer C.free(unsafe.Pointer(hostModuleNameBuf))

	var hostModuleNameLen C.uint32_t

	funcNameBuf := (*C.char)(C.malloc(bufSize))
	defer C.free(unsafe.Pointer(funcNameBuf))

	var funcNameLen C.uint32_t

	result := C.ZenGetImportFuncName(
		m.ptr,
		C.uint32_t(funcIdx),
		hostModuleNameBuf,
		&hostModuleNameLen,
		C.uint32_t(bufSize),
		funcNameBuf,
		&funcNameLen,
		C.uint32_t(bufSize),
	)

	if result == 0 {
		return "", "", fmt.Errorf("failed to get import function name for index %d", funcIdx)
	}

	hostModuleName = C.GoStringN(hostModuleNameBuf, C.int(hostModuleNameLen))
	funcName = C.GoStringN(funcNameBuf, C.int(funcNameLen))

	return hostModuleName, funcName, nil
}

// CreateIsolation creates an isolation environment
func (r *Runtime) CreateIsolation() *Isolation {
	isolationPtr := C.ZenCreateIsolation(r.ptr)

	isolation := &Isolation{
		ptr:     isolationPtr,
		runtime: r,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(isolation, (*Isolation).Delete)
	return isolation
}

// Delete releases Isolation resources
func (i *Isolation) Delete() {
	if i.ptr != nil && i.runtime != nil && i.runtime.ptr != nil {
		C.ZenDeleteIsolation(i.runtime.ptr, i.ptr)
		i.ptr = nil
		runtime.SetFinalizer(i, nil)
	}
}

// CreateInstance creates an instance
func (i *Isolation) CreateInstance(mod *Module) (*Instance, error) {
	errorBuf := (*C.char)(C.malloc(C.size_t(ErrorBufSize)))
	defer C.free(unsafe.Pointer(errorBuf))

	instancePtr := C.ZenCreateInstance(i.ptr, mod.ptr, errorBuf, C.uint32_t(ErrorBufSize))
	if instancePtr == nil {
		errMsg := C.GoString(errorBuf)
		return nil, fmt.Errorf("failed to create instance: %s", errMsg)
	}

	instance := &Instance{
		ptr:       instancePtr,
		isolation: i,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(instance, (*Instance).Delete)
	return instance, nil
}

// CreateInstanceWithGas creates an instance with gas limit
func (i *Isolation) CreateInstanceWithGas(mod *Module, gasLimit uint64) (*Instance, error) {
	errorBuf := (*C.char)(C.malloc(C.size_t(ErrorBufSize)))
	defer C.free(unsafe.Pointer(errorBuf))

	instancePtr := C.ZenCreateInstanceWithGas(i.ptr, mod.ptr, C.uint64_t(gasLimit), errorBuf, C.uint32_t(ErrorBufSize))
	if instancePtr == nil {
		errMsg := C.GoString(errorBuf)
		return nil, fmt.Errorf("failed to create instance with gas: %s", errMsg)
	}

	instance := &Instance{
		ptr:       instancePtr,
		isolation: i,
	}

	// Use finalizer to ensure resources are released
	runtime.SetFinalizer(instance, (*Instance).Delete)
	return instance, nil
}

// Delete releases Instance resources
func (inst *Instance) Delete() {
	if inst.ptr != nil && inst.isolation != nil && inst.isolation.ptr != nil {
		C.ZenDeleteInstance(inst.isolation.ptr, inst.ptr)
		inst.ptr = nil
		runtime.SetFinalizer(inst, nil)
	}
}

// CallWasmFuncByName calls a WASM function
func (inst *Instance) CallWasmFuncByName(runtime *Runtime, funcName string, args []string) ([]Value, error) {
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	// Convert arguments
	var cArgs []*C.char
	var cArgsPtr **C.char

	if len(args) > 0 {
		cArgs = make([]*C.char, len(args))
		for i, arg := range args {
			cArgs[i] = C.CString(arg)
			defer C.free(unsafe.Pointer(cArgs[i]))
		}
		cArgsPtr = (**C.char)(unsafe.Pointer(&cArgs[0]))
	}

	// Prepare for return values
	var outResults [10]C.ZenValueExtern // Assume max 10 return values
	var outNumResults C.uint32_t

	result := C.ZenCallWasmFuncByName(
		runtime.ptr,
		inst.ptr,
		cFuncName,
		(**C.char)(unsafe.Pointer(cArgsPtr)),
		C.uint32_t(len(args)),
		&outResults[0],
		&outNumResults,
	)

	if result == 0 {
		// Get error information
		errorBuf := (*C.char)(C.malloc(C.size_t(ErrorBufSize)))
		defer C.free(unsafe.Pointer(errorBuf))

		C.ZenGetInstanceError(inst.ptr, errorBuf, C.uint32_t(ErrorBufSize))
		errMsg := C.GoString(errorBuf)

		return nil, fmt.Errorf("failed to call WASM function: %s", errMsg)
	}

	// Process return values
	numResults := int(outNumResults)
	values := make([]Value, numResults)

	for i := 0; i < numResults; i++ {
		valueType := int(outResults[i].value_type)
		var value interface{}

		switch valueType {
		case TypeI32:
			value = int32(outResults[i].value)
		case TypeI64:
			value = int64(outResults[i].value)
		case TypeF32:
			// Handle floating point numbers, requires special bit representation handling
			value = *(*float32)(unsafe.Pointer(&outResults[i].value))
		case TypeF64:
			// Handle floating point numbers, requires special bit representation handling
			value = *(*float64)(unsafe.Pointer(&outResults[i].value))
		}

		values[i] = Value{
			Type:  valueType,
			Value: value,
		}
	}

	return values, nil
}

// GetGasLeft gets the remaining gas
func (inst *Instance) GetGasLeft() uint64 {
	return uint64(C.ZenGetInstanceGasLeft(inst.ptr))
}

// SetGasLeft sets the remaining gas
func (inst *Instance) SetGasLeft(gas uint64) {
	C.ZenSetInstanceGasLeft(inst.ptr, C.uint64_t(gas))
}

// EnableLogging enables logging
func EnableLogging() {
	C.ZenEnableLogging()
}

// DisableLogging disables logging
func DisableLogging() {
	C.ZenDisableLogging()
}
