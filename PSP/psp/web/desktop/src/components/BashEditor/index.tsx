import * as React from 'react';
import { observer } from 'mobx-react-lite'

import {CodeEditor} from '@/components'
interface IProps {
  code?: string;
  placeholder?: string;
  width?: string;
  height?: string;
  onChange?: (code: string) => void;
  readOnly?: boolean; // 新增只读属性
}

let  editorRef = null
export const BashEditor = observer((props: IProps) => {
  const [editingCode, setEditingCode] = React.useState('');

  React.useEffect(() => {
    setEditingCode(props.code)
  },[props.code])
  
  const handleCodeChange = () => {
    const newCode = editorRef && editorRef.getValue();
    setEditingCode(newCode);

    if (props.onChange) {
      props.onChange(newCode);
    }
  };

  
  return (

    <div style={{ 
      padding: '20px 0',
      width: props.width || '550px',
      height: props.height || '200px',
      maxHeight: '600px',
      overflow: 'auto',
      border: '1px solid rgba(0,0,0,0.1)',
      backgroundColor: props.readOnly ? '#f8f8f8' : 'white',
    }}>
    <CodeEditor
      ref={ref => (editorRef = ref)}
      value={editingCode}
      language='shell'
      readOnly={props.readOnly}
      onChange={handleCodeChange}
    />
  </div>
  );
});

