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
4. Create a configuration file called config.yaml and set the following options:
    ```makefile
      image_url_file: image_urls.txt
      download_directory: ./downloads/
      batch_size: 2
      min_wait_time: 0.8
      max_wait_time: 3.0
      max_image_size_mb: MAX
      replace_downloaded_file_size: true
      skip_if_file_exists: false
    ```
Adjust the values according to your requirements. Ensure that the image_url_file option points to the correct file name and path. The max_image_size_mb option can be set to "MAX" to skip the size check and download all images regardless of their size.

5. Run it using `go run main.go`. The program will download the images from the URLs specified in the image_urls.txt file. The specified batch size and wait time between batches will be respected. Images exceeding the specified size (if max_image_size_mb is not set to "MAX") will be skipped.

6. The downloaded images will be saved in the downloads directory.


## Configuration Options
- image_url_file: The path to the file containing the list of image URLs to download.
- download_directory: The directory where the downloaded images will be saved.
- batch_size: The number of images to download concurrently in each batch.
- min_wait_time: The minimum wait time between batches (in seconds).
- max_wait_time: The maximum wait time between batches (in seconds).
- max_image_size_mb: The maximum allowed size (in megabytes) for an image. Set to "MAX" to skip the size check and download all images regardless of their size.
- replace_downloaded_file_size: Set it to true to replace already downloaded files if their size differs from the newly downloaded ones. Set it to false to keep the existing files without replacement.
- skip_if_file_exists: Set it to true to skip downloading if the file already exists. Set it to false to allow downloading even if the file exists.
