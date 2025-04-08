import pkg from '@/package'

if (!pkg.bin) {
  throw new Error(
    'Not found item "bin" in "package.json", maybe you need run "npm start" or "npm build".'
  )
}

export const bin = Object.keys(pkg.bin)[0]
