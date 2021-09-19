import * as JSON from "../../node_modules/assemblyscript-json/assembly/JSON";

import * as Pubsub from "../../lib/pubsub";
import * as Log from "../../lib/log";

const threshold: f64 = 20.0;
/*
 * Entry point
 */
export function run(): void {
  const event = Pubsub.getMessage()
  Log.println("received event: " + event);
  let jsonObj: JSON.Obj = <JSON.Obj>(JSON.parse(event));
  if (jsonObj != null) {
    let temp = jsonObj.getFloat("temp");
    if (temp != null) {
      Log.println("temp: " + temp.toString());
      if (temp.valueOf() <= threshold) {
        Log.println("temperature below threshold, forwarding event")
        Pubsub.send(event)
      } else {
        Log.println("temperature above threshold, not forwarding event")
      }
    } else {
      Log.println("no temperature found in event")
    }
  } else {
    Log.println("no JSON object found, not sending")
  }
}