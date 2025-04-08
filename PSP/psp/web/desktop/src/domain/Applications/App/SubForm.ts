import { observable, action } from 'mobx'

import Section, { IRequest as ISectionRequest } from './Section'

export interface IRequest {
  section: ISectionRequest[]
}

interface ISubForm {
  sections: Section[]
}

export default class SubForm implements ISubForm {
  @observable sections: Section[] = []

  constructor(request?: IRequest) {
    request && this.init(request)
  }

  @action
  init = (request: IRequest) => {
    Object.assign(this, {
      sections: (request.section || []).map((item) => new Section(item))
    })
  }

  // add section
  @action
  add = (section: Section, index?: number) => {
    if (index === undefined) {
      this.sections.push(section)
    } else {
      this.sections.splice(index, 0, section)
    }
  }

  // delete section
  @action
  delete = (index: number) => this.sections.splice(index, 1)

  // toggle section
  @action
  toggle = (sourceIndex: number, targetIndex: number) => {
    if (sourceIndex === targetIndex) {
      return
    }

    const sourceSection = this.sections[sourceIndex]
    if (!sourceSection) {
      throw new Error(`SourceIndex: ${sourceIndex} is not valid index`)
    }

    const targetSection = this.sections[targetIndex]
    if (!targetSection) {
      throw new Error(`TargetIndex: ${targetIndex} is not valid index`)
    }

    this.sections = this.sections.map((item, index) => {
      if (index === sourceIndex) {
        return targetSection
      }

      if (index === targetIndex) {
        return sourceSection
      }

      return item
    })
  }

  toRequest(): IRequest {
    return {
      section: this.sections.map((item) => item.toRequest())
    }
  }
}
