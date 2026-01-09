// Simple build script for copying files
const fs = require('fs');
const path = require('path');

const srcDir = path.join(__dirname, 'src');
const distDir = path.join(__dirname, 'dist');

// Ensure dist directory exists
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true });
}

// Copy files
const files = ['index.html', 'style.css', 'app.js'];
files.forEach(file => {
  const srcFile = path.join(srcDir, file);
  const distFile = path.join(distDir, file);
  if (fs.existsSync(srcFile)) {
    fs.copyFileSync(srcFile, distFile);
    console.log(`Copied ${file}`);
  }
});

console.log('Build complete!');
