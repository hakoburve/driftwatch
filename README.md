# driftwatch

> CLI tool that detects configuration drift between deployed services and their declared state in version control.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your version-controlled config directory and a running environment to compare against:

```bash
# Check for drift against a live Kubernetes cluster
driftwatch scan --source ./config/prod --env kubernetes --namespace production

# Output drift report as JSON
driftwatch scan --source ./config/prod --env kubernetes --format json

# Watch continuously and alert on changes
driftwatch watch --source ./config/prod --env kubernetes --interval 60s
```

Example output:

```
[DRIFT DETECTED] service/api-gateway
  replicas: declared=3, actual=5
  image tag: declared=v1.4.2, actual=v1.4.1

[OK] service/auth-service
[OK] service/worker
```

---

## Configuration

`driftwatch` can be configured via a `.driftwatch.yaml` file in your project root:

```yaml
source: ./config/prod
environment: kubernetes
namespace: production
format: table
interval: 30s
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)