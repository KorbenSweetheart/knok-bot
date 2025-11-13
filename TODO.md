## Ideas

- [x] implement long polling for getUpdates (Fetch)
- [ ] implement suggested improvements
- [ ] implement SQLlite
- [ ] try something with loging (LOKI or other)
- [ ] fix loop with "[ERR] consumer: can't get events: empty updates list"
> don’t treat an empty getUpdates result as an error and use long-polling (timeout) + exponential backoff;
> also use a controlled goroutine pool for concurrent processing. The docs show getUpdates is designed for long polling and that updates are stored up to 24 hours;
> you should rely on that instead of spinning when the server returns [].
- [ ] fix potential memory leak