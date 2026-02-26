# Purpose

This file is used to track what I learn while working on this project, including my thought process, experiments, mistakes, and insights as my understanding of networking and TCP grows.

## Learning Goals

1. Understand the core ideas behind HTTP and implement parts of the HTTP/1.1 protocol.
2. Learn how internet protocols are built by implementing them from scratch.
3. Develop a deeper understanding of how web applications communicate over TCP/IP.
4. Build strong networking fundamentals relevant to backend and cybersecurity work.

---

## CHAPTER 1: HTTP streams

To simulate how the internet sends data as a stream, I created a simple scenario using a local file (`message.txt`) as the data source.  
I wrote a program that reads the file in fixed-size chunks of 8 bytes and prints each chunk as it is received.

At first, I questioned why such a small buffer size was used, since real-world systems can read much larger amounts of data very quickly.  
This design choice is intentional: using a small buffer makes it easier to observe how stream-based communication works and highlights the fact that data does not arrive as a complete message.

In real-world protocols, data is usually read in larger buffers (for example, 1024 bytes or more) for performance reasons.  
However, regardless of buffer size, an important insight is that a read operation does not guarantee the requested number of bytes — applications must be able to handle partial reads and correctly reconstruct the data.

### Experiments

- changed the buffer size from 8 -> 9 or 10 or more
- observed fewer read calls but same stream behavior

---

## CHAPTER 2: TCP (Transmission Control Protocol)

TCP is one of the core communication protocols of the internet.  
Its main purpose is to provide **reliable, ordered, and error-checked** data transmission between hosts.

For example, if I want to send the message `"Hello, world!"` from one host to another, the data is not sent as a single unit. Instead, it is split into smaller pieces and transmitted across the network.

As the data travels through the network, packets may take different paths and can arrive **out of order**.  
TCP is responsible for reassembling the data in the correct order, detecting missing packets, and requesting retransmission when necessary.

Without a protocol like TCP, there would be no guarantee that the received data is complete or ordered correctly.

> Note:
>
> While studying networking fundamentals, I clarified the difference between packets and frames.
>
> TCP creates segments at Layer 4, which are encapsulated into IP packets at Layer 3.
>
> These packets are then encapsulated into frames at Layer 2 before being transmitted over the network.

### From File Streams to TCP Streams

In earlier experiments, I simulated streaming by reading data from a local file (`message.txt`) in fixed-size chunks.  
In this chapter, I replaced the file-based stream with a real TCP stream using Go’s standard library.

Instead of reading lines from a file, the program now reads data directly from a TCP connection.  
This reinforces the idea that **files and network connections can both be treated as streams**, and that the same read logic applies regardless of where the data originates.

This change helped solidify my understanding that TCP provides a continuous byte stream rather than discrete messages, and that applications are responsible for parsing and reconstructing meaningful data.

### TCP vs UDP

UDP stands for **User Datagram Protocol**.  
TCP and UDP are both **transport-layer protocols**, but they are designed for different use cases.

#### Key Differences

| Feature             | TCP    | UDP    |
| ------------------- | ------ | ------ |
| Connection-oriented | Yes    | No     |
| Handshake           | Yes    | No     |
| Guaranteed order    | Yes    | No     |
| Reliability         | Yes    | No     |
| Performance         | Slower | Faster |

UDP is generally faster because it does not establish a connection, perform a handshake, or guarantee delivery or ordering of packets.

### Conceptual Difference

**TCP** can be thought of as a careful delivery service:  
the sender verifies the receiver, ensures each piece arrives correctly and in order, and retransmits data if something is lost.

**UDP**, on the other hand, simply sends data without checking whether it arrives, arrives in order, or arrives at all.  
The responsibility for handling loss or ordering is left to the application.

### easy to understand memes

![TCP vs UDP diagram](Images/TCPvsUDP.jpeg)
![TCP vs UDP diagram](Images/tcpudp.jpg)

> This are the funny ways i understand TCP and UDP differences

### When to Use Each

- TCP: HTTP/HTTPS, file transfers, emails
- UDP: DNS, video streaming, online gaming, VoIP

In this project, I am using `nc` (netcat) as a **TCP sender**.  
Netcat requires a connection to be established between the sender and receiver before any data can be transmitted.

This behavior directly reflects one of TCP’s core properties: it is **connection-oriented**.  
Data is only sent after the TCP handshake is completed, ensuring that both sides are ready to communicate reliably.

### Files vs Network

One important things from this project is that **files and network connection behave very similarly**.
I started by simple reading and writings to files, then updated my code to be a bit more abstract so it can handle both.
From the perspective of my code, files and network connections are both just streams of bytes that you can read from and write to.

The core difference comes down to **pull vs push**:

- **Files (pull):**  
  When reading from a file, I am in control of the process:
  - **When** to read (e.g., when the program runs)
  - **How much** to read (e.g., 8 bytes at a time)
  - **When** to stop reading (EOF)

> Note: With files, you _pull_ the data at your own pace.

- **Network connections (push):**  
  When reading from a network connection, the data is pushed to me by the remote server.  
  I have no control over:
  - When data arrives
  - How much data arrives at a time
  - When the stream ends

I must be ready to receive data whenever it comes and handle it correctly.  
This distinction reinforced my understanding that **network streams are asynchronous and unpredictable**, unlike files which are synchronous and controlled.

## CHAPTER 3: Requests

**HTTP/1.1** is a text based protocol that works over TCP

> Note:
>
> I am following RFC: 9112 and 9110
>
> 9110 is the semantic
>
> 9112 is Message Syntax and Routing

### TCP to HTTP

> Why am I using HTTP and not just TCP?

TCP provides **reliable, ordered, and complete** delivery of bytes between hosts.  
However, TCP does **not tell us what type of data** is being sent — it could be text, an image, a video, or an email. TCP only guarantees that the bytes arrive correctly and in order.

HTTP, on the other hand, is an **application-layer protocol** built on top of TCP.  
It gives us a way to specify **what kind of data** is being sent and received (e.g., text/html, image/png, application/json), and provides **semantic meaning** like requests, responses, headers, and status codes.

In short:

- **TCP** → ensures data is delivered reliably
- **HTTP** → ensures data is meaningful and interpretable by applications

### This is what HTTP requests looks like

```
GET /index.html HTTP/1.1
Host: DevLai.dev
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
Accept: text/html
```

HTTP allows us to specify a **destination within the server** and provides metadata about the request.

### Breakdown of an HTTP Message

| Part         | Example                    | Description                                                                                                                 |
| ------------ | -------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| Start-line   | GET /index.html HTTP/1.1   | The first line of the request (or response). It specifies the **HTTP method**, **resource path**, and **protocol version**. |
| Header lines | Host: DevLai.dev           | Zero or more lines containing headers. Headers are key-value pairs providing metadata about the request.                    |
|              | User-Agent: Mozilla/5.0    | Another example of a header line.                                                                                           |
| Blank line   | (empty line)               | Separates headers from the message body. Required even if there is no body.                                                 |
| Message body | (none in this GET request) | Optional. Contains data sent to the server (e.g., JSON, form data).                                                         |

> Both HTTP requests and responses follow this format. Collectively, these are called **HTTP Messages**.

### Key Takeaways

1. **Start-line:** Declares the request or response.
2. **Header lines:** Zero or more lines containing metadata about the message.
3. **Empty line:** Separates headers from the body.
4. **Message body:** Optional data payload.

#### Formal Representation

```go
METHOD /resource-path PROTOCOL-VERSION\r\n
field-name: value\r\n
field-name: value\r\n
...\r\n
\r\n
[message-body]
```

### cURL

- stands for client URL
- is a command line tool for making http requests.

## CHAPTER 4: Request Lines

I created a simple test file for request parsing.
In this project, I am not using table-driven tests.

Following ThePrimeagen’s approach, " instead of writing tests for every function, I focus on testing the parts I am unlikely to get right on the first attempt—" especially parsing logic and protocol boundaries.

If, in the future, something becomes unclear or error-prone, I can always add more tests.

### Parsing the Request Line

At this point, I already have code that handles plain-text TCP data.
The next step is to convert that raw text into structured data, while ensuring it conforms to the HTTP/1.1 protocol as defined in RFC 9110 and RFC 9112.

for example, given

```go
POST /rice HTTP/1.1
Host: localhost:42069
User-Agent: curl;/7.81.0
Accept: */*
content-length: 17

{"Type": "Basmati"}
```

The HTTP parser should return a struct that looks like this

```go
type Request struct {
    RequestLine RequestLine
    Headers     map[string]string
    Body        []byte
}
```

#### Goal of This section

The goal of this section is to parse the request start-line correctly, using the server’s raw input and following the HTTP message parsing rules defined in RFC 9112.

```go
METHOD SP request-target SP HTTP-version CRLF
```

> Examples

```go
POST /rice HTTP/1.1
GET /index.html HTTP/1.1
```

#### What I Learned from RFC 9112 Section 3

- Request line must be: `METHOD SP request-target SP HTTP-version CRLF`
- Methods are case-sensitive (GET not get)
- Only single space (SP) allowed between components
- Must end with `\r\n` (CRLF)
-

### Parsing Strategy

To begin parsing, I created a function called `ParseRequestLine`.

This function:

- Takes the raw request text as a string

- parses the request line

- Returns:
  - A structured RequestLine

  - The remaining unparsed text (headers + body)

  - An error if the request line is invalid

This allows the parser to:

1. Validate protocol correctness early
2. Keep request-line parsing independent from headers and body parsing
3. Fail fast if the request if malformed

> Note:
>
> Request line is composed of three components
>
> - Method
> - RequestTarget
> - HTTPVersion
>
> Components are separated by space and not tabs
>
> The line must end by CRLF( **\r\n** ). Prime calls this “registered nurse”, which is honestly kind of funny

#### Validation & Testing

I implemented a parsing method for the request start-line and added targeted tests to validate edge cases, including:

- Ensuring the HTTP method is uppercase
  - Lowercase methods like get are rejected

- Rejecting invalid HTTP versions
  - Only HTTP/1.1 is accepted

- Detecting invalid request-line formats

- Rejecting request lines with more than three components

### Parsing a Stream

TCP (and HTTP over TCP) is a stream-based protocol, not a message-based one.

Data arrives as an ordered stream of bytes, which may be split or combined arbitrarily by the transport layer. TCP guarantees that data is delivered in order, but it does **not** guarantee that application-level messages will be received in complete units.

Therefore, an HTTP parser must handle incomplete reads and incrementally parse incoming data until enough bytes are available to determine a complete structure.

---

#### Incomplete vs complete data

Instead of receiving a full HTTP request at once, I might receive only the first few characters:

#### Incomplete

```bash
GE
```

So, I create a parser that can handle incomplete reads. It must be smart enough to know that parsing is not finished yet and keep reading until it receives the full request line:

#### Complete

```bash
GET /rice HTTP/1.1

```

---

#### Old approach

- `ReadAll()` → parse everything at once
- Works in toy cases, but bad practice for real networking

**Why it’s bad:**

- TCP is a stream, not message-based
- I might get `"GE"` now and `"T /"` later
- I don’t know:
  - how big the request is
  - when it ends
  - how fast it arrives

---

#### New approach

- Read small chunks
- Parse incrementally
- Keep parser state

---

#### About the `chunkReader` and tests

- `chunkReader` simulates slow / fragmented network reads
- Tests use:
  - `numBytesPerRead: 1`
  - `numBytesPerRead: 3`

- This forces the code to handle:
  - incomplete request lines
  - partial `\r\n`
  - split tokens like `GE` + `T`

The tests are basically saying:

> “If your parser assumes it gets the full line at once, it will fail.”

---

#### Parser state (`initialized` / `done`)

- The `Request` struct becomes a small state machine
- Two states for now:
  - `initialized` → still parsing
  - `done` → request line parsed

My understanding: **track the state of the parser** so it knows whether it should keep consuming data or stop.

---

#### Buffer and byte tracking

This lesson emphasizes tracking:

- bytes read from the reader
- bytes consumed by the parser
- bytes remaining unparsed in the buffer

I created a buffer with these key points in mind:

- The buffer size (8, 1024, etc.) is **not**:
  - how much data exists
  - how much will be parsed

- It is just:
  - temporary storage for incoming bytes

What really matters is how many bytes were:

- read
- parsed
- left unparsed and shifted forward for the next read

#### Reading vs Parsing

**Reading** is moving data from the reader into our program.

**Parsing** is interpreting that data — for example, transforming raw `[]byte` into a `RequestLine` struct.

> Note: once the data is parse we can discard it from buffer to save memory.

### State Machine

A state machine is a system that can be in one of many possible states, and its behavior depends on its current state.

```go
func add(a, b int) int {
  return a + b
}
```

The example above is **not** state machine as it does not have any internal state.t does not maintain any internal state. It takes two inputs and returns a result without remembering anything.

```go
type Counter struct {
  count int
}

func (c *counter) Add(a int) int {
  c.count += a
  return c.count
}
```

This example has internal state (`count`). Each time we call Add, it changes the state of the `Counter`. The output depends not only on the input, but also on the current internal state. Therefore, it is stateful.

### Connect the Parsing

In this lesson, I felt like I built something real. It almost felt like magic.

I sent a `curl` request from one terminal, and on the other side I could see the parsed request line — even though it was just `localhost`. I could see exactly what was being sent and how it was interpreted.

In the first chapter, I was reading messages from files and processing them.

In the second chapter, I switched to reading data from the network. I was receiving raw bytes over TCP and processing them.

But in the fourth chapter, something changed. Instead of just reading lines, I parsed the HTTP start line and extracted its three components:

- Method
- Target
- Version

Previously, I was using a `GetLineChannel` function that only split input by newline.

This time, I removed that function and replaced it with my own `RequestFromReader` function. Instead of just printing a raw line, I parsed the request into a structured `RequestLine` and printed it like this:

```bash
Request line:
- Method: GET
- Target: /
- Version: 1.1
```

That’s when it clicked — I wasn’t just reading data anymore. I was interpreting a protocol.

#### I moved from

##### Text Processing

" Read a line, Print a line "

##### To

#### Protocol Parsing

" Read a byte stream, Interpret structured meaning"

#### Before and After

**Before**:

```bash
Network → Read → Split on newline → Print
```

**After**:

```bash
Network → Read (stream) → State Machine → Structured Data → Print

```

`I moved from reading lines to parsing protocols.`

## CHAPTER 5: HTTP Headers

**Headers** , a metadata fields that accompany HTTP requests and responses.

The RFC does not call them header, The RFC uses the term `Field-line`

> Each field line consists of a case-insensitive field name followed by a colon (`":"`), optional leading whitespace, the field line value, and optional trailing whitespace."

```bash
field-line   = field-name ":" OWS field-value OWS
```

There can be an unlimited amount of whitespace before and after the `field-value`. However, when parsing a `field-name`, there must be no space betwixt the colon and the `field-name`.In other words, these are valid:

```bash
'Host: localhost:42069'
'          Host: localhost:42069    '
```

But this is not:

```bash
Host : localhost:42069
```

### Header Parser — My Learning

Creating the Headers parser gave me quite a few hiccups — I was stuck for a long time.

Here’s how I approached it:

### setup

- I was given a `Headers` type (a `map[string]string`) and a `Parse` function which returns (int, bool, error).

- I set a global variable `rn` to represent `CRLF` (`\r\n`).

- Similar to the `RequestLine` parser, I created:
  - a `read` variable to track how many bytes I have already consumed

  - a loop to iterate through the data

  - an index `i` to find the next CRLF

### Handling Incomplete Data

- While reading, I check for CRLF in the current unparsed slice:

```go
i := bytes.Index(data[read:], rn)
if i == -1 {
    break // need more data
}

  - data -> the full byte slice i got from the network
  - read -> how many bytes i have already parsed and consumed
  - i    -> index of the next \r\n in the current unparsed slice
```

- If no CRLF is found, the parser cannot continue yet, so we wait for more data and parse again.

### Detecting End of Headers

- If we find an empty line (i.e., i == 0), it means we’ve reached the end of headers:

```go
if i == 0 {
    done = true
    read += len(rn)
    break
}

```

- We advance the `read` counter to consume the empty line and signal to the caller that headers are finished.

### parsing One Header Line

- We take one header line from the buffer:

```go
name, value, err := parseHeader(data[read:read+i])
if err != nil {
return 0, false, err
}
```

- parseHeader splits the line into name and value, trims whitespace, and immediately handles malformed headers.

- data[read:read+i] is exactly one header line without the trailing CRLF.

- Then we advance the read pointer:

```go
read += i + len(rn) // move past this line
```

- After parsing, we store the header in the map:

```go
h[name] = value
```

Example:

```go
h["Host"] = "example.com"
h["User-Agent"] = "curl/8.0"
```

### Constraints

The header parsing mostly working but not according to RFC

#### Case Insensitivity

Field names are case-insensitive so Content-Length and content-length are the same and we have to account for this.

#### Valid Characters

Field-name has implicit definition of a token, as defined in RFC 9110, TOken are short textual identifiers that do not include whitespace or delimiters.

```bash
token          = 1*tchar

  tchar          = "!" / "#" / "$" / "%" / "&" / "'" / "*"
                 / "+" / "-" / "." / "^" / "_" / "`" / "|" / "~"
                 / DIGIT / ALPHA
                 ; any VCHAR, except delimiters
```

so, a `field-name` must only contain:

- Uppercase Letters A-Z
- Lowercase Letters a-z
- Digits: 0-9
- Special characters: ``!, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~``

### Multiple Values

According to RFC: 9110, in **[5.2](https://datatracker.ietf.org/doc/html/rfc9110#name-field-lines-and-combined-fi)** it is mentioned, that its is perfectly valid to have multiple values for single header key.
for example:

```bash
Example-Field: Foo, Bar
Example-Field: Baz
```

### Add to Parse

While adding to my Parse function, I ran into a few problems , similar to what the video course showed — where my code panicked multiple times.

The causes were tiny mistakes, like:

- Adding extra spaces while testing

- Miscalculating the number of bytes to slice

#### Using slog for Debugging

I discovered slog for structured logging, which I found very helpful. I used it multiple times in my code to debug issues like this:

```go
// debugging
slog.Info("Read from reader",
    "n", n,
    "bufLen", bufLen,
)

readN, err := request.parse(buf[:bufLen])
if err != nil {
    return nil, err
}

// debugging
if readN > bufLen {
    slog.Info("Parse returned more than buffer length",
        "readN", readN,
        "bufLen", bufLen,
    )
    return nil, fmt.Errorf("parse returned readN > bufLen: %d > %d", readN, bufLen)
}
```

- The first log shows how many bytes were read and the current buffer length.

- After parsing, the second log checks if `parse` returned more bytes than the buffer length, which could indicate a bug.

#### What I Learned

The panic happened because of Go slice rules:

```go
buf[:bufLen+n] // This caused panic
```

- `bufLen` already includes the newly read bytes

- By adding `n` again, I doubled the length, making the start index larger than the end index → slice panic

In simple terms:

> My start index was greater than my end index because I incremented `bufLen` before slicing and then added `n` again.

### Live headers

I was starting to really understand how HTTP headers work, how they are parsed, and what constraints they must follow according to the RFCs.

- I could see live headers coming from my curl requests in real time.

- I saw exactly how each header line was split into name and value, and how the parser enforced rules:
  - No spaces in the field name

  - Optional whitespace around the value

  - Empty line marking the end of headers

This was the moment everything clicked:

> I wasn’t just reading data or parsing lines anymore It was interpreting structured protocol data live, exactly like a real HTTP server does.

I got full insight into why headers are structured the way they are, and how incremental parsing, CRLF detection, and buffer management all work together.

## CHAPTER 6: HTTP Body

An HTTP/1.1 message consists of:

- A start-line
- Followed by CRLF
- Zero or more header field lines
- An empty line (CRLF)
- An optional message body

From the RFC:

```bash
  HTTP-message   = start-line CRLF
                   *( field-line CRLF )
                   CRLF
                   [ message-body ]
```

A message can be.

- A **request** (client -> server)
- A **Response** (server -> client)

Syntactically, they are almost identical.
The main difference is:

- The `start-line`
- The `field-line`(headers)
- The extra `CRLF` separating headers an body
  Now I need to parse the **Message Body**

There are many edge cases

According to RFC 9110 (**[section 8.6)](https://datatracker.ietf.org/doc/html/rfc9110#section-8.6)**)

> A user agent should send Content-Length in a request when the method defines a meaning for enclosed content and it is not sending Transfer-Encoding.

In my Implementation:

- If there is no `Content-Length` header -> I assume there is **no body**.
- If `Content-Length` exists -> I must read exactly the many bytes

### Setup

To support body parsing:

- I update `(r *Request)`
- Added a new state in the state machine: `StateBody`
- Added a `.Body` field to the `Request` struct (type: `[]byte`)

Code:

```go
type Request struct {
    RequestLine RequestLine
    Headers     Headers
    Body        []byte
    state       parserState
}
```

### Parsing Strategy

1. Parse headers.
2. Extract Content-Length using GetInt.
3. Store that value as length.
4. If length == 0 → there is no body.
5. Otherwise:

- Read incrementally from the stream.
- Append only what is needed.
- Stop once total body size equals Content-Length.

Because TCP is a stream:

- I might not receive the entire body in one read.
- I must handle partial body reads.

### Core Logic

```go
remainingData := min(length - len(r.Body), len(currentData))
```

where:

- `length` -> Total body size required (from `Content-Length`)

- `len(r.Body)` -> How many bytes you've already collected

- `length - len(r.Body)` -> How many bytes are still missing

- `len(currentData)` -> How many bytes are available right now in this chunk

- min(...) ensures:
  - I never read more than required
  - I never read more than what i currently have
  - I respect exact body size defined by `Content-Length`

#### Example

suppose:

```bash
Content-Length: 10
```

**First read:**

Get 6 bytes:

```bash
len(r.Body) = 0
len(currentData) = 6
```

so:

```bash
length - len(r.Body) = 10 - 0 = 10
min(10, 6) = 6
```

**Second Read:**

I get 10 more bytes:

```bash
len(r.Body)  = 0
len(currentData) = 10
```

Now:

```bash
length - len(r.Boy) = 10 - 6 = 4
min(4, 10) = 4
```

you ony append 4 bytes. Even though i received 10 bytes
Because the body should only be 10 total

### Appending the Body

```go
r.Body = append(r.Body, currentData[:remainingData]...)
```

Why This?:

- HTTP bodies are raw bytes.
- Converting to string is unnecessary.
- `append` avoids extra allocations and conversions

## CHAPTER 7: HTTP Response

Now I understand what an HTTP request is.

A request is a structured block of text that contains:

- A start line (request line)
- Headers (field-lines)
- An optional body

When a client sends a request, it travels through multiple layers of the network stack and eventually reaches the destination server.

Example HTTP request

```bash
GET /coffee.html HTTP/1.1
User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)
Host: www.coffee.com
Accept-Language: en-us
Accept-Encoding: gzip, deflate
Connection: Keep-Alive
Content-Length: 23

{
  "coffee": "is good"
}

```

so when this request reaches the server, server parses the status line , like asking to server, "hey server, i need to route out to this place" and server parses all the things , queries, headers and things, and server will be going to either a database, third party services like aws or any GPT to get `response` back and that response is sent back with response message which has similar structure to request where it has status line, field lines and response body, as follows:

`What Happens on the Server

When the request reaches the server:

1. The server parses the request line:
   - Method (`GET`)

   - Target (`/coffee.html`)

   - Version (`HTTP/1.1`)

2. The server parses all headers.

3. If a body exists (based on `Content-Length` or `Transfer-Encoding`), it reads the body.

4. The server then decides what to do:
   - Route to a handler

   - Query a database

   - Call third-party services (e.g., cloud services, APIs)

   - Perform internal logic

After processing, the server sends back a `response`.

### HTTP Response structure

An HTTP response has a very similar structure:

- A status line

- Headers (field-lines)

- An empty line

- An optional response body

Example:

```bash
HTTP/1.1 200 OK
Accept-Ranges: bytes
Age: 294510
Cache-Control: max-age=604800
Content-Type: text/html; charset=UTF-8
Date: Fri, 21 Jun 2024 14:18:33 GMT
Etag: "3147526947"
Expires: Fri, 28 Jun 2024 14:18:33 GMT
Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
Server: ECAcc (nyd/D10E)
X-Cache: HIT
Content-Length: 1256

Coffee is good!
```

- `TTP/1.1` → Version
- `200` → Status code
- `OK` → Reason phrase

Common status codes:

- `200 OK` → Success
- `404 Not Found` → Resource does not exist
- `500 Internal Server Error` → Server failed to handle request

**The Request/Response Model:**

```bash
Client → HTTP Request → Server
Server → HTTP Response → Client
```

### Server

Now I am upgrading from a simple `tcplistener` to an actual `httpserver.`

Previously, I was:

- Reading a single line
- Reading data from a file
- Reading raw bytes from a TCP connection

Now, I am moving to the next level:

Building a server that:

- Accepts valid HTTP requests
- Parses them correctly
- Sends back valid HTTP responses

#### Server Setup

I was given `cmd/httpserver/main.go` file:

```go
const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
```

This file:

- starts the server in `42069`
- waits for the system signals(`SIGINT`/`SIGTERM`) to gracefully stop
- Calls `server.Close()` on exit

#### creating server

Inside my `server` package, i implemented following things:

1. `type Server struct`

```go
type Server struct {
  listener net.Listener
}
  - Keeps track of the server state
  - stores the TCP listener
```

2. `func Serve(port int) (*server, error)`

```go
func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.Listen()
	return server, nil
}
```

- Creates a TCP listener
- Instantiates a new `server`
- starts the `Listen` loop in a goroutine
- returns the server instance

3. `func(s *server) Close() error`

- Closes the TCP Listener
- Sets the server state so that further erros from Accepts can be ignored

4. `func (s *Server) listen()`

```go
func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}
```

- Continuously accepts new connections
- Each connection is handled in its own goroutine
- Ensures the server can handle multiple simultaneous Clients
  2

5. `func (s *Server) handle(conn net.Conn)` -

```go
func (s *Server) handle(conn net.Conn) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()
}

```

- Handles a single connection
- Writes a raw HTTP response
- closes the connection

> As of in this stage, i am hard coding the response msg, in future will add dynamic responses

### Response

An HTTP response follows the sme overall HTTP message format:

```bash
HTTP-message   = start-line CRLF
                 *( field-line CRLF )
                 CRLF
                 [ message-body ]
```

#### status line

The status line has the following structure:

```bash
status-line = HTTP-version SP status-code SP [ reason-phrase ]
```

Examples:

```bash
HTTP/1.1 200 OK
```

```bash
HTTP/1.1 404 Not Found
```

> A server MUST send the space that separates the status code from the phrase
>
> Even if the reason phrase is empty, the trailing space after the status code must still be present

##### my Implementation

At this stage, i implemented `StatusCode` enum-like type supporting:

- `200`
- `400`
- `500`

#### writing the Status line

I used a `switch` on the status code and wrote the raw bytes:

```go
func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
  switch statusCode {
  case StatusOK:
    _, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
    return err
  // other cases...
  }
}
```

Current approach:

- Hard-coded reason phrase
- writes directly to the provided writer
- keep things simple for now

#### Writing Headers

After the status line, the server must write the headers(field lines)
I implemented a helper that iterates through my header collection using `ForEach` and formats each header line:

```go
func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.Write(b)
	return err
}
```

what this does:

- Iterates overall headers
- formats each headers as

```bash
  Name: Value\r\n
```

#### Other Common Headers

Besides the core headers i implemented there are several other important HTTP headers as well:

- **Content-Encoding:**
  Indicates whether the response body has been encoded or compressed (for example: gzip, br).
  If present, the client must decode the body using the specified encoding before using the content.

- **Date:**
  Specifies the date and time at which the message was generated by the server.
  This helps clients and intermediaries understand response freshness and is commonly included in HTTP responses.

- **Cache-Control:**
  Provides directives that control how the response may be cached by browsers and intermediary caches.
  It is useful for defining caching behavior such as:
  - whether the response can be cached

  - how long it stays fresh (max-age)

  - whether revalidation is required

#### Flow so far

At this point server response flow is:

- write status line
- Write headers
- (next step) write body

Right now the server is still fairly low-level and manual, but it correctly follows the HTTP message structure.

### Handler

It's time to define the handler function. This will allow the user's of our `server` package to define their own logic for handling requests:

```go
type Handler func(w io.Writer, req *request.Request) *HandlerError
```

- similar to of Go standard library:

```Go
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
```

- **Difference**: I use a general io.Writer instead of http.ResponseWriter for simplicity and flexibility.

#### Implementation

1. Creating `Handler` function and `HandleError`
   ```Go
   Type HandlerError struct {
   	StatusCode response.StatusCode
   	Message    string
   }
   type Handler func(w io.Writer, req *request.Request) *HandlerError
   ```

- `HandleError` allows the handler to return an HTTP status code and message in case of errors

2. Accept `handler` in Server
   Update `server.Serve` to accept a handler function

   ```Go
    server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
   ...
   }
   ```

3. The server now uses the handler to process requests and generate responses:

```go
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, *headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)
	if handlerError != nil {
		// updating the body and Content-Length
		body := []byte(handlerError.Message)
		headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))

		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, *headers)
		conn.Write(body)
		return
	}

	// 200 OK PATH
	body := writer.Bytes()
	headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, *headers)
	conn.Write(body)
}
```

#### Flow Explanation:

- Parse the request from connection
- if parsing fails - return **400 Bad Request**
- Pass request to handler:
  - if handler returns and error - > sends corresponding status + message
  - otherwise -> write handler's response to the body and return **200 OK**

1. Example Handler in `main.go`

```go
func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "your problem, not my problem",
			}
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "My bad, sorry",
			}
		}
		w.Write([]byte("we all good, frfr"))
		return nil
	})
  ...
}
```

- Routes requests to different paths
- Returns appropriate errors or writes a normal response

> Note:
>
> The `Handler` is responsible for reporting errors or writing the
> body.

### Refactor

As mentioned in the course, the previous design constraining the library:

- Errors were always returned as plain text
- Headers were always the same
- Users had limited control over the response

To fix this, the handler signature needed to be more flexible.

Previously:

```go
type Handler func(w io.Writer, req *request.Request) *HandlerError
```

New Design

```go
type Handler func(w *response.Writer, req *request.Request)
```

Now the handler receives a custom response.Writer, which gives full control over:

- Status line
- Headers
- Body

This encapsulates boilerplate while giving users flexibility.

#### Goal of the Refactor

The main goal was to allow the `httpserver` to return HTML responses instead of being locked to plain text.

To achieve this:

- Stop depending directly on io.Writer
- Create a custom response writer
- Let handlers fully control the HTTP response

#### response.Writer

**struct**

```Go
type Writer struct {
	writer io.Writer
	state  writerState
}
```

**Purpose:**

- writer → underlying connection
- state → enforces correct write order using a state machine

This prevents invalid HTTP responses.

**WriteStatusLine**

```Go
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateStatusLine {
		return fmt.Errorf("cannot write status line in current state")
	}
	switch statusCode {
	case StatusOK:
		_, err := w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		w.state = stateHeaders
		return err
	case StatusBadRequest:
		_, err := w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.state = stateHeaders
		return err
	case StatusInternalServerError:
		_, err := w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.state = stateHeaders
		return err
	default:
		return fmt.Errorf("Great, you found new Status. Unrecognized error code")
	}
}
```

**What it does**:

- Writes the HTTP status line
- Advances state → stateHeaders
- Prevents writing status line twice

**WriteHeaders**

```Go

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateHeaders {
		return fmt.Errorf("cannot write headers in current state")
	}
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.writer.Write(b)
	w.state = stateBody
	return err
}
```

**What it does**:

- Formats header field lines
- Writes the blank line after headers
- Moves state → stateBody

**WriteBody**

```Go
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write body in current state")
	}
	n, err := w.writer.Write(p)
	w.state = stateDone
	return n, err
}

```

**Important:**

- p is the payload (body bytes)
- Ensures body is written only at the correct time
- Moves state → stateDone

#### Updated Server handle

```Go
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*headers)
		return
	}
	s.handler(responseWriter, r)
}

```

**Flow**:

- Create response writer
- Parse request
- On parse error → send 400
- Otherwise → delegate to handler

#### Update main Handler

```Go
func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := body200()
		status := response.StatusOK

		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = body400()
			status = response.StatusBadRequest
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			body = body500()
			status = response.StatusInternalServerError
		}
		h.Set("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
  ...
}
```

## CHAPTER 8: Chunk Encoding

At this point, I assumed HTTP messages were parsed and sent as a single block of data. This works fine for many use cases, but what if I want to send updates bit by bit?

then i Remember the basic structure of an `HTTP message`:

```zsh
     HTTP-message   = start-line CRLF
                      *( field-line CRLF )
                      CRLF
                      [ message-body ]
```

The `[message-body]` can be flexible. Instead of using `Content-Length`, HTTP allows sending the body in chunks using the `Transfer-Encoding: chunked` header.

### chunked transfer-Encoding

Format of a chunked response

```bash
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
... repeat ...
0\r\n
\r\n
```

- `<n>` → hexadecimal number indicating the size of the chunk in bytes
- `<data of length n>` → the actual data of that chunk
- Repeat until all data is sent
- Last chunk has size 0 followed by an empty line (\r\n) to indicate end of message

### Example: Plain Text Response

```bash
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

1E
I could go for a cup of coffee
C
But not Java
12
Never go full Java
0
```

- 1E → 30 bytes
- C → 12 bytes
- 12 → 18 bytes
- 0 → signals end of body

Chunked encoding is most often used for:

- streaming large amount of data (like big files)
- real-time updates(like chat-style application)
- sending data of unknown size(like a live feed)

### chunk Format

In RFC 9112, section 7.1

each chunk looks like this:

```bash
chunk-size [chunk-extension] CRLF
chunk-data
CRLF
```

- `chunk-size` → hexadecimal number representing chunk length
- `chunk-extension` → optional, rarely used
- `chunk-data` → the actual bytes of this chunk

This allows HTTP to send messages incrementally without knowing the total size in advance.

### Implementation of chunked Encoding:

To support sending chunked responses, I added two methods to the response package:

- `func (w *Writer) WriteChunkedBody(p []byte) (int, error)` → writes a single chunk

- `func (w *Writer) WriteChunkedBodyDone() (int, error)` → writes the final zero-length chunk to signal the end

For `WriteChunkedBody`:

1. make sure we are in body sate
2. writing the chunking size in hex, followed by `CRLF`

```go
	n := len(p)
	_, err := w.writer.Write([]byte(fmt.Sprintf("%x\r\n", n)))
	if err != nil {
		return 0, err
	}
```

3. write the actual chunk itself

```go
	_, err = w.writer.Write(p)
	if err != nil {
		return 0, nil
	}
```

4. printing the CRLF after the chunk

```go
	_, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
```

For `writeChunkedBodyDone`:

1. check if we are in body state
2. write the final zero-length chunk with trailer

```go
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}
```

**Updating the Server to Use Chunked ResponsesUpdate the `server`**:

1. Check if the request target starts with /httpbin/stream.
2. Build the external URL using the request target:

```go
target := req.RequestLine.RequestTarget
res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
```

3. If there is an error, send a `500 Internal Server Error`. Otherwise:

- Write the status line (`200 OK`)
- Delete `Content-Length` and add `Transfer-Encoding: chunked`
- `Set Content-Type`
- Write headers

4. Read from the external response in a loop using a small buffer, writing each chunk with `WriteChunkedBody`:

```go
buf := make([]byte, 32)
for {
    n, err := res.Body.Read(buf)
    if n > 0 {
        w.WriteChunkedBody(buf[:n])
    }
    if err != nil {
        break
    }
}
```

5. After all chunks are sent, call `WriteChunkedBodyDone()`

```go
w.WriteChunkedBodyDone()
return
```

This way, the server can stream data to the client incrementally, without knowing the total size in advance.

## Mistakes & Realizations

- Initially assumed each `Read()` returns a full message → wrong, learned TCP is stream-based.
- Thought UDP might be "unreliable" for all small messages → realized some apps handle reliability at the application layer.
- In the third lesson of request, I thought 0 and nil represents failure but it's more like "needs more data"
- Buffer size is not equal to data size
- Parsing is incremental, not "loop through buffer"
- Confused "pull vs push" → really about synchronous (files) vs asynchronous (network) data availability
- Initially said "TCP splits into packets" → learned TCP creates segments, IP creates packets

**Chapter 5**:

Lesson 1

- Partial reads are normal over TCP — your parser must handle them.

- Empty line detection is key to know when headers are done.

- Advancing the read pointer after each line prevents parsing the same line twice.

- Immediate error handling ensures malformed headers don’t propagate.

Lesson 4

- Structured logging (slog) is incredibly useful for debugging incremental parsing.
- Small mistakes (extra space, wrong byte count) can cause panics in network parsers.

**Chapter 6**:

- Parsing the body is different from parsing headers:
  - Headers are line-based
  - Body is **byte-count based**
- Headers stop at `\r\n\r\n`, Body stops at exact byte count
- Header is delimiter-base and Body is length-based

**Chapter 7**:

- In a request, the first line is called **request line**
  In a response, the first line is called **status line** and it contains:
  - HTTP version
  - Status code
  - Reason phrase
- The request and response share the same high-level structure

## Security Insights

- TCP’s connection-oriented nature is reliable, but also vulnerable to SYN flood attacks.
- Partial read handling is important to prevent buffer overflows or request smuggling.

---

```

```

```

```
