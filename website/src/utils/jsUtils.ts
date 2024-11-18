// Inspired by https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore?tab=readme-ov-file#_sortby-and-_orderby
export function sortBy<T>(
  array: T[],
  getter: (item: T) => string | number | boolean,
): T[] {
  const sortedArray = [...array];
  sortedArray.sort((a, b) =>
    getter(a) > getter(b) ? 1 : getter(b) > getter(a) ? -1 : 0,
  );
  return sortedArray;
}

export function toggleListItem<T>(list: T[], item: T): T[] {
  const itemIndex = list.indexOf(item);
  if (itemIndex === -1) {
    return list.concat(item);
  }
  const newList = [...list];
  newList.splice(itemIndex, 1);
  return newList;
}
