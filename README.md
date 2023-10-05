# PSGC - Philippine Standard Geographic Code

PSGC (Philippine Standard Geographic Code) is a Go application for managing geographic data in the Philippines. This command-line tool provides functionality for running a RESTful API and generating JSON files from CSV input.

## Table of Contents

- [Usage](#usage)
  - [Running the RESTful API](#running-the-restful-api)
  - [Generating JSON Files](#generating-json-files)
- [Options](#options)
  - [API Command Options](#api-command-options)
  - [Generator Command Options](#generator-command-options)
- [Profiling](#profiling)
- [Contributing](#contributing)
- [License](#license)

## Usage

### Building

Before using PSGC, you need to build the executable binary. Make sure you have Go installed on your system.

To build PSGC, use the following command:

```bash
go build -o psgc ./cmd/http
```

### Running the RESTful API

To run the PSGC RESTful API, use the following command:

```bash
./psgc api
```

By default, the API will run on port 5000. You can specify a different port using the `--port` option (see [API Command Options](#api-command-options)).

### Running the Json Generator

To generate json files from csv, use the following command:

```bash
./psgc generate
```

By default, the generator will use the default file located at `files/csv/psgc.csv`.You can specify a different port using the `--port` option (see [Generator Command Options](#generator-command-options)).

## Options

### Common Options

- `--profile, -p`: Record CPU profiling data. This option can be used with any PSGC command to enable CPU profiling.

### API Command Options

- `--port, -P`: Specify the port on which the API will run (default is 5000).

### Generator Command Options

- `--file, -f`: Specify the path to the CSV input file for JSON generation. If not provided, the generator will use the default file located at `files/csv/psgc.csv`.
