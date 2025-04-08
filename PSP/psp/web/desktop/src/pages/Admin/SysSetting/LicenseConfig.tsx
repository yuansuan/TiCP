import * as React from 'react'
import { observer } from 'mobx-react'
import { Button } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import license from '@/domain/License'
import { sysConfig } from '@/domain'
import { Modal } from '@/components'
import { observable } from 'mobx'
import LicenseForm from './LicenseForm'

let Confetti = null

if (process.env.BABEL_ENV === 'development') {
  Confetti = require('@/utils/Confetti')
}

@observer
export default class LicenseConfig extends React.Component<any> {
  @observable loading = false
  textRef = null

  constructor(props) {
    super(props)
    this.textRef = React.createRef()
  }

  async componentDidMount() {
    this.loading = true
    try {
      await license.getLicenseInfo()
    } finally {
      this.loading = false
    }
  }

  updataLicense = () => {
    Modal.show({
      title: '续期',
      footer: null,
      content: ({ onCancel }) => {
        const successCallback = async () => {
          await license.getLicenseInfo()
          onCancel()
        }
        return <LicenseForm successCallback={successCallback} />
      },
      width: 800
    })
  }

  toggle = () => {
    if (process.env.BABEL_ENV === 'development') {
      Confetti && Confetti.toggle()
      this.textRef.current.classList.toggle('animate-charcter')
    }
  }

  componentWillUnmount(): void {
    if (process.env.BABEL_ENV === 'development') {
      Confetti && Confetti.destory()
    }
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label align={'left'}>软件名</Label>
          <span
            ref={this.textRef}
            className='textField'
            onClick={() => this.toggle()}>
            {sysConfig.websiteConfig?.title || '远算云仿真平台'}
          </span>
        </div>
        {/* <div className='item'>
          <Label align={'left'}>版本</Label>
          <span className='textField'>{license.info?.version || '无'}</span>
        </div> */}
        <div className='item'>
          <Label align={'left'}>过期时间</Label>
          <span className='textField'>
            {license.info?.expiry || '无'}
            <span style={{ paddingLeft: 20, paddingRight: 20 }} />
            软件还剩 {license.info?.available_days || 0} 天过期
            <Button
              style={{ marginLeft: 20, height: 30 }}
              loading={this.loading}
              onClick={this.updataLicense}>
              续期
            </Button>
          </span>
        </div>
      </ConfigWrapper>
    )
  }
}
