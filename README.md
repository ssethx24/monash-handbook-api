# Monash Handbook API

This is a simple API that scrapes and serves Monash University handbook. It is written in Go and uses MongoDB and Redis for storing and caching handbook data.

This was created as a personal project for a course mapper and timetabler. However, the timetabler is now a separate project and the course mapper is not yet implemented.

This repository contains the scraper and API needed to fetch handbook data, which would be useful in creating a course mapper.

The timetabler is available on [timetabler.jasondev.me](https://timetabler.jasondev.me).

## Page Content
- [How it works](#how-it-works)
- [Setup](#setup)
- [API Endpoints](#api-endpoints)
  - [Handbook Data](#handbook-data)
    - [Get Unit Information](#get-unit-information)
    - [Get Course Information](#get-course-information)
    - [Get Area of Study Information](#get-area-of-study-information)
    - [Check Unit Requisites](#check-unit-requisites)
    - [Get Handbook Search API URL](#get-handbook-search-api-url)
  - [Health Check](#health-check)

A simple API, purely written in Go, that scrapes and serves Monash University handbook and timetable data.

## How it works

There are two main parts of this API:
- **Handbook Data:** Scrapes and serves detailed information about units, courses, and areas of study, including prerequisites and enrolment rules checking.

The handbook data can be used with no setup; just run the server and start querying. It depends on `handbook.monash.edu`

## Setup

1. Install Go: https://golang.org/doc/install
2. Install MongoDB: https://docs.mongodb.com/manual/installation/
3. Install Redis: https://redis.io/download
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Set up environment variables, as shown in sample.env
6Start the server:
   ```bash
    go run main.go
    ```
6. **Done!** The server will run on port 8080. You can test it by visiting `http://localhost:8080/v1/health` or by testing the API routes using Postman.

## Docker Setup

1. Install Docker: https://docs.docker.com/get-docker/
2. Install Docker Compose: https://docs.docker.com/compose/install/
3. Set up environment variables, as shown in sample.env
4. Start the server:
   ```bash
   docker-compose up
   ```

## API Endpoints

### Handbook Data

#### Get Unit Information
- **Endpoint:** `/v1/:year/units/:code`
- **Method:** `GET`
- **Description:** Retrieves detailed information about a specific unit
- **Parameters:**
  - `year`: The year of the handbook, ranging from `2020` to `2025`, or `current`
  - `code`: The unit code (e.g., `FIT3175`)
- **Examples:**
```bash
curl 'localhost:8080/v1/2025/units/FIT2004'
curl 'localhost:8080/v1/current/units/FIT3175'
```

#### Get Course Information
- **Endpoint:** `/v1/:year/courses/:code`
- **Method:** `GET`
- **Description:** Retrieves detailed information about a specific course
- **Parameters:**
  - `year`: The year of the handbook, ranging from `2020` to `2025`, or `current`
  - `code`: The course code (e.g., `C2000` or `S2000`)
```bash
curl 'localhost:8080/v1/2024/courses/C2000'
```


#### Get Area of Study Information
- **Endpoint:** `/v1/:year/aos/:code`
- **Method:** `GET`
- **Description:** Retrieves detailed information about a specific area of study (e.g. minor, major)
- **Parameters:**
  - `year`: The year of the handbook, ranging from `2020` to `2025`, or `current`
  - `code`: The area of study code (e.g., `SFTWRDEV07`)
```bash
curl 'localhost:8080/v1/current/aos/SFTWRDEV07'
```


#### Check Unit Requisites
- **Endpoint:** `/v1/:year/units/:code/check`
- **Method:** `POST`
- **Description:** Checks if a student meets the prerequisites for a given unit
- **Parameters:**
  - `year`: The year of the handbook, ranging from `2020` to `2025`, or `current`
  - `code`: The unit code (e.g., `FIT3175`)
- **Request Body:**
  - A JSON array of completed units, each with a `code` field
  - **Example:**
    ```json
    [
        {"code": "FIT1045"},
        {"code": "FIT2004"}
    ]
    ```
- **Response:**
  - A JSON object with:
    - `met_requisites`: boolean indicating if requirements are met
    - `message`: array of unmet requisites messages
    - `warning`: enrolment rule warnings, if any. Usually appears in first year units, where they warn about the minimum VCE scores. This basically represents mini-enrolment rules but we can't parse them into a data structure.
  - **Example:**
    ```json
    {
        "met_requisites": true,
        "message": [],
        "warning": "Some warning text"
    }
    ``` 
- **Sample Usage**
```bash
curl 'localhost:8080/v1/2024/units/FIT3152/check' \
--header 'Content-Type: application/json' \
--data '    [
        {"code": "FIT1045"},
        {"code": "FIT2004"}
    ]'
```
Response:
```json
{
  "message": [
    "Requires one of: FIT1006 or ETC1010 or STA1010 or ETC1000 or ETW1000 or ETF1100 or ETW1010 or FIT2086 or ETW2111",
    "Requires one of: FIT2094 or FIT3171"
  ],
  "met_requisites": false,
  "warning": ""
}
```

#### Get Handbook Search API URL
- **Endpoint:** `/v1/handbook/search_url`
- **Method:** `GET`
- **Description:** Retrieves the current handbook search API URL
- **Response:**
  - A JSON object with the `url` field
  - **Example usage:**
    ```bash
    curl http://localhost:8080/v1/handbook/search_url
    ```
    ```json
    {
    "url": "api-ap-southeast-2.prod.courseloop.com"
    }
    ```

### Health Check
- **Endpoint:** `/v1/health`
- **Method:** `GET`
- **Description:** Simple health check endpoint
- **Response:**
  - A JSON object with `status` field
  - **Example:**
    ```json
    {"status": "ok"}
    ```
