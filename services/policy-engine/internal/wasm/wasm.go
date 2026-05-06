package wasm

// RunPlugin executes a WASM scoring plugin at pluginPath with the provided
// input bytes and returns the computed score.
// TODO: integrate Wasmtime (github.com/bytecodealliance/wasmtime-go) in Phase 3.
func RunPlugin(pluginPath string, input []byte) (float64, error) {
	return 0.0, nil
}
