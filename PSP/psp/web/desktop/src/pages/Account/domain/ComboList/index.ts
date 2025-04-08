import { observable, action, runInAction } from 'mobx'
import { Combo, BaseCombo } from './Combo'
import { accountServer } from '@/server'

export class BaseComboList {
  @observable list: Combo[] = []
}

type Request = Omit<BaseComboList, 'list'> & {
  list: BaseCombo[]
}

export class ComboList extends BaseComboList {
  constructor(props?: Partial<Request>) {
    super()

    if (props) {
      this.update(props)
    }
  }

  @action
  update({ list, ...props }: Partial<Request>) {
    Object.assign(this, props)

    if (list) {
      this.list = list.map(item => new Combo(item))
    }
  }

  fetch = async (product_id?: string) => {
    const { data } = await accountServer.getComboList(product_id)

    runInAction(() => {
      this.update({
        list: data.item
      })
    })
  }
}
