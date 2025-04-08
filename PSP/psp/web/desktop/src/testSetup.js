/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import '@testing-library/jest-dom'
import Enzyme from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'

Enzyme.configure({ adapter: new Adapter() })
