/* Copyright (C) 2016-present, Yuansuan.cn */

import { Clipper } from '.'
import utils from './utils'
import constants from './constants'

export type ClipRect = {
  startX: number
  startY: number
  width: number
  height: number
}

export type ClipNode = {
  x: number
  y: number
  width: number
  height: number
  cursor: string
  direction?: string
}

export type Region = {
  x: number
  y: number
  width: number
  height: number
  style?: any
  state?: any
}

export type RegionFactory = () => Region | Region[]

export type PluginProps = {
  core: Clipper
  utils: typeof utils
  constants: typeof constants
}
