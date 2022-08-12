# lexos
A port of the lexos-cli tool in a Go package, for ease of use in Go projects.

### API: <br/>
```Get```: Takes a ISBN(string) as input and returns the Lexile(```int```), Atos(```float64```), Ar(```float64```), and error if any.<br/>
```Install```: Installs the necessary driver and browser. Called in ```Get```.
