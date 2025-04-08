import test from 'ava'
import { $spawn } from '../utils/_child'
import { version } from '@/package'

test('option version should print version', async t => {
  try {
    const data = await $spawn(['--version'])
    t.regex(data, new RegExp(version))
  } catch (e) {
    t.fail()
  }
})
