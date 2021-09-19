@external("cache", "set")
declare function _set(key: usize, key_len: i32, value: usize, value_len: i32): void;

@external("cache", "get")
declare function _get(key: usize, key_len:i32, value: usize): void;

export function get(key: string): string {
  let outKey = String.UTF8.encode(key);
  const ptr = memory.data(128);
  _get(
    changetype<usize>(outKey), outKey.byteLength,
    ptr
  );
  return String.UTF8.decodeUnsafe(ptr, 128, true);
}

export function set(key: string, value: string): void {
  let outKey = String.UTF8.encode(key);
  let outValue = String.UTF8.encode(value);
  _set(
    changetype<usize>(outKey), outKey.byteLength,
    changetype<usize>(outValue), outValue.byteLength
  )
}