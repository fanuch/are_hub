# Connecting client procedure (protocol)
1. The channel requested is found by the channel ID in the URL parameters passed in when upgrading to a websocket. I.e. `/<publish or subscribe>/<id>`. If no channel is found, then the client is disconnected with a WS_ERROR_NOT_FOUND status code.

2. The client is then prompted for the password to the channel they would like to connect to by sending a WS_CHALLENGE_PASSWORD status. If the passwords do not match the client is disconnected with a WS_ERROR_UNAUTHORISED status code.

3. If the telemetry channel is live and the connecting client is a subscriber (i.e. upgrading from `/subscribe/<id>`), then the client is added to the telemetry channel and will start receiving forwarded messages from the publisher. If the telemetry channel is live and the connecting client is a publisher (i.e. upgrading from `/publish/<id>`) and there is no currently connected publisher, then the client will be set as the new publisher.
