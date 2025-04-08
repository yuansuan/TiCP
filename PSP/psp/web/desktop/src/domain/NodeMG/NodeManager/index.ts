import { observable, action } from 'mobx'
import { Http } from '@/utils'
import Node from './Node'

export default class NodeManager {
  @observable nodeList: Node[] = []
  @observable nodeName: string = ''
  @observable pageIndex: number = 1
  @observable pageSize: number = 10

  @observable total: number = 0

  @action
  async getNodeList() {
    const res = await Http.get('/node/list', {
      params: {
        node_name: this.nodeName,
        page_index: this.pageIndex,
        page_size: this.pageSize
      }
    })
    this.total = res.data.total
    this.nodeList = res.data?.list?.map(item => {
      const node = new Node(item)
      return node
    })

    return res
  }

  operate = (node_names: string[], action) => {
    return Http.post('/node/operate', {
      node_names: node_names,
      operation: action
    })
  }

  getClusterCoreInfo() {
    return Http.get('/node/coreNum')
  }
}
