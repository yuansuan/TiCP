/* Copyright (C) 2016-present, Yuansuan.cn */

/**
 * 计算传进来的数据，不让其移出可视区域
 * @param data 需要计算的数据
 * @param trimDistance 裁剪框宽度
 * @param canvasDistance 画布宽度
 */
export function confine(
  data: number,
  trimDistance: number,
  canvasDistance: number
) {
  if (noNegative(data) + trimDistance > canvasDistance) {
    return noNegative(canvasDistance - trimDistance)
  } else {
    return noNegative(data)
  }
}

/**
 * 对参数进行处理，小于0则返回0
 */
export function noNegative(data: number) {
  return data > 0 ? data : 0
}

export default {
  confine,
  noNegative
}
