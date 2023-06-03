# advanced-image-downloader-golang

A Go program to download images from a list of URLs in parallel, with configurable batch size and tsnfom wait timed between batches.

## Prerequisites

- Go (version 1.16 or higher)

## Getting Started

1. Clone the repository:

   ```shell
   git clone https://github.com/your-username/image-downloader.git
   ```
2. Navigate to the project directory:

    ```shell
    cd image-downloader
    ```

3. Create a text file called image_urls.txt and add the URLs of the images you want to download, each on a separate line.
4. Create a configuration file called config.txt and set the following options:
    ```makefile
    ImageURLFile = image_urls.txt
    DownloadDirectory = ./downloads/
    BatchSize = 2
    MinWaitTime = 0.8
    MaxWaitTime = 3.0
    ```
Adjust the values according to your requirements. Ensure that the ImageURLFile option points to the correct file name and path.

5. Run it using `go run main.go`
6. The downloaded images will be saved in the downloads directory.


## Configuration Options
- ImageURLFile: The path to the file containing the list of image URLs to download.
- DownloadDirectory: The directory where the downloaded images will be saved.
- BatchSize: The number of images to download concurrently in each batch.
- MinWaitTime: The minimum wait time between batches (in seconds).
- MaxWaitTime: The maximum wait time between batches (in seconds).
