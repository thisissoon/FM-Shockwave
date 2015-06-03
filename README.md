# FM Shockwave

SOON\_ FM Volume Management Service. This service subscribes to a Redis
Pub/Sub service listening for volume change events. Once an event has been
recieved it changes the volume and publishes a volume chnaged event back
to the channel.

Note: Currently this can only run on Linux.
