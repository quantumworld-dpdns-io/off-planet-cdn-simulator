pub struct WasmRuntime;

impl WasmRuntime {
    pub fn new() -> Self { Self }

    pub fn run_plugin(&self, _plugin_path: &str, _input: &[u8]) -> anyhow::Result<f64> {
        // TODO: implement Wasmtime/WASI execution (Phase 5)
        Ok(0.0)
    }
}
