# Directory Synchronization Tool

A command-line tool to efficiently synchronize the contents of a source directory with a destination directory. The tool ensures files are copied, updated, or optionally removed in the destination directory to match the source directory's structure and contents.

---

## Features

- **One-Time Synchronization**:
  - Copies files from the source directory to the destination directory if they are missing.
  - Updates files in the destination directory if the size or modification date differs from the source.
  - Optionally deletes extra files in the destination directory, ensuring a mirror of the source's state.

- **Error Handling**:
  - Gracefully logs any file access errors or permission issues without halting the program.
  - Ensures critical validation errors (e.g., invalid source path) stop the application with an appropriate error message.

- **CLI Arguments**:
  - Specify source and destination directories via command-line flags.
  - Control deletion of extra files in the destination using an optional flag.

- **Logging**:
  - Logs operations such as file copying, updating, and deletion, helping users keep track of changes made during synchronization.

---

## Prerequisites

- **Go Language**: Version 1.18 or later should be installed on your system to build and run the tool.

---

## Usage

### Command-Line Arguments

| Argument           | Description                                                                                     |
|--------------------|-------------------------------------------------------------------------------------------------|
| `-src`             | (Required) Specifies the path to the source folder. This directory must exist.                  |
| `-dst`             | (Required) Specifies the path to the destination folder. If the folder does not exist, it will be created. |
| `-delete-missing`  | (Optional) Deletes files in the destination directory that are missing in the source directory. Defaults to `false`. |

### Example Commands

1. **Basic Synchronization**
   Synchronize files from `srcDir` to `dstDir`, without deleting any extra files in `dstDir`:
   ```bash
   dirsync -src=/path/to/source -dst=/path/to/destination
   ```

2. **Synchronization with File Deletion**
   Synchronize files from `srcDir` to `dstDir` and delete any files in `dstDir` that are not present in `srcDir`:
   ```bash
   dirsync -src=/path/to/source -dst=/path/to/destination -delete-missing
   ```

---

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd dirsync
   ```

2. Build the executable:
   ```bash
   go build -o dirsync
   ```

3. Run the tool using the built executable:
   ```bash
   ./dirsync -src=/path/to/source -dst=/path/to/destination
   ```

---

## Logging

The tool logs all file operations such as file copying, overwriting, and deletion to the console. Errors encountered during the synchronization process are also logged.

---

## Error Handling

- **Non-Critical Errors**:
  - Issues such as missing files or permission problems in individual files are logged but do not terminate the program.

- **Critical Errors**:
  - Missing or invalid source directory paths will terminate the program with an appropriate error message.
