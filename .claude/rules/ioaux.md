---
globs: ["ioaux/**"]
---

## ReadSeekBuffer

`ReadSeekBuffer` combines `bytes.Buffer` and `bytes.Reader` from the Go standard library: it carries the full read/write buffer behavior of `bytes.Buffer` and adds `Seek`, `ReadAt`, and `Size` from `bytes.Reader`, plus `Position`.

**Constraint**: Do not modify the advertised behavior of `bytes.Buffer` or `bytes.Reader`, meaning what their documentation promises. Changes to `ReadSeekBuffer` are strictly limited to the methods sourced from `bytes.Reader` (`Seek`, `ReadAt`, `Size`), `Position`, any future additions in that same vein, and the internal details those methods depend on. If a proposed change would also apply to `bytes.Buffer` or `bytes.Reader` unchanged, it does not belong here.

**Retention**: Internal side effects of `bytes.Buffer` that are not part of its documented contract are not protected. In particular, `bytes.Buffer` drops bytes already read to reuse space (in `grow`, and through `Reset` calls in `Read`, `ReadByte`, `ReadRune` and `WriteTo`). `ReadSeekBuffer` must never drop bytes before the read position, because `Seek` and `ReadAt` can reach them. Only `Reset` and `Truncate` may discard them, because discarding is what those calls are for. The cost is that memory is not reused until one of them is called; this is accepted.

**Concurrency**: `ReadSeekBuffer` is not safe for concurrent use. This is a known and accepted limitation — do not flag it as a bug or attempt to add synchronization.
