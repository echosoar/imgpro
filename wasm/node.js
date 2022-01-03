
const { readFileSync } = require('fs');
const { resolve } = require('path');
const { execSync } = require('child_process');
const wasmExec = execSync('go env GOROOT').toString().replace(/\n/g, '') + '/misc/wasm/wasm_exec.js';
require(wasmExec);
const wasmFile = resolve(__dirname, './imgpro.wasm');
const run = async (file) => {
  const buf = readFileSync(file);
  const go = new global.Go();
  return WebAssembly.instantiate(buf, go.importObject).then(result => {
    go.run(result.instance);
  }).catch(e => {
  });
}
run(wasmFile).then(() => {
  const cacheFileBuffer = readFileSync(resolve(__dirname, '../test/imgs/cool_88.webp'))
  let result;
  try {
    result = global.imgExec([
      "size",
      "type",
      "width",
      "height",
      "frame",
      "hue",
      "exif"
    ], cacheFileBuffer);
    console.log(JSON.stringify(result, null, 2));
  } catch(e) {
    console.log('error: ' + e.message)
  }
});