// Package integtest provides integration tests for the monotask CLI.
//
// The tests use binary of the application. In case of missing binary they skip.
//
// Binary for tests must be provided using: MONOTASK_BINARY:
//
//   - specify binary manually.
//   - use dorenv files.
//   - use make goal
package integtest
