import { Command } from 'func'

@Command({ name: 'start' })
export class Start {
  constructor() {
    require('@/public/scripts/start')
  }
}
