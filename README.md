# Vacuum 🧹

Vacuum is a powerful command-line tool designed to help you efficiently move old files from a source directory (`source`) to a designated target directory (`target`). It offers various customizable options to suit your archiving needs.

---

## 🚀 Features

- **File Archiving**: Seamlessly moves files from your source to target directories.
- **Customizable Parameters**: Define file age, subdirectory inclusion, and more.
- **Dry Run Option**: Safely preview actions before executing them.
- **Optional Shredding**: Securely delete original files post-transfer if required.

---

## 🔧 Basic Usage

```bash
vacuum -source=<source directory> -target=<target directory> [options]
```

Example:

```bash
vacuum -source=/home/user/documents -target=/archive/documents -age=5 -shred
```

---

## 📚 Command-Line Options

| Flag          | Description                                                        |
| ------------- | ------------------------------------------------------------------ |
| `-source`     | Path to the source directory from which files will be moved.       |
| `-target`     | Path to the target directory where files will be archived.         |
| `-dry`        | Perform a dry run without executing any file operations.           |
| `-shred`      | Delete the original file after copying it to the target directory. |
| `-recurse`    | Recursively include all subdirectories.                            |
| `-age`        | Minimum file age in years to consider for archiving.               |

---

## 📝 Example Commands

1. **Dry Run**: Simulate the process without moving any files.

   ```bash
   vacuum -source=/path/to/source -target=/path/to/target -dry
   ```

2. **Archiving Files Older than 5 Years**:

   ```bash
   vacuum -source=/path/to/source -target=/path/to/target -age=5
   ```

3. **Move and Shred**: Archive and securely delete original files.

   ```bash
   vacuum -source=/path/to/source -target=/path/to/target -age=3 -shred
   ```

4. **Recursively Include Subdirectories**:

   ```bash
   vacuum -source=/path/to/source -target=/path/to/target -recurse
   ```

---

## ⚠️ Important Notes

- **Dry Run**: Highly recommended to use the `-dry` flag before performing actual file operations to ensure everything works as expected.
- **File Shredding**: If the `-shred` flag is used, files will be permanently deleted from the source directory after transfer. This action cannot be undone.
- **Recursive Option**: The `-recurse` flag enables processing of subdirectories.

---

## 📄 License

This project is licensed under the BSD-3-Clause License. See the LICENSE file for more details.

---

### 🌟 Contributing

Feel free to submit issues or pull requests. Contributions are always welcome!
