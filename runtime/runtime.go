package runtime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:generate stringer --type RuntimeEngine
type RuntimeEngine int

const (
	RuntimeEngineGo RuntimeEngine = iota
	RuntimeEngineJS
)

// Args defines the configuration for creating a new Runtime instance.
type Args struct {
	Stdout       io.Writer
	DeploymentID uuid.UUID
	Engine       RuntimeEngine
	Blob         []byte
	Cache        wazero.CompilationCache
}

// Runtime manages the WebAssembly execution environment.
type Runtime struct {
	stdout       io.Writer
	ctx          context.Context
	deploymentID uuid.UUID
	engine       RuntimeEngine
	mod          wazero.CompiledModule
	runtime      wazero.Runtime
}

// New initializes a new WebAssembly runtime with the given arguments.
func New(ctx context.Context, args Args) (*Runtime, error) {
	config := wazero.NewRuntimeConfigCompiler().WithCompilationCache(args.Cache)
	rt := wazero.NewRuntimeWithConfig(ctx, config)

	// Ensure WASI is available
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	mod, err := rt.CompileModule(ctx, args.Blob)
	if err != nil {
		rt.Close(ctx) // Cleanup on failure
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	return &Runtime{
		runtime:      rt,
		ctx:          ctx,
		deploymentID: args.DeploymentID,
		engine:       args.Engine,
		stdout:       args.Stdout,
		mod:          mod,
	}, nil
}

// Invoke executes the compiled WebAssembly module with provided input and environment variables.
func (r *Runtime) Invoke(stdin io.Reader, env map[string]string, script []byte, args ...string) error {
	switch r.engine {
	case RuntimeEngineGo:
		// do nothing
	case RuntimeEngineJS:
		var jsRuntimeArguments []string
		if len(script) == 0 {
			return fmt.Errorf("script argument is required for JS runtime")
		}
		jsRuntimeArguments = append(jsRuntimeArguments, "", "-e", string(script))
		args = append(jsRuntimeArguments, args...)
	default:
		return fmt.Errorf("invalid runtime engine %d", r.engine)
	}

	modConf := wazero.NewModuleConfig().
		WithStdin(stdin).
		WithStdout(r.stdout).
		WithStderr(os.Stderr).
		WithArgs(args...)

	// Set environment variables
	for k, v := range env {
		modConf = modConf.WithEnv(k, v)
	}

	_, err := r.runtime.InstantiateModule(r.ctx, r.mod, modConf)
	if err != nil {
		return fmt.Errorf("failed to instantiate module: %w", err)
	}

	return nil
}

// Close shuts down the runtime and releases resources.
func (r *Runtime) Close() error {
	if err := r.runtime.Close(r.ctx); err != nil {
		return fmt.Errorf("failed to close runtime: %w", err)
	}
	return nil
}
