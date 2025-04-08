import { toTimeDuration } from '@/utils/Formatter'

const transformer = {
  runTime: str => toTimeDuration(str),
  estimatedRunTime: str => toTimeDuration(str),
  remainingTime: str => toTimeDuration(str),
}
export default transformer
