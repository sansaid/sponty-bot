# sponty-bot

A humble little bot with a God complex that randomly generates a "night out" plan for you, if you can't be bothered to think of one. Locally used by a Reading Discord server, so this isn't for general use unless you're in Reading, UK.

## Usage

To ask Sponty to generate a night out for you, just use the command:

```sh
/rng-party
```

There are optional options that can be used to further customise your event:

* `generate_chaplin` (`bool`, default: `true`) - whether to randomly assign a "party chaplin" (i.e. someone that ensures the event goes as planned). This will randomly pick someone in the `Party Chaplin` role on Discord.
* `location_type` ([`park`, `pub`], default: `pub`) - the location type to meet up in. The party perks change per location.