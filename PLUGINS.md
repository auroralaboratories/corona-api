# Plugin Development Primer


Plugins are built in two phases; one for implementation of functionality (think of these as models), and one for doing the necessary work to return a useful REST response (i.e.: controllers).

## Models

Models are where the functional implementation of the tasks a plugin performs are built.  By convention, these are in files prefixed with "`plugin-`" (e.g.: `plugin-session.go`)

Models must do the following things:

* Implement the `IPlugin` inteface
  * Define an _Init()_ method to perform necessary setup/connection
* Mix-in the `BasePlugin` struct
* Instantiate the plugin instance in the `CoronaAPI.Plugins` map (in `web.go`'s _Init()_ method)


## Controllers

Controllers are relatively thin, almost boilerplate, methods that are called by the REST routing engine.  By convention, these are in files prefixed with "`api-`" (e.g.: `api-session-workspaces.go`).  These methods are implemented on `CoronaAPI` and are responsible for calling the plugin's model implementations, using data passed in via POST data, query strings, or other per-call and environmental factors.

Controllers must do the following things:

* Perform discrete operations that represent logical API endpoints
* Output structured data via `go-json-rest`'s _WriteJson()_ writer method, or errors via the _rest.Error()_ method.
  * Very occasionally, other types of data will be returned, but the convention is to try to represent state via JSON as much as possible.
  * 


## Routes

Routes are the actual URL paths that map to your controllers.  At present, they are defined in `web.go`'s _Init()_ method.


# See Also

* https://godoc.org/github.com/ant0ine/go-json-rest/rest
