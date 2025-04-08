import { runInAction } from "mobx";
import { Hardware, ListHardwareRequest } from ".";
import { vis } from "..";

export class HardwareList {
  list: Array<Hardware> = []

  fetch() {
    runInAction(async () => {
      this.list = (await vis.listHardware(new ListHardwareRequest({}))).items
    })
    return this.list
  }
}
