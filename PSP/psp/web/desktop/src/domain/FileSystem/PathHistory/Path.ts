import { observable, action } from 'mobx'
import { Point } from '../Points'

export default class Path {
  @observable source: Point
  @observable path: string

  constructor({ source, path }) {
    this.source = source
    this.path = path
  }

  @action
  updatePath = (path) => (this.path = path)
}
