# Go-Redis

A simple, persistent key-value store inspired by Redis, written in Go.

## Features

- **TCP Server**: Listens on port `6379`, the default Redis port.
- **Key-Value Store**: Supports `SET`, `GET`, and `DELETE` operations.
- **Persistence**: The database is saved to a `redis.db` file and reloaded on startup.
- **Concurrent**: Handles multiple client connections concurrently.

## Getting Started

### Prerequisites

- Go (1.x or later)

### Building and Running

1.  **Clone the repository (or just use the existing code):**

2.  **Build the application:**

    ```sh
    go build -o go-redis
    ```

3.  **Run the server:**

    ```sh
    ./go-redis
    ```

    The server will start and listen on port `6379`.

## Usage

You can connect to the server using a TCP client like `telnet` or `netcat`.

```sh
telnet localhost 6379
```

### Supported Commands

-   **SET key value**

    Sets the string value of a key.

    ```
    > SET name John
    OK
    ```

-   **GET key**

    Gets the value of a key.

    ```
    > GET name
    John
    ```

-   **DELETE key**

    Deletes a key.

    ```
    > DELETE name
    OK
    ```
    ```
    > GET name
    (nil)
    ```

### Persistence

The database is automatically persisted to the `redis.db` file in the same directory whenever a `SET` or `DELETE` command is executed. When the server starts, it will load the data from this file if it exists. 