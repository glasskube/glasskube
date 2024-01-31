# Client

The client is an executable binary written in Go. This is the component the user primarily interacts with. It accepts user inputs via a Web UI and CLI. 
Both kinds of user interface are supposed to provide the same functionality, and it's up to the user to choose their favorite way of doing things. 

The client's tasks revolve around handling user input and consequently managing the required operations with the cluster (via the Kubernetes API) and the package repository. 
For example, when installing a package, the client creates a `Package` custom resource via the Kubernetes API. When listing the available packages, it fetches the package list from the repository.

## CLI

The CLI commands make use of the [Cobra Library](https://github.com/spf13/cobra), which helps in crafting easy to use command line applications.

## GUI

The GUI itself is contained in the `glasskube serve` command, which spins up a local webserver.
For the technical preview we decided to render the pages server side with Go templates. The web technology stack might change in future versions.
