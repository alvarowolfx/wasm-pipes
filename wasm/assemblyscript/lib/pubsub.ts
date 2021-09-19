
@external("pubsub", "getMessage")
declare function _getMessage(ptr: usize): void;

@external("pubsub", "send")
declare function _send(ptr: usize, len: i32): void;

export function getMessage(): string {
  const ptr = memory.data(1024);
  _getMessage(ptr);
  return String.UTF8.decodeUnsafe(ptr, 1024, true);
}

export function send(str: string): void {
  let out = String.UTF8.encode(str);
  _send(changetype<usize>(out), out.byteLength);
}