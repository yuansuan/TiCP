import { observable, action, computed } from 'mobx'
import nanoid from 'nanoid'

import Field, { IRequest as IFieldRequest } from './Field'

export interface IRequest {
  name: string
  field: IFieldRequest[]
}

interface ISection {
  _key: string
  name: string
  fields: Field[]

  editing: boolean
}

export default class Section implements ISection {
  @observable _key = nanoid()
  @observable name = ''
  @observable fields: Field[] = []

  @observable editing = false
  @action
  updateEditing = (editing) => (this.editing = editing)

  constructor(request?: IRequest) {
    if (request) {
      this.init(request)
    }

    if (!request) {
      this.updateEditing(true)
    }
  }

  @computed
  get computedEditing() {
    return this.editing || this.fields.some((item) => item.editing)
  }

  @action
  init = (request: IRequest) => {
    Object.assign(this, {
      name: request.name,
      fields: request.field.map((item) => new Field(item))
    })
  }

  @action
  updateName = (name) => (this.name = name)

  @action
  add(field: Field, index?: number) {
    if (index === undefined) {
      this.fields.push(field)
    } else {
      const fields = [...this.fields]
      fields.splice(index, 0, field)
      this.fields = fields
    }
  }

  @action
  delete = (index: number) => this.fields.splice(index, 1)

  // toggle field
  @action
  toggle = (sourceIndex: number, targetIndex: number) => {
    if (sourceIndex === targetIndex) {
      return
    }

    const sourceField = this.fields[sourceIndex]
    if (!sourceField) {
      throw new Error(`SourceIndex: ${sourceIndex} is not valid index`)
    }

    const targetField = this.fields[targetIndex]
    if (!targetField) {
      throw new Error(`TargetIndex: ${targetIndex} is not valid index`)
    }

    this.fields = this.fields.map((item, index) => {
      if (index === sourceIndex) {
        return targetField
      }

      if (index === targetIndex) {
        return sourceField
      }

      return item
    })
  }

  toRequest(): IRequest {
    return {
      name: this.name,
      field: this.fields.map((item) => item.toRequest())
    }
  }
}
