/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import WorkStationList from './WorkStationList'
import WorkTaskList from './WorkTaskList'
import VirtualMachineList from './VirtualMachineList'
import MachineNodeList from './MachineNodeList'
import SoftwareAppList from './SoftwareAppList'
import VirtualMachineSetting from './VirtualMachineSetting'
export const workStationList = new WorkStationList()
export const workTaskList = new WorkTaskList()
export const vmList = new VirtualMachineList()
export const nodeList = new MachineNodeList()
export const sofwareAppList = new SoftwareAppList()
export const machineSetting = new VirtualMachineSetting()
export const VM_OS_TYPE = {
  centos7: 'CentOS 7',
  win7: 'Windows 7',
  win10: 'Windows 10',
}
