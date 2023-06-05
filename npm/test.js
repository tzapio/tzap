require('./wasm_exec.js')
const fs = require('fs');
const go = new Go(); // Defined in wasm_exec.js. Don't forget to add this in your index.html.
const wasmBuffer = fs.readFileSync('./release/tzap.wasm');
const runWasmAdd = async () => {
  // Get the importObject from the go instance.
  const importObject = go.importObject;
  go.argv = process.argv.slice(1)
  // Instantiate our wasm module
  const wasmModule = await WebAssembly.instantiate(wasmBuffer, importObject);
  // Allow the wasm_exec go instance, bootstrap and execute our wasm module
  go.run(wasmModule.instance);

};
runWasmAdd();