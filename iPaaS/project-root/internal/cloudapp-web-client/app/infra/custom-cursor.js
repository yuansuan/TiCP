const DEFAULT_CURSOR_IMG_IDX = -1

export default class CustomCursor {
  /**
   * @param {HTMLElement} elem
   */
  constructor(elem) {
    this._elem = elem
    // _imgM is used to store cursor image for cache
    // here every image has an unique index
    this._currentIdx = DEFAULT_CURSOR_IMG_IDX
    this._imgM = new Map([[DEFAULT_CURSOR_IMG_IDX, 'default']])
  }

  static genImgData(data, hotspotX = 0, hotspotY = 0) {
    return `url(data:image/png;base64,${data}) ${hotspotX} ${hotspotY}, auto`
  }

  /**
   * set an image for the cursor.
   * (hotspotX, hotspotY, data) must be set together
   * the imageData of idx will be set as the new if (hotspotX, hotspotY, data) exist
   * @param {Number} idx cursor image idx and >= 0 must be satisfied
   * @param {Number|null} hotspotX cursor's hotspot x
   * @param {Number|null} hotspotY cursor's hotspot y
   * @param {String|null} data base64 encoded png data
   */

  hideCursor(idx){
    this._elem.style.cursor = 'none'
    this._imgM.set(idx,null)
    this._currentIdx = idx
  }
  setImage(idx, hotspotX = null, hotspotY = null, data = null) {
    if (this._currentIdx === idx) {
      return
    }

    if (data != null) {
      const imgData = this.constructor.genImgData(data, hotspotX, hotspotY)
      this._imgM.set(idx, imgData)
    }

    this._currentIdx = idx
    if (!this._imgM.has(idx)) {
      this._currentIdx = DEFAULT_CURSOR_IMG_IDX
    }

    this._elem.style.cursor = this._imgM.get(idx)
  }
}
