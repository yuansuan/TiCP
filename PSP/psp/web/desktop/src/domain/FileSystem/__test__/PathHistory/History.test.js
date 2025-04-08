import History from "@/domain/FileSystem/PathHistory/History";

beforeEach(() => {
  jest.resetModules();
});

describe("create History", () => {
  test("create History with no initialValue", () => {
    const history = new History();

    expect(history.list.length).toBe(0);
    expect(history.cursor).toBe(-1);
  });

  test("create History with initialValue", () => {
    const history = new History("test");

    expect(history.list.length).toBe(1);
    expect(history.list[0]).toBe("test");
    expect(history.cursor).toBe(0);
  });
});

describe("get current from history", () => {
  test("when list is empty, current is undefined", () => {
    const history = new History();

    expect(history.current).toBeUndefined();
  });

  test("when list is't empty, current depends on cursor", () => {
    const history = new History("test");

    expect(history.cursor).toBe(0);
    expect(history.current).toBe("test");

    history.cursor = 1;
    expect(history.cursor).toBe(1);
    expect(history.current).toBeUndefined();
  });
});

describe("judge that if history can go back", () => {
  test("if cursor is less than 0, history.prevDisabled is true", () => {
    const history = new History();

    expect(history.cursor).toBe(-1);
    expect(history.prevDisabled).toBe(true);
  });

  test("if cursor is more than 0, history.prevDisabled is false", () => {
    const history = new History();

    history.cursor = 1;
    expect(history.cursor).toBe(1);
    expect(history.prevDisabled).toBe(false);
  });
});

describe("judge that if history can go forward", () => {
  test("if cursor is more than history.list.length - 1, history.nextDisabled is true", () => {
    const history = new History();

    expect(history.cursor >= history.list.length - 1).toBe(true);
    expect(history.nextDisabled).toBe(true);
  });

  test("if cursor is less than history.list.length - 1, history.nextDisabled is false", () => {
    const history = new History("test");

    history.cursor = -1;
    expect(history.cursor >= history.list.length - 1).toBe(false);
    expect(history.nextDisabled).toBe(false);
  });
});

describe("push item to history", () => {
  test("if push item in middle of history, item will replace the remain", () => {
    const history = new History();

    history.list = [0, 1, 2, 3];
    history.cursor = 0;
    history.push(1);
    expect(history.cursor).toBe(1);
    expect(history.list).toEqual([0, 1]);
  });

  test("push item to the tail", () => {
    const history = new History(0);

    history.push(1);
    expect(history.cursor).toBe(1);
    expect(history.list).toEqual([0, 1]);
  });
});

describe("go back to prev item", () => {
  test("call history.prev when history.prevDisabled is true will be omiited and return undefined", () => {
    const history = new History();

    expect(history.prevDisabled).toBe(true);
    expect(history.cursor).toBe(-1);
    const prevItem = history.prev();
    expect(prevItem).toBe(undefined);
    expect(history.cursor).toBe(-1);
  });

  test("call history.prev when history.prevDisabled is false will go back and return prev item", () => {
    const history = new History();

    history.push(0);
    history.push(1);

    expect(history.prevDisabled).toBe(false);
    expect(history.cursor).toBe(1);
    const prevItem = history.prev();
    expect(prevItem).toBe(0);
    expect(history.cursor).toBe(0);
  });
});

describe("go forward to next item", () => {
  test("call history.next when history.nextDisabled is true will be omiited and return undefined", () => {
    const history = new History();

    expect(history.nextDisabled).toBe(true);
    expect(history.cursor).toBe(-1);
    const nextItem = history.next();
    expect(nextItem).toBe(undefined);
    expect(history.cursor).toBe(-1);
  });

  test("call history.next when history.nextDisabled is false will go forward and return next item", () => {
    const history = new History();

    history.push(0);
    history.push(1);
    history.prev();

    expect(history.nextDisabled).toBe(false);
    expect(history.cursor).toBe(0);
    const nextItem = history.next();
    expect(nextItem).toBe(1);
    expect(history.cursor).toBe(1);
  });
});

describe("iterate history", () => {
  test("iterate history", () => {
    const history = new History();

    const list = [1, 2, 3, 4];
    history.list = list;
    expect([...history]).toEqual([...list]);
  });
});
