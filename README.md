# whatsup
Is your site or service up? if not, create a pagerduty alert. Resolves
alerts when things come back up.

## Installation

With `go get`:

```
go install github.com/r0fls/whatsup
```

From source:

```
git clone https://github.com/r0fls/whatsup.git
cd whatsup
go install
```


## Usage:

You'll need to have a pagerduty account, and create a service or use an
existing one. To find your services, go to your pagerduty account, then go to:

Configuration > Services

Once your there click, on the gear and click view, or create a new service with
the green button. When viewing your service, click the Integrations tab. Finally, you'll see the integration key.


```
./whatsup -k <YOUR_SERVICE_INTEGRATION_KEY> -r http://localhost:8000 -p 30s
```

the `--period` or `-p` flag defaults to 60s.
