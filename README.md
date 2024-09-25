# Newsteam SDK

The Newsteam SDK helps developers integrate Newsteam APIs and import content from other sources.

## Wire Feeds Integration

Use the SDK to bring in wire feeds from any source. You can build custom extensions to import these feeds into Newsteam.

### How to Set It Up

1. In Newsteam Desk, go to your feed and click "Wire" to configure.

2. On your computer:

    - Make a new directory
    - Run `go mod init app`
    - Install the SDK: `go get github.com/feight/newsteam-sdk`

3. Create an importer with these methods:

    - `Id()`: Returns a string to identify your feed
    - `GetLogfiles()`: Gets log files as byte slices
    - `ProcessLogfile([]byte)`: Turns each log file into an `*admin.ArticleInput`

4. In `main.go`, set up wire feeds like this:

```go
newsteam.InitializeFeeds([]newsteam.Feed{
    &cosmos.CosmosImporter{
        Feed: "bd", Host: "https://businesslive.co.za/apiv1",
    },
})
```

That's it! You're ready to start importing wire feeds.
