import List from "@/domain/FileSystem/List";
import Store from "@/domain/FileSystem/Store";

beforeEach(() => {
  jest.resetModules();
});

describe("create fileList", () => {
  test("subscribe store when create fileList", () => {
    Store.hooks.afterDelete.tapAsync = jest.fn(Store.hooks.afterDelete.tapAsync);
    Store.hooks.afterUpdate.tapAsync = jest.fn(Store.hooks.afterUpdate.tapAsync);
    Store.hooks.afterAdd.tapAsync = jest.fn(Store.hooks.afterAdd.tapAsync);

    new List();

    expect(Store.hooks.afterDelete.tapAsync).toBeCalled();
    expect(Store.hooks.afterUpdate.tapAsync).toBeCalled();
    expect(Store.hooks.afterAdd.tapAsync).toBeCalled();
  });
});

describe("fileList.onAfterDelete", () => {
  test("fileList.onAfterDelete", () => {
    const list = new List();
    list.parentPath = "home";

    const nodePath = "home/test";
    list.mount([
      {
        path: nodePath
      }
    ]);
    expect(list.filterFirstNode((item) => item.path === nodePath)).toBeTruthy();
    list.onAfterDelete([nodePath]);
    expect(
      list.filterFirstNode((item) => item.path === nodePath)
    ).toBeUndefined();
  });
});

describe("fileList.onAfterUpdate", () => {
  let list;
  const nodePath = "home/test";

  beforeEach(() => {
    list = new List();
    list.parentPath = "home";

    list.mount([
      {
        path: nodePath,
        name: "name"
      }
    ]);
  });

  test("fileList.onAfterUpdate update child node", () => {
    const oldNode = list.filterFirstNode((item) => item.path === nodePath);
    expect(oldNode).toBeTruthy();
    expect(oldNode.name).toBe("name");

    list.onAfterUpdate([
      {
        path: nodePath,
        newProps: {
          path: nodePath,
          name: "newName"
        }
      }
    ]);

    const newNode = list.filterFirstNode((item) => item.path === nodePath);
    expect(newNode).toBeTruthy();
    expect(newNode.name).toBe("newName");
  });

  test("fileList.onAfterUpdate update node is't in list will be omitted", () => {
    const oldNode = list.filterFirstNode((item) => item.path === nodePath);
    expect(oldNode).toBeTruthy();
    expect(oldNode.name).toBe("name");

    list.onAfterUpdate([
      {
        path: "invalidPath",
        newProps: {
          path: "invalidPath",
          name: "newName"
        }
      }
    ]);

    const newNode = list.filterFirstNode((item) => item.path === nodePath);
    expect(newNode.name).not.toBe("newName");
  });
});

describe("fileList.onAfterAdd", () => {
  test("fileList.onAfterAdd add child nodes", () => {
    const list = new List();
    list.parentPath = "home";

    list.onAfterAdd([
      {
        path: "home/node1"
      },
      {
        path: "home/node2"
      }
    ]);

    expect(list.children.length).toBe(2);
    expect(
      list.filterFirstNode((item) => item.path === "home/node1")
    ).toBeTruthy();
    expect(
      list.filterFirstNode((item) => item.path === "home/node2")
    ).toBeTruthy();
  });

  test("fileList.onAfterAdd add nodes are not belong list will be omitted", () => {
    const list = new List();
    list.parentPath = "home";

    list.onAfterAdd([
      {
        path: "invalidPath/node1"
      },
      {
        path: "invalidPath/node2"
      }
    ]);

    expect(list.children.length).toBe(0);
  });
});

describe("judge a path is fileList's child path", () => {
  test("fileList.isChild", () => {
    const list = new List();

    list.parentPath = "home";
    expect(list.isChild("home/test")).toBe(true);
    expect(list.isChild("test")).toBe(false);
  });
});

describe("mount specific nodes to list", () => {
  test("fileList.mount only mount child nodes", () => {
    const list = new List();

    list.parentPath = "home";

    list.mount([
      {
        path: "home/node1",
        name: "node1",
        isdir: false
      },
      {
        path: "home/node2",
        name: "node2",
        isdir: true
      },
      {
        path: "node3",
        name: "node3"
      }
    ]);

    expect(list.childrenMap.size).toBe(2);
    expect(list.children[0].path).toBe("home/node1");
    expect(list.children[1].path).toBe("home/node2");
  });
});

describe("update list when path changed", () => {
  test("call fileList.update with same path will be omitted", () => {
    const list = new List();

    list.clear = jest.fn(list.clear);
    list.mount = jest.fn(list.mount);

    list.parentPath = "testPath";
    list.update("testPath");

    expect(list.clear).not.toBeCalled();
    expect(list.mount).not.toBeCalled();
  });

  test("call fileList.update with different path", () => {
    const list = new List();

    list.clear = jest.fn(list.clear);
    list.mount = jest.fn(list.mount);

    list.parentPath = "testPath";
    list.update("newPath");

    expect(list.clear).toBeCalled();
    expect(list.mount).toBeCalled();
  });
});
