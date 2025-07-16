# Instagram Scraper

A simple and efficient Instagram scraper built with Go. This tool allows users to programmatically extract data from Instagram profiles, including posts, followers, and more.

## Features

- Scrape user profiles and gather insights
- Fetch posts, comments, and likes
- Option to set proxies
- Easy-to-use command-line interface


## Contact

For questions or suggestions, contact here: [Telegram](https://t.me/shiny0103) | [Twitter](https://x.com/0xTan1319)


## Installation

To get started, clone the repository to your local machine:

```bash
git clone https://github.com/0xTan1319/instagram-scraper.git
cd instagram-scraper
```
## Prerequisites

Go installed on your machine (Version 1.18 or later)

## Usage

Set up your environment:

Create a .env file based on .env.example and configure your parameters.

Run the scraper:

You can run the application using the command:

```bash

go run main.go
Command-line options:
```

Use the following command-line options to customize your scraping experience:

```text

-user [USERNAME]       Specify the Instagram username
-proxy [PROXY]        Use a proxy to avoid IP bans
```

## Example

To scrape data from a user, execute:

```bash

go run main.go -user example_username
```

## Proxy Configuration

You can set up proxies for scraping by adding proxy addresses in the proxies.txt file. Each proxy should be listed on a new line.

Contributing
Contributions are welcome! If you have suggestions for improvements or new features, feel free to fork the repository and create a pull request.

