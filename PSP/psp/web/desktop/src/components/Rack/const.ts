import Gradient from 'javascript-color-gradient'

export const rowHeight = 20 // 1 U 的高度
export const rackWidth = 215 // 机架宽度
const colorRange = ['#10E617', '#E6D610', '#E61010']

const phase1 = new Gradient()
phase1.setMidpoint(51)
phase1.setGradient(colorRange[0], colorRange[1])

const phase2 = new Gradient()
phase2.setMidpoint(51)
phase2.setGradient(colorRange[1], colorRange[2])

export const colors = [...phase1.getArray(), ...phase2.getArray()]
