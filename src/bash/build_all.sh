
#!/bin/bash
gox -ldflags="-w -s" -verbose -arch=amd64 -os="linux darwin windows" -output="build/{{.OS}}_{{.Arch}}/{{.Dir}}" github.com/Byron/godi
