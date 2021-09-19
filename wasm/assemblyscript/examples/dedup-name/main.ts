import * as JSON from "../../node_modules/assemblyscript-json/assembly/JSON";

import * as Cache  from "../../lib/cache";
import * as Pubsub from "../../lib/pubsub";
import * as Log from "../../lib/log";

/*
 * Entry point
 */
export function run(): void {
  const event = Pubsub.getMessage()
  Log.println("received event: " + event);
  let jsonObj: JSON.Obj = <JSON.Obj>(JSON.parse(event));
  if (jsonObj != null) {
    let name = jsonObj.getString("name");
    if (name != null) {
      Log.println("name: " + name.toString());
      const lastValue = Cache.get("lastName")
      if (lastValue != name.toString()) {
        Log.println("name value changed, setting new value on cache")
        Cache.set("lastName", name.toString())
        Log.println("sending name forward")
        Pubsub.send(event)
      } else {
        Log.println("repeated name, not sending")
      }
    } else {
      Log.println("no name found in event")
    }
  } else {
    Log.println("no JSON object found, not sending")
  }
}