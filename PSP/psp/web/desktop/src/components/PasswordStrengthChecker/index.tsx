import * as React from 'react'
import styled from 'styled-components'

const StrengthMap = {
  weak: { label: '强度低', color: '#D65350' },
  medium: { label: '强度中', color: '#E8B35D' },
  strong: { label: '强度高', color: '#3B9C7B' },
}

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  .strengthLines {
    display: flex;
    width: 100%;
    justify-content: space-between;
    .strengthLine {
      opacity: 0.4;
      width: 30%;
      .line {
        border-bottom: 4px solid gray;
        border-radius: 3px;
      }
      .text {
        font-size: 12px;
        padding: 5px 20px;
        text-align: center;
      }
    }
  }
  .tips {
    font-size: 12px;
    color: #333;
  }
`

function PasswordStrengthChecker({ strength, style, tips }) {
  return (
    <Wrapper>
      <div className='strengthLines' style={style}>
        {Object.keys(StrengthMap).map(key => (
          <div
            className='strengthLine'
            key={key}
            style={{
              color: StrengthMap[key].color,
              opacity: strength === key ? 1 : 0.4,
            }}>
            <div
              className='line'
              style={{ borderBottomColor: StrengthMap[key].color }}></div>
            <div className='text'>{StrengthMap[key].label}</div>
          </div>
        ))}
      </div>
      <div className='tips'>提示: {tips}</div>
    </Wrapper>
  )
}

export default PasswordStrengthChecker
