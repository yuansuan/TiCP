import * as React from 'react'
import Prism from 'prismjs'
import InfoDiff from './InfoDiff'
import FormDiff from './FormDiff'
import BinPathDiff from './BinPathDiff'
import ArrayDiff from './ArrayDiff'
import { InlineDiff } from './Diff'
import { Wrapper, StyledSection } from './style'
import ReactDiffViewer from 'react-diff-viewer'
import { Modal, Button } from '@/components'
import 'prismjs/components/prism-bash'
interface IProps {
  diff: any
  script: any
}

export default class DiffInfo extends React.Component<IProps> {
  highlightSyntax = code => {
    return (
      <pre
        style={{ display: 'inline' }}
        dangerouslySetInnerHTML={{
          __html: Prism.highlight(code || '', Prism.languages.bash, 'bash')
        }}></pre>
    )
  }

  get scriptDiffStyle() {
    return {
      diffContainer: {
        width: '1200px'
      }
    }
  }

  render() {
    const { diff, script } = this.props
    const application = diff && diff.application ? diff.application : {}
    const iconData = diff && diff.icon_data ? diff.icon_data : null
    const scriptData = diff && diff.script_data ? diff.script_data : null
    const formDiff =
      application.sub_form && application.sub_form.section
        ? application.sub_form.section
        : {
            update: [],
            add: [],
            delete: []
          }
    const binPathDiff = application?.bin_path
      ? application.bin_path
      : { update: [], add: [], delete: [] }
    const schedulerParamDiff = application?.scheduler_param
      ? application.scheduler_param
      : { update: [], add: [], delete: [] }
    const queueDiff = application?.queues
      ? application.queues
      : { update: [], add: [], delete: [] }
    const licenseDiff = application?.licenses
    ? application.licenses
    : { update: [], add: [], delete: [] }
    
    return (
      <Wrapper>
        {diff ? (
          <>
            <InfoDiff application={application} iconData={iconData} />

            <FormDiff formDiff={formDiff} />
            <BinPathDiff binPathDiff={binPathDiff} title='可执行文件路径信息' />
            <BinPathDiff binPathDiff={schedulerParamDiff} title='' />
            <ArrayDiff title='队列信息' arrayDiff={queueDiff} />
            <ArrayDiff title='许可证类型' arrayDiff={licenseDiff} />

            {scriptData && (
              <StyledSection>
                <div>
                  <span className='tag'>脚本信息</span>
                </div>
                <Button
                  title={'点击显示脚本变更详情'}
                  type='link'
                  onClick={() => {
                    Modal.show({
                      title: '脚本变更详情',
                      width: 1200,
                      footer: null,
                      content: ({ onCancel }) => (
                        <ReactDiffViewer
                          styles={this.scriptDiffStyle}
                          oldValue={script.old || ''}
                          newValue={script.new || ''}
                          splitView={true}
                          renderContent={this.highlightSyntax}
                        />
                      )
                    })
                  }}>
                  脚本变更
                </Button>
              </StyledSection>
            )}

            {application.help_doc && (
              <StyledSection>
                <div>
                  <span className='tag'>说明文档</span>
                </div>
                <InlineDiff name='说明文档变更' />
              </StyledSection>
            )}
          </>
        ) : (
          <div>暂无变更</div>
        )}
      </Wrapper>
    )
  }
}
