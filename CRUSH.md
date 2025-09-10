# Build/Lint/Test Commands
go build -v -o /dev/null .                    # Build and check compilation
go test -v ./handler                             # Run all handler tests
go test -v ./handler -run TestUV                # Run specific test
GOOS=windows go build -v -o /dev/null .         # Check Windows compatibility

# Code Style Guidelines
- Use early returns, no "else { return <expr> }"
- Prefer methods on structs over package-level functions
- Minimize package-level variables (use singletons when needed)
- Error handling: return errors with context, use fmt.Errorf for wrapping
- Naming: MixedCaps for exported, mixedCaps for unexported
- Imports: group standard, third-party, local packages; no aliases
- Types: Use concrete types unless interface abstraction is needed
