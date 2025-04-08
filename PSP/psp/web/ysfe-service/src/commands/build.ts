import { Command, SubOptions, CommandArgsProvider } from 'func'

@Command({ name: 'build' })
@SubOptions([
  {
    name: 'config',
    alias: 'C',
    description: 'custom webpack config for building',
    type: String,
  },
])
export class Build {
  constructor({ option }: CommandArgsProvider) {
    const build = require('@/public/scripts/build')
    build(option)
  }
}
