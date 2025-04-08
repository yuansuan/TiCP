[<img src="https://mengzhuo.github.io/yaml-to-go/resources/images/yaml-to-go.png" alt="YAML-to-Go converts YAML to a Go struct"></a>](https://mengzhuo.github.io/yaml-to-go)

Translates YAML into a Go type definition. [Check it out!](http://mengzhuo.github.io/yaml-to-go)

This is a sister tool to:

- [curl-to-Go](https://mholt.github.io/curl-to-go), which converts curl commands to Go code.
- [json-to-Go](https://mholt.github.io/json-to-go), which converts JSON to Go code.

Things to note:

- The script sometimes has to make some assumptions, so give the output a once-over.
- In an array of objects, it is assumed that the first object is representative of the rest of them.
- The output is indented, but not formatted. Use `go fmt`!

Contributions are welcome! Open a pull request to fix a bug, or open an issue to discuss a new feature or change.


### Credits

YAML-to-Go is brought to you by Meng Zhuo ([mengzhuo](https://twitter.com/mengzhuo)).

Fork from JSON-to-Go by Matt Holt ([mholt6](https://twitter.com/mholt6)).

The Go Gopher is originally by Renee French. This artwork is an adaptation.
