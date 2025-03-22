# Gator - RSS Feed Aggregator CLI

## Description

Gator is a guided project from [Boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang)

> Build a blog aggregator microservice in Go. Put your API, database, and web scraping skills to the test.

## Requirements

1. Go installed
2. Postgres

## Setup

1. Copy .env template file and populate with db details
2. run `goose postgres <dburl> migrate up` from `sql\schema` directory
3. Create a `.gatorconfig.json` file in `$HOME` file should contain:

```JSON
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```

## Usage

`gator <command> <arg>`

## Commands
