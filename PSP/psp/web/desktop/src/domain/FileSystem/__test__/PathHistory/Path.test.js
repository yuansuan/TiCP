import Path from "@/domain/FileSystem/PathHistory/Path";

describe("create path", () => {
  test("create Path with source and path params", () => {
    const path = new Path({
      source: "source",
      path: "path"
    });

    expect(path.source).toBe("source");
    expect(path.path).toBe("path");
  });
});

describe("update path", () => {
  test("update path", () => {
    const path = new Path({
      source: "source",
      path: "path"
    });

    path.updatePath("newPath");

    expect(path.path).toBe("newPath");
  });
});
