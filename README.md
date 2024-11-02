
# API Testing Tool

This tool automates API testing by reading a JSON configuration file and validating each API endpoint based on the specified request details and expected responses. It’s designed to streamline testing by allowing you to define tests in a structured format.

## Features
- Supports multiple HTTP methods (GET, POST, etc.)
- Validates response status codes and bodies
- Excludes specific fields from validation if needed
- Configurable with JSON for flexible and easy setup

## JSON Configuration Format

The tool reads a JSON file structured as follows:

```json
{
  "tests": {
    "base_url": "http://localhost:8080/api",
    "apis": [
      {
        "skip": true,
        "name": "create user",
        "method": "POST",
        "route": "/users",
        "payload": {
          "name": "kebede",
          "email": "kebede@gmail.com"
        },
        "expects": {
          "status": 200,
          "exclude": ["id", "created_at"],
          "body": {
            "user": {
              "name": "kebede",
              "email": "kebede@gmail.com"
            }
          }
        }
      },
      ...
    ]
  }
}
```

### JSON Fields

- **`base_url`**: The root URL for all API endpoints.
- **`apis`**: An array of test cases, where each test has the following fields:
  - **`skip`** (optional): If `true`, the test will be skipped.
  - **`name`**: A descriptive name for the test.
  - **`method`**: The HTTP method (`GET`, `POST`, etc.).
  - **`route`**: Endpoint path relative to the base URL.
  - **`payload`**: JSON object to send as the request payload for methods like `POST`.
  - **`expects`**:
    - **`status`**: Expected HTTP status code.
    - **`exclude`** (optional): Fields to ignore in the response validation.
    - **`body`**: Expected response body structure.

## How to Use

1. **Prepare the JSON Configuration**:
   - Define all test cases in a JSON file, specifying each API’s details and expected results.

2. **Run the Tool**:
   - Execute the testing tool with the JSON configuration as input.
   - The tool will send requests to each endpoint and validate responses based on the specified expectations.

3. **Review Results**:
   - The tool provides feedback on each test, highlighting any mismatches in status code or response structure.

## Example Usage

To run the tool, you can use a command like:

```bash
go run .\cmd/main.go .\test.json
```

Replace `tests.json` with the path to your JSON configuration file.

## Output

The tool outputs a report indicating:
- Passed tests with matching responses
- Failed tests with details of mismatches

## Contributions

Contributions are welcome! Feel free to submit issues or pull requests to improve the tool’s functionality.
