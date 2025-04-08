#!/usr/bin/env node

/**
 * module-alias 在这里有两个用处
 * 1. 提供模块 require 的 alias 解析；
 * 2. func-service build 以后会改变 src 的目录结构，通过 module-alias 保证 src 与 static 文件的引用一致性；
 */
const moduleAlias = require('module-alias')
moduleAlias.addAlias('@', require('path').resolve(__dirname, '..'))

import * as commands from './commands'
import * as options from './options'
import { Container } from 'func'

const modules = Object.assign({}, commands, options)
const params = Object.keys(modules).map(key => modules[key])
new Container(params)
