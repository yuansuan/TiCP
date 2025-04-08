import uuid from 'uuid/v4'

class IdGen {
  /**
   * @return {String}
   */
  gen() {
    return uuid()
  }
}

export default new IdGen()
