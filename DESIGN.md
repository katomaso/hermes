# Hermes design

Hermes has multiple relays for messages.

In the first draft implementation (that might be as well final) the agents are
compiled together in the big binary and messages are handled by go-routines.
Therefore you cannot find list of services in the database even though they
should be there. You will not see usage of an external message queue because
all is build in.