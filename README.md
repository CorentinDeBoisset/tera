# 🌲 Tera

Boost your development workflow by at least `10^12`.

Tera is a simple CLI tool that help you manage interdependent services, and can run series of tasks without hassle.
This tool was created specifically to simplify the management of local development servers, or the launching of test routines, but you can use it for many other use cases.

It was made possible by the TUI libraries built by the people at [Charmbracelet](https://charm.land/).

## ❯ Installation steps

You can use a pre-built binary with:

```bash
wget "https://github.com/CorentinDeBoisset/tera/releases/latest/download/tera_$(uname)_$(uname -m)" -O tera

# On linux, you can install on /usr/local/bin (requires sudo), or in ~/bin
install -m 0755 tera /path/to/install
```

Alternatively, you can run the following command (you will need to install golang v1.26+):

```bash
go install github.com/corentindeboisset/tera
```

## ❯ Usage

You can save a configuration in a `.tera.yml` file at the root of your project. The configuration is composed of two main parts, the job list (for `tera run <job-name>`) an the service list (for `tera services`).

Here is an example configuration:

```yaml
jobs:
  run-unittests:
    steps:
      - name: Python unittests
        run_before:
          - name: Setup the database
            cmd: ./prepare_db.sh

        # Within a step, all tasks are run simultaneously
        tasks:
          - name: Unittests with coverage
            cmd: coverage run -m pytest
            path: ./path/to/tests

          - name: Lint
            cmd: ruff check && ruff format . --check

        run_after:
          - name: Export the coverage report
            cmd: coverage report

      # You can add as many steps as you want

services:
  backend: API service
    cmd: fastapi run main.py
    path: ./path/to/api
    healthcheck:
      port: 8000 # The service is considered "up" once this port is available

  frontend: Frontend service
    cmd: node run serve-app
    path: ./path/to/frontend
    healthcheck:
      port: 4200
    open_target: "http://localhost:4200" # This link will be opened in a browser once the service is up
    dependencies:
        - target: backend
```

For extensibility, you can add a `.tera.extra.yml` in the same folder as the `.tera.yml` to set up overrides.
If working with git, you can commit the main configuration and put the extra config in the `.gitignore`.

## ❯ Contributing

### Development setup

You can build the development binary:

```bash
make dev
./bin/tera_dev <args>
```

### Tests

If you want to run the tests, you can execute:

```
make test
```

Aditionnaly, if you want a coverage report:

```
make coverage
```

### Translation management

Install the `gotext` excutable:

```bash
go install golang.org/x/text/cmd/gotext@latest
```

Then update the translation catalogs with:

```bash
go generate ./...
```
