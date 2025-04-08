import Favorite from "@/domain/FileSystem/Points/Favorite";
import { Http } from "@/utils";

describe("create favorite point", () => {
  test("create favorite point will set rootPath", () => {
    const favorite = new Favorite({ path: "rootPath" });

    expect(favorite.rootPath).toBe("rootPath");
  });
});

describe("favorite.fetch", () => {
  test("favorite.fetch success will invoke mount", async () => {
    const res = {
      data: [
        {
          path: "rootPath",
          name: "node1"
        }
      ]
    };
    Http.get = jest.fn(() => Promise.resolve(res));
    const favorite = new Favorite({ path: "rootPath" });
    favorite.mount = jest.fn();

    await favorite.fetch();

    expect(Http.get).toBeCalledWith("/web/data/favorite/list");
    expect(favorite.mount).toBeCalledWith(favorite, res.data);
  });
});

describe("favorite.mount", () => {
  test("favorite.mount", () => {
    const favorite = new Favorite({ path: "rootPath" });
    favorite.clear = jest.fn();
    favorite.mount(favorite, [
      {
        path: "rootPath",
        name: "node1"
      },
      {
        path: "rootPath",
        name: "node2"
      }
    ]);

    expect(favorite.clear).toBeCalled();
    expect(favorite.children.length).toBe(2);
  });
});
