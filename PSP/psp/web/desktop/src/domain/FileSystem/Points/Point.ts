import { BaseDirectory } from '@/utils/FileSystem'

export interface IPoint {
  rootPath: string
  service: {
    get?: (path: string) => Promise<any[]>
    fetch: (path: string) => Promise<any[]>
    fetchTree: (options: { path: string; rootPath: string }) => Promise<any[]>
    delete: (options: { paths: string[] }) => Promise<boolean>
    rename: (options: { path: string; newName: string }) => Promise<boolean>
    createDir: (options: { path: string }) => Promise<boolean>
    move: (options: {
      dstpath: string
      srcpaths: string[]
      overwrite: boolean
    }) => Promise<boolean>
    copy: (options: {
      path: string
      names: string[]
      toPath: string
    }) => Promise<boolean>
    preDownload: (options: { paths: string[]; userId: number }) => Promise<any>
    download: ({ token, userId }) => void
    compress: (options: {
      path: string
      names: string[]
      compressType: string
      zipName: string
    }) => Promise<boolean>
    extract: (options: {
      file: string
      toPath: string
      fileType: string
    }) => Promise<boolean>
    view: (options: {
      path: string
      offset: number
      len: number
    }) => Promise<string>
    edit: (options: { path: string; content: string }) => Promise<any>
    exist: (paths: string[]) => Promise<boolean[]>
  }
  [key: string]: any
}

type Point = IPoint & BaseDirectory
export default Point
