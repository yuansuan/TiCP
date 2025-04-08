import * as React from 'react'
import { observable, action, computed } from 'mobx'
import { sysConfig } from '@/domain'
import { observer } from 'mobx-react'
import { Input, Checkbox, message } from 'antd'
import { Button, Modal } from '@/components'
import { PasswdForm } from '@/components'
import { Subject, combineLatest, from, empty } from 'rxjs'
import Vcode from 'react-vcode'
import {
  withLatestFrom,
  filter,
  tap,
  map,
  switchMap,
  finalize,
  catchError
} from 'rxjs/operators'
import { untilDestroyed } from '@/utils/operators'
import { Auth } from '@/domain'
import { Wrapper } from './style'
import LicenseForm from '@/pages/Admin/SysSetting/LicenseForm'
const codes = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789'
const vCodeNum = sysConfig.userConfig.show_verify_code
const MAX_LOGIN_ERROR_NUM = typeof vCodeNum !== 'number' ? 3 : vCodeNum

@observer
class Login extends React.Component<any> {
  stateModel = observable({
    loading: false,
    disabled: false,
    update: action(function (this: any, props) {
      Object.assign(this, props)
    })
  })

  @observable inputVCode = null
  @observable loginErrNum = 0
  vcodeRef = null
  vcode = null

  @computed
  get isShowVCode() {
    // 0: 永久不显示
    // -1: 永久显示
    // 1: 出错 1 次后显示
    // 3: 出错 3 次后显示

    if (MAX_LOGIN_ERROR_NUM === 0) return false

    return this.loginErrNum >= MAX_LOGIN_ERROR_NUM
  }

  username$ = new Subject()
  password$ = new Subject()
  login$ = new Subject()

  componentDidMount() {
    const currentPath = localStorage.getItem('CURRENTROUTERPATH')
    if (currentPath) {
      localStorage.removeItem('CURRENTROUTERPATH')
    }
    // 用户名、密码信息
    const loginInfo$ = combineLatest(this.username$, this.password$).pipe(
      untilDestroyed(this)
    )

    // 用户名/密码为空，登录按钮置灰
    loginInfo$
      .pipe(filter(([username, password]) => !username || !password))
      .subscribe(() => this.stateModel.update({ disabled: true }))

    loginInfo$
      .pipe(filter(([username, password]) => !!(username && password)))
      .subscribe(() => this.stateModel.update({ disabled: false }))

    // 登录
    this.login$
      .pipe(
        untilDestroyed(this),
        withLatestFrom(loginInfo$),
        map(([, loginInfo]) => loginInfo),
        filter(([username, password]) => !!(username && password)),
        tap(() =>
          this.stateModel.update({
            loading: true
          })
        ),
        switchMap(([username, password]) =>
          from(Auth.login(username, password)).pipe(
            // 捕获错误，防止错误中断 login$ 流
            catchError(code => {
              this.loginErrNum += 1
              this.vcodeRef && this.vcodeRef.onClick()
              if (code === 16036) {
                Modal.show({
                  title: '您的许可证无效，请按照下面的步骤进行激活',
                  content: ({ onCancel, onOk }) => {
                    const msg = '新许可证激活成功, 请重新登录'
                    const successCallback = async () => {
                      onCancel()
                    }
                    return (
                      <LicenseForm
                        successCallback={successCallback}
                        successMsg={msg}
                      />
                    )
                  },
                  bodyStyle: { background: '#fff' },
                  width: 800,
                  footer: null
                })
              }
              return empty()
            }),
            finalize(() =>
              this.stateModel.update({
                loading: false
              })
            )
          )
        )
      )
      .subscribe(res => {
        const { user } = res.data
        if (user.id) {
          localStorage.setItem('userId', user.id)
          localStorage.removeItem('needLogin')
          window.location.reload()
        } else if (user.changePwd) {
          Modal.show({
            title: '修改密码',
            content: ({ onCancel, onOk }) => (
              <PasswdForm
                tips={'为了确保您的安全，请修改您的密码'}
                username={user.name}
                onCancel={onCancel}
                onOk={onOk}
              />
            ),
            bodyStyle: { background: '#fff' },
            width: 800,
            footer: null
          })
        } else if (user.licenseExpired) {
          Modal.show({
            title: '您的许可证已过期，请按照下面的步骤进行续期',
            content: ({ onCancel, onOk }) => {
              const msg = '许可证续期成功, 请重新登录'
              const successCallback = async () => {
                onCancel()
              }
              return (
                <LicenseForm
                  successCallback={successCallback}
                  successMsg={msg}
                />
              )
            },
            bodyStyle: { background: '#fff' },
            width: 800,
            footer: null
          })
        }
      })

    // get token from url params
    // atob() btoa() "YWRtaW46YWRtaW4="
    // http://0.0.0.0:8080/#/login?token=YWRtaW46YWRtaW4=
    const params = new URLSearchParams(location.hash.split('?')[1])

    if (params.get('token')) {
      const token = atob(params.get('token'))
      const temps = token.split(':')
      this.username$.next(temps[0])
      this.password$.next(temps[1])
      this.login$.next()
    }
  }

  handleLogin = e => {
    e.preventDefault()
    if (
      this.isShowVCode &&
      this.vcode?.toLowerCase() !== this.inputVCode?.toLowerCase()
    ) {
      message.error('验证码错误')
      this.vcodeRef && this.vcodeRef.onClick()
      return
    }
    this.login$.next()
  }

  onInputVCodeChange = e => {
    this.inputVCode = e.target.value
  }

  onUsernameChange = e => this.username$.next(e.target.value.trim())

  onPasswordChange = e => this.password$.next(e.target.value)

  render() {
    const { loading, disabled } = this.stateModel

    const logo =
      sysConfig.websiteConfig?.loginPage?.copyRightLogoUrl ||
      require('@/assets/images/loginlogo.png')

    const bigBg =
      sysConfig.websiteConfig?.loginPage?.bgUrl ||
      require('@/assets/images/bigbg.png')

    const smallBg =
      sysConfig.websiteConfig?.loginPage?.leftBgUrl ||
      require('@/assets/images/smallbg.jpg')

    return (
      <Wrapper bigBg={bigBg} smallBg={smallBg} logo={logo}>
        <div className='centerBox'>
          <div className='bgBox' />
          <div
            className='loginBox'
            style={{
              marginBottom: 'inherit',
              paddingBottom: this.isShowVCode ? 0 : 60
            }}>
            <div className='title'>
              <p>欢迎使用</p>
              <h2>{sysConfig.websiteConfig?.title || '远算云仿真平台'}</h2>
            </div>
            <div className='form' onSubmit={this.handleLogin}>
              <div className='field'>
                <label className='label'>用户名:</label>
                <Input
                  autoFocus
                  placeholder='请输入'
                  onChange={this.onUsernameChange}
                  data-testid='login-name'
                />
              </div>
              <div className='field'>
                <label className='label'>密码:</label>
                <Input.Password
                  placeholder='请输入'
                  onPressEnter={this.handleLogin}
                  onChange={this.onPasswordChange}
                  data-testid='login-password'
                />
              </div>
              {/* TODD 验证码功能 */}
              {this.isShowVCode && false ? (
                <div className='field'>
                  <label className='label'>验证码:</label>
                  <Input
                    placeholder='请输入'
                    value={this.inputVCode}
                    onPressEnter={this.handleLogin}
                    onChange={this.onInputVCodeChange}
                    style={{ width: 150 }}
                  />
                  <Vcode
                    length={6}
                    onChange={v => {
                      this.vcode = v
                    }}
                    style={{ margin: '0 10px' }}
                    width={180}
                    options={{
                      codes: codes.split(''),
                      fontSizeMin: 20,
                      fontSizeMax: 22,
                      fonts: [
                        'Times New Roman',
                        'Georgia',
                        'Serif',
                        'sans-serif',
                        'arial'
                      ]
                    }}
                    ref={obj => (this.vcodeRef = obj)}
                  />
                </div>
              ) : null}
              <div>
                {/* 需要与后端配合生成 rememberToken 实现,
                如果单纯前端明文存储密码，有安全隐患，因此暂时先隐藏 */}
                <Checkbox className='remember'>记住密码</Checkbox>
              </div>
              <div>
                <Button
                  htmlType='submit'
                  className='loginBtn'
                  loading={loading}
                  disabled={disabled}
                  onClick={this.handleLogin}
                  data-testid='login-button'
                  type='primary'>
                  登录
                </Button>
              </div>
            </div>
            <div style={{ textAlign: 'center', color: 'red' }}>
              {sysConfig.customConfig?.loginPage?.texts?.warning || ''}
            </div>
          </div>
        </div>
        <div className='footerBox'>
          <div className='ysLogo'>
            <div className='logo' />
          </div>
          <p className='copyRight'>
            {sysConfig.websiteConfig?.loginPage?.copyRightText ||
              '© 2016-2023 Yuansuan.cn All rights reserved.'}
          </p>
        </div>
      </Wrapper>
    )
  }
}

export default Login
