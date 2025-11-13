## Ideas

- [x] implement long polling for Updates and doRequest (Fetch)
- [x] implement exponential backoff
- [ ] improve error handling and bad request cases:
    - [ ] Max consecutive failures
    - [ ] errors classification and program failure in case of really critical erros (badToken or Unauthorized)
    - [ ] [circuit breaker pattern](https://learn.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)) in case of transient faults
- [ ] implement controlled goroutine pool for concurrent event processing.
- [ ] implement SQLlite
- [ ] try something with loging (LOKI or other)
- [x] fix issue with "[ERR] consumer: can't get events: empty updates list"