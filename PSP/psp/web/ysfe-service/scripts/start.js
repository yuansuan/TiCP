const TscWatchClient = require('tsc-watch/client')
const watch = new TscWatchClient()
const child_process = require('child_process')
const fs = require('fs-extra')

fs.emptyDirSync('dist')
fs.ensureDirSync('dist')
fs.copySync('public', 'dist/public')
fs.copy('package.json', 'dist/package.json')

watch.on('first_success', () => {
  console.log('npm unlink')
  let stdout = child_process.execSync('npm unlink')
  console.log(stdout.toString())

  console.log('npm link')
  stdout = child_process.execSync('npm link')
  console.log(stdout.toString())
})

watch.on('subsequent_success', () => {
  console.log('tsc-watch compile successful!')
})

watch.start('--project', '.')
