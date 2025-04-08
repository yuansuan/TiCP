const child_process = require('child_process')
const fs = require('fs-extra')

fs.emptyDirSync('dist')
fs.ensureDirSync('dist')
child_process.execSync('npx tsc')
fs.copySync('public', 'dist/public')

// write package.json
const pkg = fs.readJSONSync('package.json')
delete pkg.files
fs.writeJSONSync('dist/package.json', pkg)
