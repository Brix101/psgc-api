# PSGC TOOL / PSGC API

PSGC (Philippine Standard Geographic Code) Tool is a Go application for managing geographic data in the Philippines. This command-line tool provides functionality for running a RESTful API and generating JSON files from CSV input. All the data used in this app is sourced from
**[ Philippine Statistics Authority](https://psada.psa.gov.ph/psgc).**

> **Desctiption:** I created this tool as a side project to address the need for efficient data management. Many existing APIs were slow and lacked pagination features, so this tool was developed to provide a **faster and more robust solution** for data management.

**API URL:** [PSGC- API](https://psgc-api.onrender.com) (Redirects to Swagger Documentation)

## Table of Contents

- [API Documentation](#api-documentation)
- [Usage](#usage)
  - [Usage with Air](#usage-with-air)
  - [Building](#building)
  - [Running the RESTful API](#running-the-restful-api)
  - [Running the Json Generator](#running-the-json-generator)
- [Options](#options)
  - [Common Options](#common-options)
  - [API Command Options](#api-command-options)
  - [Generator Command Options](#generator-command-options)

## API Documentation

The PSGC API is documented using Swagger Docs, providing detailed information about available endpoints, request parameters, and example responses. You can explore the API documentation by following the link below:

[Swagger API Documentation](https://psgc-api.onrender.com/docs/index.html)

> **Note:** Please refer to the Swagger documentation for a comprehensive guide on how to use the API effectively.

## Usage

### Usage with Air

When developing PSGC, you can use Go's Air tool for hot-reloading during development to run the API, which helps streamline the development process. First, ensure that you have Go installed on your system.

To run the API with Air, follow these steps:

Start the API with Air:

```bash
air api
```

Start the API with Makefile:

```bash
make dev
```

### Building

Before using PSGC, you need to build the executable binary. Make sure you have Go installed on your system.

To build PSGC, use the following command:

```bash
make build
```

### Running the RESTful API

To run the PSGC RESTful API, use the following command:

```bash
./psgc api
```

You can specify a different port using the `--port` option (see [API Command Options](#api-command-options)).

> **Note:** The API's default port is set to 5000. If a port is specified in an environment file (e.g., .env), that port will take precedence as the default. However, you can also use the `--port` flag when running the program, and it will override both the default port and the value specified in the environment file.

### Running the Json Generator

To generate json files from csv, use the following command:

```bash
./psgc generate
```

You can specify a different port using the `--port` option (see [Generator Command Options](#generator-command-options)).

> **Note:** By default, the generator will use the default file located at `files/csv/psgc.csv`.

## Options

### Common Options

- `--profile, -p`: Record CPU profiling data. This option can be used with any PSGC command to enable CPU profiling.

### API Command Options

- `--port, -P`: Specify the port on which the API will run (default is 5000).

### Generator Command Options

- `--file, -f`: Specify the path to the CSV input file for JSON generation. If not provided, the generator will use the default file located at `files/csv/psgc.csv`.
