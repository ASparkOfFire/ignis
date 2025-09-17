package runtime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ASparkOfFire/ignis/internal/cache"
	"github.com/ASparkOfFire/ignis/internal/runtime/js"

	"github.com/google/uuid"
	"github.com/ignis-runtime/wasi-go/imports"
	"github.com/ignis-runtime/wasi-go/imports/wasi_http"
	"github.com/ignis-runtime/wazero"
)

//go:generate stringer --type RuntimeEngine
type RuntimeEngine int

const (
	RuntimeEngineWASM RuntimeEngine = iota
	RuntimeEngineJS
)

// NetworkConfig defines network-related configuration
type NetworkConfig struct {
	Listens []string // Network addresses to listen on
	Dials   []string // Network addresses allowed for dialing
}

// WasiConfig defines WASI-specific configuration
type WasiConfig struct {
	Dirs         []string // Directories to mount
	MaxOpenFiles int      // Maximum open files limit
	EnableHttp   bool     // Enable HTTP support
}

// Args defines the configuration for creating a new Runtime instance.
type Args struct {
	Stdout       io.Writer
	DeploymentID uuid.UUID
	Engine       RuntimeEngine
	Blob         []byte
	Cache        cache.ModCache[uuid.UUID]
	Network      *NetworkConfig // Optional network configuration
	Wasi         *WasiConfig    // Optional WASI configuration
}

// Runtime manages the WebAssembly execution environment.
type Runtime struct {
	stdout       io.Writer
	ctx          context.Context
	deploymentID uuid.UUID
	engine       RuntimeEngine
	mod          wazero.CompiledModule
	runtime      wazero.Runtime
	cache        cache.ModCache[uuid.UUID]
	network      *NetworkConfig
	wasi         *WasiConfig
	wasiHTTP     *wasi_http.WasiHTTP
}

// defaultNetworkConfig returns default network configuration
func defaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Listens: []string{},
		Dials:   []string{},
	}
}

// defaultWasiConfig returns default WASI configuration
func defaultWasiConfig() *WasiConfig {
	return &WasiConfig{
		Dirs:         []string{},
		MaxOpenFiles: 1024,
		EnableHttp:   false,
	}
}

// New initializes a new WebAssembly runtime with the given arguments.
func New(ctx context.Context, args Args) (*Runtime, error) {
	// Set defaults for optional configurations
	network := args.Network
	if network == nil {
		network = defaultNetworkConfig()
	}

	wasiConfig := args.Wasi
	if wasiConfig == nil {
		wasiConfig = defaultWasiConfig()
	}

	if !args.Cache.Has(args.DeploymentID) {
		args.Cache.Add(args.DeploymentID, wazero.NewCompilationCache())
	}

	config := wazero.NewRuntimeConfigCompiler().WithCompilationCache(args.Cache.Get(args.DeploymentID))
	rt := wazero.NewRuntimeWithConfig(ctx, config)

	blob := args.Blob
	switch args.Engine {
	case RuntimeEngineWASM:
		// For WASM, we'll set up enhanced WASI after compiling the module
	case RuntimeEngineJS:
		blob = js.Runtime
	default:
		rt.Close(ctx)
		return nil, fmt.Errorf("unsupported runtime engine %q", args.Engine)
	}

	mod, err := rt.CompileModule(ctx, blob)
	if err != nil {
		rt.Close(ctx) // Cleanup on failure
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	runtime := &Runtime{
		runtime:      rt,
		ctx:          ctx,
		deploymentID: args.DeploymentID,
		engine:       args.Engine,
		stdout:       args.Stdout,
		mod:          mod,
		network:      network,
		wasi:         wasiConfig,
	}

	// Set up enhanced WASI for WASM modules (this must happen after module compilation)
	if args.Engine == RuntimeEngineWASM {
		if err := runtime.setupEnhancedWASI(); err != nil {
			runtime.Close()
			return nil, fmt.Errorf("failed to setup enhanced WASI: %w", err)
		}
	}

	return runtime, nil
}

// setupEnhancedWASI configures the enhanced WASI environment exactly like the reference
func (r *Runtime) setupEnhancedWASI() error {
	// Match the reference implementation exactly
	wasmName := fmt.Sprintf("deployment-%s", r.deploymentID.String())

	builder := imports.NewBuilder().
		WithName(wasmName).
		WithSocketsExtension("auto", r.mod).
		WithCustomDNSServer("8.8.8.8")

	// instantiate without stdio here; just setup context and system for WASI HTTP if enabled
	ctx, system, err := builder.Instantiate(r.ctx, r.runtime)
	if err != nil {
		return fmt.Errorf("failed to instantiate enhanced WASI: %w", err)
	}

	// Update the context
	r.ctx = ctx

	// Setup cleanup
	go func() {
		<-ctx.Done()
		system.Close(ctx)
	}()

	// Setup HTTP like the reference
	importWasi := false
	switch r.wasi.EnableHttp {
	case true:
		importWasi = wasi_http.DetectWasiHttp(r.mod)
	default:
		importWasi = false
	}

	if importWasi {
		r.wasiHTTP = wasi_http.MakeWasiHTTP()
		if err := r.wasiHTTP.Instantiate(ctx, r.runtime); err != nil {
			return fmt.Errorf("failed to setup WASI HTTP: %w", err)
		}
	}

	return nil
}

// Invoke executes the compiled WebAssembly module with provided input and environment variables.
func (r *Runtime) Invoke(stdin io.Reader, env map[string]string, script []byte, args ...string) error {
	defer r.Close()

	switch r.engine {
	case RuntimeEngineWASM:
		return r.invokeWASM(stdin, env, args...)
	case RuntimeEngineJS:
		return r.invokeJS(stdin, env, script, args...)
	default:
		return fmt.Errorf("invalid runtime engine %d", r.engine)
	}
}

// invokeWASM handles WASM module execution - preserving stdin/stdout args,
// but uses pipes under the hood to connect to WASI stdio.
func (r *Runtime) invokeWASM(stdin io.Reader, env map[string]string, args ...string) error {
	// Create OS pipes for stdin
	rStdin, wStdin, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Create OS pipes for stdout
	rStdout, wStdout, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Asynchronously copy from provided stdin to pipe writer (wStdin)
	go func() {
		defer wStdin.Close()
		if _, err := io.Copy(wStdin, stdin); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed copying to WASI stdin: %v\n", err)
		}
	}()

	// Asynchronously copy from pipe reader (rStdout) to provided stdout
	go func() {
		defer rStdout.Close()
		if _, err := io.Copy(r.stdout, rStdout); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed copying from WASI stdout: %v\n", err)
		}
	}()

	// Setup WASI with pipe file descriptors for stdio
	builder := imports.NewBuilder().
		WithName(fmt.Sprintf("deployment-%s", r.deploymentID.String())).
		WithSocketsExtension("auto", r.mod).
		WithCustomDNSServer("1.1.1.1").
		WithStdio(int(rStdin.Fd()), int(wStdout.Fd()), int(os.Stderr.Fd()))

	ctx, system, err := builder.Instantiate(r.ctx, r.runtime)
	if err != nil {
		return fmt.Errorf("failed to instantiate enhanced WASI: %w", err)
	}
	defer system.Close(ctx)

	// Update context to one returned by builder (might contain WASI setup)
	r.ctx = ctx

	// Configure module config with args/env, no explicit stdin/stdout/stderr as WASI system handles it
	modConf := wazero.NewModuleConfig()
	if len(args) > 0 {
		modConf = modConf.WithArgs(args...)
	}
	for k, v := range env {
		modConf = modConf.WithEnv(k, v)
	}

	instance, err := r.runtime.InstantiateModule(ctx, r.mod, modConf)
	if err != nil {
		return fmt.Errorf("failed to instantiate WASM module: %w", err)
	}
	defer instance.Close(ctx)

	return nil
}

// invokeJS handles JavaScript execution
func (r *Runtime) invokeJS(stdin io.Reader, env map[string]string, script []byte, args ...string) error {
	if len(script) == 0 {
		return fmt.Errorf("script argument is required for JS runtime")
	}

	var jsRuntimeArguments []string
	jsRuntimeArguments = append(jsRuntimeArguments, "", "-e", string(script))
	args = append(jsRuntimeArguments, args...)

	modConf := wazero.NewModuleConfig().
		WithStdin(stdin).
		WithStdout(r.stdout).
		WithStderr(os.Stderr).
		WithArgs(args...)

	// Add environment variables
	for k, v := range env {
		modConf = modConf.WithEnv(k, v)
	}

	instance, err := r.runtime.InstantiateModule(r.ctx, r.mod, modConf)
	if err != nil {
		return fmt.Errorf("failed to instantiate JS module: %w", err)
	}
	defer instance.Close(r.ctx)

	return nil
}

// Close shuts down the runtime and releases resources.
func (r *Runtime) Close() error {
	if r.runtime != nil {
		if err := r.runtime.Close(r.ctx); err != nil {
			return fmt.Errorf("failed to close runtime: %w", err)
		}
	}
	return nil
}
