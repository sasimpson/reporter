# CSP Reporter

This is an example of using a Content Security Policy (and Report Only) Header, the Reporting To mechanism, and how to receive that report and log it. 

## Purpose

These are now required by PCI 4.0 for any page handling payment data, but is an interesting mechanism to detect things being injected on your page. The CSP header helps to mitigate cross-site scripting and packet sniffing attacks by limiting and specifying source and protocols. 

## Content-Security-Policy/Content-Security-Policy-Report-Only

Much of the documentation can be found [here](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP).

This allows broad to very fine-grained control of what is allowed to run within your page. 

A normal CSP HTTP Header looks like:

    Content-Security-Policy: <policy here>

### Writing a policy

a full list of the valid policy headers is [here](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy).

example: 

    Content-Security-Policy: default-src: 'site'; script-src: https://cdn.jquery.com; 

`default-src` [🔗](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/default-src) is the fallback for other directives.  at the minimum this should be defined.  

### Additional controls:

* [child-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/child-src)
* [connect-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/connect-src)
* [font-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/font-src)
* [frame-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/frame-src)
* [img-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/img-src)
* [manifest-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/manifest-src)
* [media-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/media-src)
* [object-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/object-src)
* [prefetch-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/prefetch-src)
* [script-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/script-src)
* [script-src-elem](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/script-src-elem)
* [script-src-attr](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/script-src-attr)
* [style-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/style-src)
* [style-src-elem](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/style-src-elem)
* [style-src-attr](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/style-src-attr)
* [worker-src](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/worker-src)

## Reporting-To and reporting violations

Using the `Reporting-To` header, violations of the CSP or CSP-Report-Only can be sent to a waiting endpoint that can capture the content from the browser. 

In order to use this, add the `Reporting-To` header to the HTTP header:

    Reporting-To: csp="http://localhost:8080/csp"

And similarly add the report-uri field to your CSP/CSPRO header:

    Content-Security-Policy: default-src: 'site'; script-src: https://cdn.jquery.com; report-uri=csp

This will generate a report for the CSP violation and POST that JSON to the `csp` endpoint which is `http://localhost:8080/csp` in this case.  

**note** - this requires a secure disposition to use, so https must be used for the reporting *except* when using `localhost` or `127.0.0.1`.

It is possible to have multiple `Reporting-To` endpoints, they need to be separated with a comma(,).

## Tips

Use a `Content-Security-Policy-Report-Only` intially and report the violations to see what parts of your policy might not pass and could break your app.  Additionally you can set the `Content-Security-Policy` to something less strict, report to one endpoint, and make the `Content-Security-Policy-Only` to a stricter set of policies to get data about who is doing what. 
