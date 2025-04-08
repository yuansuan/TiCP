import PathHistory from "@/domain/FileSystem/PathHistory";
import Path from "@/domain/FileSystem/PathHistory/Path";
import Store from "@/domain/FileSystem/Store";

beforeEach(() => {
  jest.resetModules();
});

describe("create pathHistory", () => {
  test("subscribe Store after create pathHistory", () => {
    Store.hooks.afterUpdate.tapAsync = jest.fn(Store.hooks.afterUpdate.tapAsync);
    new PathHistory();

    expect(Store.hooks.afterUpdate.tapAsync).toBeCalled();
  });
});

describe("onStoreUpdate", () => {
  test("onStoreUpdate will omit extraneous event", () => {
    const history = new PathHistory();

    history.push({
      path: "home/name"
    });

    history.onStoreUpdate([
      {
        path: "home2",
        oldProps: {
          name: "home2"
        },
        newProps: {
          name: "newHome"
        }
      }
    ]);

    expect(history.current.path).not.toBe("newHome/name");
  });

  test("onStoreUpdate name", () => {
    const history = new PathHistory();

    history.push({
      path: "home/name"
    });

    history.onStoreUpdate([
      {
        path: "home",
        oldProps: {
          name: "home"
        },
        newProps: {
          name: "newHome"
        }
      }
    ]);

    expect(history.current.path).toBe("newHome/name");
  });
});

describe("get current path from history", () => {
  test("pathHistory.currentPath will return empty string when pathHistory.list is empty", () => {
    const history = new PathHistory();

    expect(history.list.length).toBe(0);
    expect(history.currentPath).toBe("");
  });

  test("pathHistory.currentPath will return current.path", () => {
    const history = new PathHistory();

    const item = {
      path: "testPath"
    };
    history.push(item);
    expect(history.currentPath).toBe(item.path);
  });
});

describe("push item to pathHistory", () => {
  test("push item to pathHistory will create Path instance", () => {
    const history = new PathHistory();

    const params = {
      path: "path",
      source: "source"
    };
    const item = history.push(params);
    expect(item instanceof Path).toBe(true);
    expect(item).toMatchObject(params);
  });

  test("if the current is same with the new Path, it will be omitted", () => {
    const history = new PathHistory();

    const params = {
      path: "path",
      source: "source"
    };
    history.push(params);
    const item = history.push(params);
    expect(item).toBeNull();
  });
});
