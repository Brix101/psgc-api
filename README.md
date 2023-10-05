# PSGC - Philippine Standard Geographic Code TOOL

PSGC (Philippine Standard Geographic Code) is a Go application for managing geographic data in the Philippines. This command-line tool provides functionality for running a RESTful API and generating JSON files from CSV input.

**Desctiption:** I created this API as a side project to address the need for efficient geographic data management in the Philippines. Many existing APIs were slow and lacked pagination features, so PSGC was developed to provide a **faster and more robust solution** for geographic data management.

**The CSV data used in PSGC app is sourced from [https://psada.psa.gov.ph/psgc](https://psada.psa.gov.ph/psgc) and is converted from XLSX format.**

**API URL:** [https://psgc-api.onrender.com](https://psgc-api.onrender.com) (Redirects to Swagger Documentation)

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

> **Note:** The API primarily uses the port specified in the environment file. If a port is defined in the environment file, that port will be used as the default, and the `--port` flag or the default port will not be used.

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
