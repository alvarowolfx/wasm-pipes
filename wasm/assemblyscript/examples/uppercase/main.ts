
import * as Pubsub from "../../lib/pubsub";
import * as Log from "../../lib/log";

/*
 * Entry point
 */
export function run(): void {
  const event = Pubsub.getMessage()
  Log.println("received event: " + event);
  const uppercased = event.toUpperCase()
  Log.println("uppercased event: " + uppercased);
  Pubsub.send(uppercased)
}