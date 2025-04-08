import FavoriteDirectory from "@/domain/FileSystem/Points/Favorite/FavoriteDirectory";

describe("create favoriteDirectory", () => {
  test("create favoriteDirectory will generate favoriteId", () => {
    const props = {
      path: "path",
      name: "name"
    };
    const directory = new FavoriteDirectory(props);

    expect(directory.favoriteId).toBe(
      window.btoa(`${props.name}::${props.path}`)
    );
  });
});
