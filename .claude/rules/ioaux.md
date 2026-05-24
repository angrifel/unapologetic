---
globs: ["ioaux/**"]
---

## ReadSeekBuffer

`ReadSeekBuffer` combines `bytes.Buffer` and `bytes.Reader` from the Go standard library: it carries the full read/write buffer behavior of `bytes.Buffer` and adds `Seek`, `ReadAt`, and `Size` from `bytes.Reader`.

**Constraint**: Do not modify any behavior inherited from `bytes.Buffer` or `bytes.Reader`. Changes to `ReadSeekBuffer` are strictly limited to the methods sourced from `bytes.Reader` (`Seek`, `ReadAt`, `Size`) and any future additions in that same vein. If a proposed change would also apply to `bytes.Buffer` or `bytes.Reader` unchanged, it does not belong here.

**Concurrency**: `ReadSeekBuffer` is not safe for concurrent use. This is a known and accepted limitation — do not flag it as a bug or attempt to add synchronization.
