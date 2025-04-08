import { FileStore } from "@/domain/FileSystem/Store";

beforeEach(() => {
  jest.resetModules();
});

describe("create FileStore", () => {
  test("init hooks when create FileStore", () => {
    const store = new FileStore();

    expect(store.hooks).toBeTruthy();
    expect(store.hooks.afterAdd).toBeTruthy();
    expect(store.hooks.afterDelete).toBeTruthy();
    expect(store.hooks.afterUpdate).toBeTruthy();
  });
});

describe("get node from store", () => {
  test("fileStore.get", () => {
    const store = new FileStore();

    const node = {
      path: "path"
    };
    store.nodeMap.set("testNode", node);
    expect(store.get("testNode")).toEqual(node);
  });
});

describe("delete node from store", () => {
  test("fileStore.delete with string", () => {
    const store = new FileStore();

    store.nodeMap.set("path", { path: "path" });
    expect(store.nodeMap.size).toBe(1);

    // mock nodeMap's delete
    store.nodeMap.delete = jest.fn(store.nodeMap.delete);

    store.delete("path");
    expect(store.nodeMap.size).toBe(0);
    expect(store.nodeMap.delete).toBeCalledWith("path");
  });

  test("fileStore.delete with string[]", () => {
    const store = new FileStore();

    store.nodeMap.set("path1", { path: "path" });
    store.nodeMap.set("path2", { path: "path" });
    expect(store.nodeMap.size).toBe(2);

    // mock nodeMap's delete
    store.nodeMap.delete = jest.fn(store.nodeMap.delete);

    store.delete(["path1", "path2"]);
    expect(store.nodeMap.size).toBe(0);
    expect(store.nodeMap.delete).toBeCalledTimes(2);
  });
});

describe("harmony node's struct", () => {
  test("fileStore.harmony", () => {
    const store = new FileStore();

    expect(
      store.harmony({
        mdate: "modified time"
      })
    ).toEqual({
      modifiedTime: "modified time"
    });
  });
});

describe("update multiple nodes", () => {
  test("call fileStore.update with null will throw error ", () => {
    const store = new FileStore();

    expect(() => store.update(null)).toThrow(Error);
  });

  test("fileStore.update will new node when node is't existed", () => {
    const store = new FileStore();

    store.harmony = jest.fn(store.harmony);

    const node = { path: "path", name: "name" };
    store.update([node]);

    expect(store.harmony).toBeCalled();
    expect(store.nodeMap.size).toBe(1);
    expect(store.nodeMap.get(node.path)).toEqual(store.harmony(node));
  });

  test("fileStore.update will update node when node is existed", () => {
    const store = new FileStore();

    const node = { path: "path", name: "name" };
    store.nodeMap.set(node.path, node);
    expect(store.nodeMap.size).toBe(1);
    expect(store.nodeMap.get(node.path)).toEqual(node);

    const newNode = {
      path: "path",
      name: "newName"
    };
    store.update([newNode]);

    expect(store.nodeMap.size).toBe(1);
    expect(store.nodeMap.get(node.path)).toEqual(newNode);
  });
});

describe("rename specific node", () => {
  test("fileStore.rename will omit update when node is't existed", () => {
    const store = new FileStore();

    store.harmony = jest.fn(store.harmony);

    const node = { path: "path", name: "name" };
    store.rename(node.path, "newName");

    expect(store.harmony).not.toBeCalled();
    expect(store.nodeMap.size).toBe(0);
  });

  test("fileStore.rename update existed node", () => {
    const store = new FileStore();

    const node = { path: "path/name", name: "name" };
    store.nodeMap.set(node.path, node);

    store.rename(node.path, "newName");

    expect(store.nodeMap.size).toBe(1);
    expect(store.nodeMap.get(node.path)).toBeUndefined();
    expect(store.nodeMap.get("path/newName")).toEqual({
      path: "path/newName",
      name: "newName"
    });
  });

  test("fileStore.rename delete child nodes", () => {
    const store = new FileStore();

    const node = { path: "path/name", name: "name" };
    const childNode = { path: "path/name/child", name: "name" };
    store.nodeMap.set(node.path, node);
    store.nodeMap.set(childNode.path, childNode);

    store.rename("path/name", "newName");

    expect(store.nodeMap.size).toBe(1);
    expect(store.nodeMap.get(childNode.path)).toBeUndefined();
  });
});
