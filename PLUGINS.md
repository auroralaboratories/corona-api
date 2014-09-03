# Plugin Development Primer


Plugins are built in two phases; one for implementation of functionality (think of these as models), and one for doing the necessary work to return a useful REST response (e.g.: controllers).

## Models

Models are where the functional implementation of the tasks a plugin performs are built.  By convention, these are in files prefixed with `plugin-` (e.g.: `plugin-session.go`)

Plugins must do the following things:

* Implement the `IPlugin` inteface
** Define an _Init()_ method to perform necessary setup/connection
* Mix-in the `BasePlugin` struct
* Instantiate the plugin instance in the `SprinklesAPI.Plugins` map (in `web.go`'s _Init()_ method)

