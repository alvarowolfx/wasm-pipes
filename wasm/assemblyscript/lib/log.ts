@external("log", "println")
declare function _println(ptr: usize, len: i32): void;

export function println(str: string): void {
  let out = String.UTF8.encode(str);
  _println(changetype<usize>(out), out.byteLength);
}