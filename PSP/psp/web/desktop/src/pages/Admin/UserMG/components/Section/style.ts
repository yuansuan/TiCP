import styled from 'styled-components'

export const SectionWrapper = styled.div`
  .header {
    display: flex;
    align-items: center;
    font-weight: bold;
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
    margin-top: 30px;
    margin-bottom: 20px;
    line-height: 22px;

    .icon {
      margin-right: 10px;
      line-height: 16px;
    }
  }

  .body {
    display: flex;
    flex-direction: column;
  }
`
