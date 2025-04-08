import RootPoint from "@/domain/FileSystem/Points/RootPoint";
import Store from "@/domain/FileSystem/Store";
import { Http } from "@/utils";

beforeEach(() => {
  jest.resetModules();
});

describe("create homePoint", () => {
  test("create homePoint with root path", () => {
    Store.hooks.afterDelete.tapAsync = jest.fn(Store.hooks.afterDelete.tapAsync);
    Store.hooks.afterUpdate.tapAsync = jest.fn(Store.hooks.afterUpdate.tapAsync);
    Store.hooks.afterAdd.tapAsync = jest.fn(Store.hooks.afterAdd.tapAsync);

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    expect(homePoint.rootPath).toBe("rootPath");
    expect(Store.hooks.afterDelete.tapAsync).toBeCalled();
    expect(Store.hooks.afterUpdate.tapAsync).toBeCalled();
    expect(Store.hooks.afterAdd.tapAsync).toBeCalled();
  });
});

describe("homePoint.onAfterDelete", () => {
  test("homePoint.onAfterDelete", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.add({
      path: "test"
    });

    expect(
      homePoint.filterFirstNode((item) => item.path === "test")
    ).toBeTruthy();
    homePoint.onAfterDelete(["test"]);
    expect(
      homePoint.filterFirstNode((item) => item.path === "test")
    ).toBeUndefined();
  });
});

describe("homePoint.onAfterUpdate", () => {
  test("homePoint.onAfterUpdate update node is not in homePoint will be omitted", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.init(homePoint, [
      {
        path: "rootPath/test",
        name: "name"
      }
    ]);

    const node = homePoint.filterFirstNode(
      (item) => item.path === "rootPath/test"
    );
    expect(node).toBeDefined();
    homePoint.onAfterUpdate([
      {
        path: "invalidPath",
        newProps: {
          path: "invalidPath",
          name: "newName"
        }
      }
    ]);
    expect(node.name).not.toBe("newName");
  });

  test("homePoint.onAfterUpdate update childNode", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.init(homePoint, [
      {
        path: "rootPath/test",
        name: "name"
      }
    ]);

    const node = homePoint.filterFirstNode(
      (item) => item.path === "rootPath/test"
    );
    expect(node).toBeDefined();
    homePoint.onAfterUpdate([
      {
        path: "rootPath/test",
        newProps: {
          path: "rootPath/test",
          name: "newName"
        }
      }
    ]);
    expect(node.name).toBe("newName");
  });
});

describe("homePoint.onAfterAdd", () => {
  test("homePoint.onAfterAdd add node is not in homePoint will be omitted", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.onAfterAdd([
      {
        path: "testPath",
        newProps: {
          path: "testPath",
          name: "newName"
        }
      }
    ]);
    expect(homePoint.children.length).toBe(0);
  });

  test("homePoint.onAfterUpdate upadddate childNode", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.onAfterAdd([
      {
        path: "rootPath/test",
        newProps: {
          path: "rootPath/test",
          name: "name"
        }
      }
    ]);
    expect(homePoint.children.length).toBe(1);
  });
});

describe("homePoint.init", () => {
  test("Mount childNodes to homePoint's node", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.init(homePoint, [
      {
        path: "rootPath/node",
        name: "node"
      }
    ]);

    expect(homePoint.children.length).toBe(1);
  });

  test("Mount node which is not belong to homePoint's node will be omitted", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.init(homePoint, [
      {
        path: "testPath/node",
        name: "node"
      }
    ]);

    expect(homePoint.children.length).toBe(0);
  });
});

describe("homePoint.service.point", () => {
  test("homePoint.service.point is homePoint", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    expect(homePoint.service.point).toBe(homePoint);
  });
});

describe("homePoint.service.fetch", () => {
  const files = [
    {
      path: "rootPath/node1",
      name: "node1"
    },
    {
      path: "rootPath/node2",
      name: "node2"
    },
    {
      path: "rootPath/node3",
      name: "node3"
    }
  ];

  beforeAll(() => {
    Store.delete = jest.fn(Store.delete);
    Store.update = jest.fn(Store.update);
  });

  afterAll(() => {
    Store.delete.mockClear();
    Store.update.mockClear();
  });

  test("homePoint.service.fetch and update nodes", async () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    Http.get = jest.fn(() =>
      Promise.resolve({
        data: files
      })
    );

    await homePoint.service.fetch("rootPath");
    expect(Http.get).toBeCalledWith("/web/data/file/list", {
      params: { path: "rootPath", columns: "NAME,TYPE,SIZE,PATH,MDATE" }
    });
    expect(Store.delete).not.toBeCalled();
    expect(Store.update).toBeCalledWith(files);
  });

  test("homePoint.service.fetch and delete nodes", async () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    Http.get = jest.fn(() =>
      Promise.resolve({
        data: files
      })
    );

    homePoint.init(homePoint, [
      {
        path: "rootPath/node4",
        name: "node4"
      }
    ]);

    await homePoint.service.fetch("rootPath");
    expect(Http.get).toBeCalledWith("/web/data/file/list", {
      params: { path: "rootPath", columns: "NAME,TYPE,SIZE,PATH,MDATE" }
    });
    expect(Store.delete).toBeCalledWith(["rootPath/node4"]);
    expect(Store.update).toBeCalledWith(files);
  });
});

describe("homePoint.service.fetchPaths", () => {
  beforeAll(() => {
    Store.delete = jest.fn(Store.delete);
    Store.update = jest.fn(Store.update);
  });

  afterAll(() => {
    Store.delete.mockClear();
    Store.update.mockClear();
  });

  test("fetchPaths return empty data will be omitted", async () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    Http.get = jest.fn(() =>
      Promise.resolve({
        data: []
      })
    );

    await homePoint.service.fetchPaths("rootPath/node1");

    expect(Http.get).toBeCalledWith("/web/data/path/list", {
      params: { path: "rootPath/node1", with_parents: true }
    });
    expect(Store.update).not.toBeCalled();
  });

  test("fetchPaths update nodes", async () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    Http.get = jest.fn(() =>
      Promise.resolve({
        data: [
          {
            path: "rootPath",
            name: "root",
            sub_file: [
              {
                path: "rootPath/node1",
                name: "node1"
              },
              {
                path: "rootPath/node2",
                name: "node2"
              },
              {
                path: "rootPath/node3",
                name: "node3"
              }
            ]
          }
        ]
      })
    );

    await homePoint.service.fetchPaths("rootPath/node1");

    expect(Store.update).toBeCalledWith([
      {
        path: "rootPath",
        name: "root"
      },
      {
        path: "rootPath/node1",
        name: "node1"
      },
      {
        path: "rootPath/node2",
        name: "node2"
      },
      {
        path: "rootPath/node3",
        name: "node3"
      }
    ]);
  });
});

describe("homePoint.service.delete", () => {
  test("homePoint will invoke fetch after delete", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);

    const params = {
      path: "testPath",
      names: []
    };
    await homePoint.service.delete(params);

    expect(Http.post).toBeCalledWith("/web/data/file/delete", params);
    expect(homePoint.service.fetch).toBeCalled();
  });
});

describe("homePoint.service.rename", () => {
  beforeAll(() => {
    Store.rename = jest.fn(Store.rename);
  });

  afterAll(() => {
    Store.rename.mockClear();
  });

  test("homePoint.service.rename with invalid name will be rejected", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);

    const params = {
      path: "testPath",
      name: "name",
      newName: "\\newName"
    };

    try {
      await homePoint.service.rename(params);
    } catch (err) {
      expect(err).toBeDefined();
    }

    expect(Http.post).not.toBeCalled();
  });

  test("homePoint.service.rename with valid name", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);

    const params = {
      path: "testPath",
      name: "name",
      newName: "newName"
    };

    await homePoint.service.rename(params);

    expect(Http.post).toBeCalledWith("/web/data/file/rename", {
      name: params.name,
      new_name: params.newName,
      path: params.path
    });

    expect(Store.rename).toBeCalledWith(
      `${params.path}/${params.name}`,
      params.newName
    );
  });
});

describe("homePoint.service.createDir", () => {
  test("homePoint.service.createDir with invalid name will be rejected", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      name: "\\newName"
    };

    try {
      await homePoint.service.createDir(params);
    } catch (err) {
      expect(err).toBeDefined();
    }

    expect(Http.post).not.toBeCalled();
  });

  test("homePoint.service.createDir with valid name will invoke fetch", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      name: "newName"
    };

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);
    await homePoint.service.createDir(params);

    expect(Http.post).toBeCalledWith("/web/data/file/create_dir", params);
    expect(homePoint.service.fetch).toBeCalled();
  });
});

describe("homePoint.service.move", () => {
  test("homePoint.service.move with valid name will invoke fetch", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      names: ["node1"],
      toPath: "rootPath/test"
    };

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);
    await homePoint.service.move(params);

    expect(Http.post).toBeCalledWith("/web/data/file/move", {
      path: params.path,
      names: params.names,
      to_path: params.toPath
    });
    expect(homePoint.service.fetch).toBeCalledTimes(2);
    expect(homePoint.service.fetch).toHaveBeenNthCalledWith(1, params.path);
    expect(homePoint.service.fetch).toHaveBeenNthCalledWith(2, params.toPath);
  });
});

describe("homePoint.service.copy", () => {
  test("homePoint.service.copy with valid name will invoke fetch", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      names: ["node1"],
      toPath: "rootPath/test"
    };

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);
    await homePoint.service.copy(params);

    expect(Http.post).toBeCalledWith("/web/data/file/copy", {
      path: params.path,
      names: params.names,
      to_path: params.toPath
    });
    expect(homePoint.service.fetch).toHaveBeenCalledWith(params.toPath);
  });
});

describe("homePoint.service.download", () => {
  test("homePoint.service.download", () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });

    document.createElement = jest.fn(document.createElement);

    homePoint.service.download({
      path: "rootPath",
      names: ["node1", "node2"]
    });
    expect(document.createElement).toBeCalledWith("a");
  });
});

describe("homePoint.service.compress", () => {
  test("homePoint.service.compress with valid name will invoke fetch", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      names: ["node1"],
      compressType: "zip",
      zipName: "test.zip"
    };

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);
    await homePoint.service.compress(params);

    expect(Http.post).toBeCalledWith("/web/data/file/compress", {
      path: params.path,
      names: params.names,
      compress_type: params.compressType,
      zip_name: params.zipName
    });
    expect(homePoint.service.fetch).toHaveBeenCalledWith(params.path);
  });

  test("homePoint.service.compress with invalid name will be rejected", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = {
      path: "rootPath",
      names: ["node1"],
      compressType: "zip",
      zipName: "\\test.zip"
    };

    try {
      await homePoint.service.compress(params);
    } catch (err) {
      expect(err).toBeDefined();
    }

    expect(Http.post).not.toBeCalled();
  });
});

describe("homePoint.service.extract", () => {
  test("homePoint.service.extract with valid name will invoke fetch", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = { path: "rootPath", toPath: "rootPath/to", name: "node1" };

    homePoint.service.fetch = jest.fn(homePoint.service.fetch);
    await homePoint.service.extract(params);

    expect(Http.post).toBeCalledWith("/web/data/file/extract", {
      path: params.path,
      name: params.name,
      to_path: params.toPath
    });
    expect(homePoint.service.fetch).toHaveBeenCalledWith(params.toPath);
  });

  test("homePoint.service.extract with invalid name will be rejected", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const homePoint = new RootPoint({
      path: "rootPath"
    });

    const params = { path: "rootPath", toPath: "rootPath/to", name: "\\node1" };

    try {
      await homePoint.service.extract(params);
    } catch (err) {
      expect(err).toBeDefined();
    }

    expect(Http.post).not.toBeCalled();
  });
});

describe("homePoint.view", () => {
  test("homePoint.view", async () => {
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    const res = {
      data: {
        Content: "content"
      }
    };
    Http.get = jest.fn(() => Promise.resolve(res));

    const params = { path: "rootPath", name: "name", lastLines: 0 };
    const viewContent = await homePoint.service.view(params);
    expect(Http.get).toBeCalledWith("/web/data/file/view", {
      params: {
        path: params.path,
        name: params.name,
        last_lines: params.lastLines
      }
    });
    expect(viewContent).toBe(res.data.Content);
  });
});

describe("homePoint.save", () => {
  test("homePoint.save", async () => {
    Http.post = jest.fn(() => Promise.resolve());

    const params = { path: "rootPath", name: "name", content: "content" };
    const homePoint = new RootPoint({
      path: "rootPath"
    });
    await homePoint.service.save(params);

    expect(Http.post).toBeCalledWith("/web/data/file/save", {
      path: params.path,
      name: params.name,
      content: params.content
    });
  });
});
