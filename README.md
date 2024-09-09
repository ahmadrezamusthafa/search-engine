# Search Engine

This project implements a search engine API where users can:
1. **Index documents** with various content types.
2. **Search documents** using multiple query terms and return matching results with relevance scores.

The system is built in Go with (BadgerDB or Redis) for fast indexing and search operations. It includes features such as stop word filtering and Time-to-Live (TTL) for indexed documents.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [API Endpoints](#api-endpoints)
  - [Index a Document](#index-a-document)
  - [Search for Documents](#search-for-documents)
- [Installation](#installation)

---

## Features

- **Document Indexing**: Index documents with both string and object content.
- **Keyword Search**: Perform searches with multiple query terms and retrieve matching documents.
- **Result Scoring**: Each search result comes with a relevance score.
- **Stop Word Filtering**: Customize stop words during indexing to exclude common terms from searches.
- **TTL (Time-To-Live)**: Indexed documents expire after a set duration, ensuring storage efficiency.

---

## Tech Stack

- **Backend**: Golang for API development.
- **Database**: BadgerDB or Redis for fast key-value storage and search operations.

---

## API Endpoints

### Index a Document

**URL**: `/index`  
**Method**: POST  
**Content-Type**: `application/json`

#### Example Request

```json
{
  "id": "tu:id:2",
  "content": {
    "string": "Journal no: 980034 TRANSFER DARI Bpk TEDDY ACHMAD ZAELANI",
    "object": {
      "id": 650282045,
      "amount": 71495000,
      "unique_code": 150,
      "total_amount": 71495150,
      "sender_bank": "mandiri",
      "sender_name": "pt people intelligence indonesia (julian alimin)",
      "status": 10,
      "flip_receiver_bank_code": 32,
      "flip_receiver_bank_type": "",
      "virtual_account_number": "8558151502502621",
      "remark": " ffb20508 via api",
      "parent_id": null,
      "parent_type": 0,
      "created_from": 10,
      "created_at": 1725261101,
      "confirmed_at": 1725261101,
      "user_id": 4416295,
      "notes": ""
    },
    "object_indexes": [
      "id",
      "sender_bank",
      "total_amount",
      "remark"
    ]
  },
  "stop_words": [
    "api"
  ]
}
```

In this example:
- The `id` field represents the document ID.
- The `content` field contains the document's content as both a string and an object.
- `object_indexes` indicates the fields within the object that are indexed.
- `stop_words` allows specific terms to be filtered out during indexing.

#### Response

- **200 OK**: Document indexed successfully.
- **400 Bad Request**: Invalid or missing fields in the request payload.

---

### Search for Documents

**URL**: `/search?query=term1&query=term2`  
**Method**: GET

#### Example Request

```bash
GET /search?query=teddy&query=achmad&query=500150&query=reza
```

In this example, the `query` parameters represent the search terms.

#### Response

Returns an array of matching documents with their IDs, relevance scores, and content.

#### Example Response

```json
{
  "status": "success",
  "data": [
    {
      "id": "tu:id:1",
      "score": 2.8520289416615365,
      "data": {
        "object": {
          "amount": 500000,
          "confirmed_at": 1725261101,
          "created_at": 1725261101,
          "created_from": 10,
          "flip_receiver_bank_code": 32,
          "flip_receiver_bank_type": "",
          "id": 650282041,
          "notes": "",
          "parent_id": null,
          "parent_type": 0,
          "remark": " ffb20508 via api",
          "sender_bank": "bca",
          "sender_name": "pt people intelligence indonesia",
          "status": 10,
          "total_amount": 500150,
          "unique_code": 150,
          "user_id": 4416295,
          "virtual_account_number": "8558151502502621"
        },
        "string": "Journal no: 980034 TRANSFER DARI Bpk TEDDY ACHMAD ZAELANI"
      }
    },
    {
      "id": "tu:id:4",
      "score": 1.2602676010180822,
      "data": {
        "object": {
          "amount": 5000000,
          "confirmed_at": 1725261101,
          "created_at": 1725261101,
          "created_from": 10,
          "flip_receiver_bank_code": 32,
          "flip_receiver_bank_type": "",
          "id": 650282047,
          "notes": "",
          "parent_id": null,
          "parent_type": 0,
          "remark": " ft232312312 via api",
          "sender_bank": "bca",
          "sender_name": "ahmad reza musthafa",
          "status": 10,
          "total_amount": 5000150,
          "unique_code": 150,
          "user_id": 1,
          "virtual_account_number": "8558151502502621"
        },
        "string": "Journal no: 3453456 TRANSFER DARI Bpk AHMAD REZA MUSTHAFA"
      }
    }
  ]
}
```

In the response:
- `id`: The document ID.
- `score`: A relevance score representing how closely the document matches the search terms.
- `data`: The actual document content.

---

## Installation
### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/ahmadrezamusthafa/search-engine.git
   cd search-engine
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Running the Service:
   - The server will run on `http://localhost:9000`
   - Web UI available on `http://localhost:9000/web`
