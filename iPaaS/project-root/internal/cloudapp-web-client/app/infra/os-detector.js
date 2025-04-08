/**
 * Detect whether the running os is mac os
 * @return {Boolean}
 */
export function isMacOS() {
  const { appVersion } = window.navigator
  return appVersion.indexOf('Mac') >= 0
}
