# godeps
golang dependency analyzer

## install
```bash
go install github.com/jcyamacho/godeps@latest
```

## commands

### *flowchart*
generate a dependency graph mermaid flowchart
#### flags:
-  -h, --help, help for flowchart
-  -s, --skip, skip modules
-  -i, --skip-indirect, skip indirect dependencies (default true)
#### args:
- [path], golang module directory to scan
#### example:
```mermaid
flowchart
        github.com/jcyamacho/godeps --> github.com/hashicorp/go-multierror
        github.com/jcyamacho/godeps --> github.com/spf13/cobra
        github.com/jcyamacho/godeps --> golang.org/x/mod
```