import React, { useState } from 'react';
import { Tree, Input } from 'antd';

const { TreeNode } = Tree;
const { Search } = Input;

const OrganizationTree = ({ treeData, onCheck,refresh }) => {
  const [expandedKeys, setExpandedKeys] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [checkedKeys, setCheckedKeys] = useState([]);

  const renderTreeNodes = (nodes) => {
    return nodes.map(node => {
      if (node.children) {
        return (
          <TreeNode title={node.name} key={node.key}>
            {renderTreeNodes(node.children)}
          </TreeNode>
        );
      }
      return <TreeNode title={node.name} key={node.key} />;
    });
  };
  const handleSearch =async (value) => {
    const lowerSearchText = value.toLowerCase(); // 转换搜索关键词为小写
    setSearchText(value);
    await refresh(lowerSearchText);

  };

  const handleCheck = (keys) => {
    setCheckedKeys(keys);
    onCheck(keys);
  };

  return (
    <div>
      <Search
        placeholder="搜索人员"
        onChange={(e) => handleSearch(e.target.value)}
        value={searchText}
      />
      <Tree
        checkable
        onCheck={handleCheck}
        expandedKeys={expandedKeys}
        onExpand={(keys) => setExpandedKeys(keys)}
        checkedKeys={checkedKeys}
        treeData={treeData}
      />
    </div>
  );
};

export default OrganizationTree;
