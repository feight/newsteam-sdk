# Newsteam SDK

The Newsteam SDK helps developers integrate Newsteam APIs and import content from other sources.

## Wire Integration

Use the SDK to bring in articles from any source. You can build custom extensions to import these articles into Newsteam.

### How to Set It Up

1. In Newsteam Desk, go to your bucket and click "Wire" to configure.

2. On your computer:

    - Make a new directory
    - Run `go mod init app`
    - Install the SDK: `go get github.com/feight/newsteam-sdk`

3. Create an importer with these methods:

    - `Id()`: Returns a string to identify your bucket
    - `GetLogfiles()`: Gets log files as byte slices
    - `ProcessLogfile([]byte)`: Turns each log file into an `*admin.Article`

4. In `main.go`, set up a wire like this:

```go
newsteam.InitializeBuckets([]newsteam.Bucket{
    &MyImporter{},
})
```

That's it! You're ready to start importing articles.
