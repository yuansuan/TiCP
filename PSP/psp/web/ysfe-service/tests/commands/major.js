import test from 'ava'
import { $spawn } from '../utils/_child'


test('major should print text', async t => {
  try {
    const data = await $spawn()
    t.regex(data, /ok/)
  } catch (e) {
    t.fail()
  }
})
