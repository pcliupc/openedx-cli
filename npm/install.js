const https = require('https');
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');

const REPO = 'pcliupc/openedx-cli';

function getDownloadInfo() {
  const version = require('./package.json').version;
  const platform = process.platform === 'win32' ? 'windows' : process.platform;
  const arch = process.arch === 'x64' ? 'amd64' : process.arch;
  const ext = platform === 'windows' ? '.exe' : '';
  const filename = `openedx-${version}-${platform}-${arch}${ext}`;
  const url = `https://github.com/${REPO}/releases/download/v${version}/${filename}`;
  return { url, filename, ext };
}

function download(url) {
  return new Promise((resolve, reject) => {
    const get = (target) => {
      https.get(target, { headers: { 'User-Agent': 'node-fetch' } }, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          get(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`Download failed: HTTP ${res.statusCode} for ${target}`));
          return;
        }
        const chunks = [];
        res.on('data', (chunk) => chunks.push(chunk));
        res.on('end', () => resolve(Buffer.concat(chunks)));
        res.on('error', reject);
      }).on('error', reject);
    };
    get(url);
  });
}

async function main() {
  const { url, ext } = getDownloadInfo();
  const binDir = path.join(__dirname, 'bin');
  const dest = path.join(binDir, `openedx${ext}`);

  fs.mkdirSync(binDir, { recursive: true });

  console.log(`Downloading ${url}...`);
  const data = await download(url);
  fs.writeFileSync(dest, data);
  fs.chmodSync(dest, 0o755);
  console.log('openedx CLI installed successfully.');
}

main().catch((err) => {
  console.error('Failed to install @openedx/cli:', err.message);
  console.error('You can download manually from: https://github.com/' + REPO + '/releases');
  process.exit(1);
});
