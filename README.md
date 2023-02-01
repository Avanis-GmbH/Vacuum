# Go Dust Vacuum

A tool to move old files from a source directory (root-dir) into a target directory (target-dir).

Basic usage:
```
go-dust-vacuum -root-dir=<root directory> -target-dir=<target directory> [additional flags]
```

Additional flags:
```
  -dry
        Should the application perform a dry run without any io operations?
  -help
        Shows usage information about this software
  -nolog
        If no log files should be written for the process. Use at own risk only!
  -min-age int
        How old the last edit of a file should be (in years) to consider it for archiving. (default 11)
  -r    Should all subdirectories be included. (default true)
  -shred
        Should the original file get deleted after copy?
```
