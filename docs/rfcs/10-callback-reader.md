- Author: Xuanwo <github@xuanwo.io>
- Start Date: 2020-02-09
- RFC PR: N/A
- Tracking Issue: N/A

# Proposal: Callback Reader

## Background

We do need progress report capability:

- For user interaction, we need progress bar to inform them current state
- For service interaction, we need to push/pull progress state

In either case, we all need the ability to do progress report.

## Proposal

So I propose following changes:

- Add Callback Reader/ReadCloser in `pkg/iowrap`
- Add `ReadCallbackFunc` in `types/pairs`
- Add `ReadCallbackFunc` pair for any method who has IO operations
- Wrap with `ReadCallbackFunc` (if `Has`) in every services

After all these work, we can work well with progress bar:

```go
import (
	"io/ioutil"
	"log"

	"github.com/schollz/progressbar/v2"

	"github.com/Xuanwo/storage/coreutils"
	ps "github.com/Xuanwo/storage/types/pairs"
)


func main()  {
 	// Init a service.
 	_, store, err := coreutils.Open("fs:///?work_dir=/tmp")
 	if err != nil {
 		log.Fatalf("service init failed: %v", err)
 	}
 
 	bar := progressbar.New(16 * 1024 * 1024)
 	defer bar.Finish()

 	r, err := store.Read("test_file", ps.WithReadCallbackFunc(func(b []byte) {
 		bar.Add(len(b))
 	}))
 	if err != nil {
 		log.Fatalf("service read failed: %v", err)
 	}
 
 	_, err = ioutil.ReadAll(r)
 	if err != nil {
 		log.Fatalf("ioutil read failed: %v", err)
 	}
}
```

## Rationale

### Builtin progress bar vs callback

Should we provide a builtin progress bar? No, I don't think so.

Progress bar is much more than print lines end with `/r`. We should take care of windows width, characters width, bar theme/style/look, thread safety, multi line support, multi platform support and so on. As a storage lib, we need to focus on storage level, and don't touch those work.

For instead, we need to provide a mechanism which every progress bar can work well with. As far as I know, callback in Read() may be a good choice.

### Callback func(int) vs func([]byte)

There are two options for callbacks in Read: `func(int)` and `func([]byte)`.

Their benchmarks don't have many differences:

```go
// Fist run
goos: linux
goarch: amd64
pkg: github.com/Xuanwo/storage/docs/design/10
BenchmarkPlainReader-8           	  269910	      4149 ns/op	 987.17 MB/s
BenchmarkIntCallbackReader-8     	  292818	      4188 ns/op	 977.98 MB/s
BenchmarkBytesCallbackReader-8   	  265878	      4176 ns/op	 980.90 MB/s
PASS

// Second run
goos: linux
goarch: amd64
pkg: github.com/Xuanwo/storage/docs/design/10
BenchmarkPlainReader-8           	  244312	      4456 ns/op	 919.22 MB/s
BenchmarkIntCallbackReader-8     	  262990	      4202 ns/op	 974.84 MB/s
BenchmarkBytesCallbackReader-8   	  240216	      4290 ns/op	 954.84 MB/s
PASS
```

The results are unstable, compare to read operation, the extra func call doesn't consume too much CPU time.

On this basis, return `n` lacks of scalability, so I choose to return `[]byte` instead.

## Compatibility

No breaking changes.

## Implementation

Most of the work would be done by the author of this proposal.
