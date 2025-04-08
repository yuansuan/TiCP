import React from 'react'
import styled from 'styled-components'

const StyledDiff = styled.div`
  margin-bottom: 30px;

  &.inline {
    display: flex;

    .name {
      margin-right: 17px;
      width: 100px;
      text-align: right;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .name {
    font-family: 'PingFangSC-Medium';
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
  }

  .mask {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    z-index: 2;
  }

  .diff {
    display: flex;
    flex: 1;

    .item {
      width: 600px;

      > div {
        width: 600px;
      }
    }

    .old {
      flex: 1;
      position: relative;
      padding: 4px 8px;
      margin-right: 12px;
      min-width: 36px;
      text-align: center;
      word-break: break-all;

      .mask {
        background-color: rgba(250, 122, 122, 0.4);

        &::after {
          content: '';
          position: absolute;
          border-top: 1px solid ${props => props.theme.borderColor};
          width: 100%;
          left: 0;
          top: 50%;
        }
      }
    }

    .new {
      flex: 1;
      position: relative;
      padding: 4px 8px;
      min-width: 36px;
      text-align: center;
      word-break: break-all;

      .mask {
        background-color: rgba(27, 218, 148, 0.1);
      }
    }
  }
`

export const PathBlockDiff = ({ className, name, Old, New }) => {
  return (
    <StyledDiff className={className}>
      {name && (
        <div className='name' title={name}>
          {name}
        </div>
      )}
      <div className='diff'>
        {Old && (
          <div className='old'>
            <div className='mask' />
            {Old}
          </div>
        )}
        {New && (
          <div className='new'>
            <div className='mask' />
            {New}
          </div>
        )}
      </div>
    </StyledDiff>
  )
}

export const InlinePathDiff = props => (
  <PathBlockDiff {...props} className={'inline'} />
)
