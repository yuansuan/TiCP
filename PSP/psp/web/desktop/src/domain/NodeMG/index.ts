import NodeManager from './NodeManager'

export const nodeManager = new NodeManager()

export enum NodeActionLabel {
  unKnown = 'NodeActionUnknown',
  open = 'node_start',
  close = 'node_close',
}
