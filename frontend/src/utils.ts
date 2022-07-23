export const basename = (path: string) => {
  return path.split(/[\\/]/).pop() || "";
};

export const dirname = (path: string) => {
  return path.replace(basename(path), "");
};
