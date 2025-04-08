/* Copyright (C) 2016-present, Yuansuan.cn */
String.prototype.strip = function (c) {
  let i = 0,
    j = this.length - 1
  while (this[i] === c) i++
  while (this[j] === c) j--
  return this.slice(i, j + 1)
}

String.prototype.count = function (c) {
  let result = 0,
    i = 0
  for (i; i < this.length; i++) if (this[i] == c) result++
  return result
}
export class Item {
  constructor({ type, name, info, data, host, id, path, size, mtime, isFile,readOnly }) {
    switch (true) {
      case /jpe?g|png|gif/gi.test(type):
        type = 'img'
        break
      case /zip|rar/gi.test(type):
        type = 'zip'
        break
      case /MP[234]|WAV|WMA|Flac|MIDI|RA|APE|AAC|CDA|MOV/gi.test(type):
        type = 'mp4'
        break
      case type === 'pdf':
        type = 'pdf'
        break
      case /txt|dmg/gi.test(type):
        break
      case !/folder/gi.test(type):
        type = 'unknownfile'
        break
      default:
        break
    }
    // if (/je?pg|png|gif/ig.test(type)) {
    //   type = 'img'
    // } else if (/MP[234]|WAV|WMA|Flac|MIDI|RA|APE|AAC|CDA|MOV/ig.test(type)) {
    //   type = 'mp4'
    // } else if (type === 'pdf') {
    //   type = 'pdf'
    // } else if (/txt|dmg/ig.test(type)) {
    //   type = type
    // } else if (!/folder/ig.test(type)) {
    //   type = 'unknownfile'
    // }
    this.type = type?.toLocaleLowerCase() || 'folder'
    this.name = name
    this.info = info || {}
    this.info.icon = this.info.icon || this.type
    this.data = data
    this.host = host
    this.id = this.gene()
    this.originId = id
    this.readOnly = readOnly
    this.path = path
    this.size = size
    this.mtime = mtime
    this.editFlag = false
    this.isFile = isFile
  }

  gene() {
    return Math.random().toString(36).substring(2, 10).toLowerCase()
  }

  getId() {
    return this.id
  }

  getData() {
    return this.data
  }

  setData(data) {
    this.data = data
  }
}

export class Bin {
  constructor() {
    this.tree = []
    this.lookup = {}
    this.special = {}
  }

  setSpecial(spid, id) {
    this.special[spid] = id
  }

  setId(id, item) {
    this.lookup[id] = item
  }

  getId(id) {
    return this.lookup[id]
  }

  getPath(id) {
    let cpath = ''
    let curr = this.getId(id)

    while (curr) {
      cpath = curr.name + '\\' + cpath
      curr = curr.host
    }

    return cpath.count('\\') > 1 ? cpath.strip('\\') : cpath
  }

  parsePath(cpath) {
    if (cpath.includes('%')) {
      return this.special[cpath.trim()]
    }

    cpath = cpath
      .split('\\')
      .filter(x => x !== '')
      .map(x => x.trim().toLowerCase())
    if (cpath.length === 0) return null

    let pid = null,
      curr = null
    for (let i = 0; i < this.tree.length; i++) {
      if (this.tree[i].name.toLowerCase() === cpath[0]) {
        curr = this.tree[i]
        break
      }
    }

    if (curr) {
      let i = 1,
        l = cpath.length
      while (curr.type?.toLocaleLowerCase() === 'folder' && i < l) {
        let res = true
        for (let j = 0; j < curr.data.length; j++) {
          if (curr.data[j].name.toLowerCase() === cpath[i]) {
            i += 1
            if (curr.data[j].type?.toLocaleLowerCase() === 'folder') {
              res = false
              curr = curr.data[j]
            }

            break
          }
        }

        if (res) break
      }

      if (i === l) pid = curr.id
    }

    return pid
  }

  parseFolder(data, name, host = null) {
    // 生成一个文档对象
    let item = new Item({
      type: data.type,
      name: data.name || name,
      info: data.info,
      host: host,
      id: data.id,
      path: data.path,
      size: data.size,
      mtime:data.mtime,
      isFile: data.isFile,
      readOnly: data.only_read
    })

    // 收集到lookup中，id是随机生成的
    this.setId(item.id, item)

    // 特殊映射
    if (data.info && data.info.spid) {
      this.setSpecial(data.info.spid, item.id)
    }

    if (item.type?.toLocaleLowerCase() !== 'folder') {
      item.setData(data.data)
    } else {
      let fdata = []
      if (data.data) {
        for (const key of Object.keys(data.data)) {
          fdata.push(this.parseFolder(data.data[key], key, item))
        }
      }

      item.setData(fdata)
    }

    return item
  }

  parse(data, id = '') {
    let drives = Object.keys(data)
    let tree = this.tree || []
    for (let i = 0; i < drives.length; i++) {
      tree.push(this.parseFolder(data[drives[i]]))
    }

    this.tree = tree
  }
}
