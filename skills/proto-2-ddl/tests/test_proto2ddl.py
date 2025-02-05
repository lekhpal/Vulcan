"""
Tests for proto2ddl utility
"""

import os
import pytest
from pathlib import Path
from click.testing import CliRunner
from src.proto2ddl import main, read_proto_file, generate_schema

# Test data
SAMPLE_PROTO = """
message User {
    required string name = 1;
    required int32 age = 2;
    optional string email = 3;
}

message Post {
    required int64 id = 1;
    required string title = 2;
    required string content = 3;
    optional int64 user_id = 4;
}
"""

@pytest.fixture
def proto_file(tmp_path):
    """Create a temporary proto file for testing."""
    proto_path = tmp_path / "test.proto"
    proto_path.write_text(SAMPLE_PROTO)
    return str(proto_path)

def test_read_proto_file(proto_file):
    """Test reading and parsing proto file."""
    descriptor = read_proto_file(proto_file)
    assert descriptor is not None
    assert len(descriptor.message_type) == 2
    assert descriptor.message_type[0].name == "User"
    assert descriptor.message_type[1].name == "Post"

def test_generate_schema(proto_file):
    """Test schema generation."""
    descriptor = read_proto_file(proto_file)
    schema = generate_schema(descriptor)
    
    # Check for essential parts of the schema
    assert "CREATE TABLE `User`" in schema
    assert "CREATE TABLE `Post`" in schema
    assert "SET NAMES utf8mb4" in schema
    assert "ENGINE=InnoDB" in schema

def test_cli_with_missing_proto_file():
    """Test CLI behavior with missing proto file."""
    runner = CliRunner()
    result = runner.invoke(main, ['--proto-file', 'nonexistent.proto'])
    assert result.exit_code != 0

def test_cli_with_valid_input(proto_file, tmp_path):
    """Test CLI with valid input."""
    schema_file = str(tmp_path / "output.sql")
    runner = CliRunner()
    result = runner.invoke(main, [
        '--proto-file', proto_file,
        '--schema-file', schema_file
    ])
    
    assert result.exit_code == 0
    assert os.path.exists(schema_file)
    
    with open(schema_file) as f:
        schema = f.read()
        assert "CREATE TABLE `User`" in schema
        assert "CREATE TABLE `Post`" in schema

def test_cli_with_default_schema_file(proto_file):
    """Test CLI with default schema file name."""
    runner = CliRunner()
    result = runner.invoke(main, ['--proto-file', proto_file])
    
    expected_schema = str(Path(proto_file).with_suffix('.sql'))
    assert result.exit_code == 0
    assert os.path.exists(expected_schema)
    
    # Cleanup
    os.remove(expected_schema)
