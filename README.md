# Newsteam SDK

The Newsteam SDK equips developers with a comprehensive toolkit designed to integrate seamlessly with the core Newsteam APIs, facilitating both the incorporation of Newsteam's native functionalities and the importation of content from external sources.

## Wire Feeds Integration

Newsteam SDK empowers developers to integrate wire feeds from any source into the Newsteam platform. By leveraging the SDK, developers can create custom extensions that enable the seamless importation of wire feeds, enhancing the versatility and functionality of the Newsteam environment.

### Implementation Guide

#### Configuring Wire Feeds in Newsteam Desk:

To initiate the process, configure your feed within Newsteam to establish a connection to a wire:

1. Navigate to Newsteam Desk, select your feed, and click on "Wire".

#### Establishing a New Go Module:

For the integration, you will need to establish a Go module:

1. Create an empty directory within your workspace on your local environment.
2. Initialize the Go module by executing `go mod init app` in your terminal.
3. Install the Newsteam SDK with the command `go get github.com/feight/newsteam-sdk`.

Next, implement an importer by defining the following methods:

-   `Id()` returning a string that identifies your feed.
-   `GetLogfiles()` which retrieves log files as slices of byte slices, with a potential error return.
-   `ProcessLogfile([]byte)` that processes each log file into a slice of `*admin.ArticleInput` pointers.

Create a `main.go` file, and initialize the wire feeds by incorporating the following snippet:

```go
newsteam.InitializeFeeds([]newsteam.Feed{
    &cosmos.CosmosImporter{
        Feed: "bd", Host: "https://businesslive.co.za/apiv1",
    },
})
```
