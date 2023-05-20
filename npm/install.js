const fs = require('fs');
const path = require('path');
const os = require('os');
const releaseDir = path.join(__dirname, 'release');
const binDir = path.join(__dirname, 'bin');
const platform = os.platform();
let arch = os.arch();

// Map x64 to amd64
if (arch === 'x64') {
  arch = 'amd64';
}

let binaryName;
const version = "v0.7.17"
if (platform === 'win32') {
  binaryName = `tzap-${version}-windows-${arch}.exe`;
} else if (platform === 'darwin') {
  binaryName = `tzap-${version}-darwin-${arch}`;
} else if (platform === 'linux') {
  binaryName = `tzap-${version}-linux-${arch}`;
} else {
  console.error(`Unsupported platform: ${platform}`);
  process.exit(1);
}

const sourcePath = path.join(releaseDir, binaryName);
const targetPath = path.join(binDir, 'tzap');
if (!fs.existsSync(binDir)){
  fs.mkdirSync(binDir);
}
fs.copyFileSync(sourcePath, targetPath);
fs.chmodSync(targetPath, 0o755);