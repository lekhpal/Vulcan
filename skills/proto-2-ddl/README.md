# Proto to DDL Converter

A utility to convert Google Protocol Buffer definitions to MySQL schema definitions.

## Requirements

- Python 3.12 or higher
- Dependencies listed in requirements.txt

## Installation

```bash
pip install -r requirements.txt
```

## Usage

```bash
python src/proto2ddl.py --proto-file input.proto --schema-file output.sql
```

### Arguments

- `--proto-file` (required): Input Protocol Buffer definition file
- `--schema-file` (optional): Output MySQL schema file. If not provided, will use the proto filename with .sql extension

## Development

### Running Tests

```bash
pytest tests/
```

### Logging

The utility uses loguru for logging. Logs are written to:
- Console (INFO level and above)
- `proto2ddl.log` file (DEBUG level and above)

## Error Handling

The utility will exit with a non-zero status code and appropriate error message if:
- Required proto file is not provided
- Proto file cannot be read or parsed
- Schema file cannot be written
