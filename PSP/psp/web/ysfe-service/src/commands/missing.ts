import { CommandMissing } from 'func'

@CommandMissing()
export class Missing {
  constructor() {
    console.log('command missing')
  }
}
