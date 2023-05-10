#!/usr/bin/env node
const { spawn } = require('child_process');
const os = require('os');
const path = require('path');

const platform = os.platform();
let arch = os.arch();

// Map x64 to amd64
if (arch === 'x64') {
  arch = 'amd64';
}

let binaryName;
const version = "v0.7.12"
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

const binDir = path.join(__dirname, 'release');
const binaryPath = path.join(binDir, binaryName);

// Run the tzap binary
const tzap = spawn(binaryPath, process.argv.slice(2));
process.stdin.pipe(tzap.stdin);
tzap.stdout.on('data', (data) => {
  process.stdout.write(`${data}`);
});

tzap.stderr.on('data', (data) => {
  process.stderr.write(`${data}`);
});

tzap.on('close', (code) => {
  process.exit(code);
});