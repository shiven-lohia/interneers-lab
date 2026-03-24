# Prompt To Recreate Project Context (Backend/Go)

Copy-paste this prompt into a new ChatGPT chat:

---

I lost my previous project guidance chat. Please help me continue my backend work and first rebuild full context from this summary.

Project scope for this chat:
- Focus only on backend/go in my interneers-lab repository.
- Ignore frontend and python backend unless I explicitly ask.

My setup:
- I am using Copilot Pro in VS Code.
- I run MongoDB using Docker.
- I test APIs with Postman collections and curl.
- I also verify data in MongoDB Compass.

Important instruction:
- Please read the PDF I attached for the week-by-week plan and align your guidance with it.

What I have already done in backend/go:
- Set up Go backend service structure with modules for hello world and products.
- Added route registration with net/http ServeMux.
- Added request logging middleware using zerolog.
- Built Products CRUD APIs:
  - GET /products
  - GET /products/{id}
  - POST /products
  - PUT /products/{id}
  - DELETE /products/{id}
- Added product validation in controller:
  - name required
  - price > 0
  - quantity >= 0
- Started with in-memory map repository, then migrated to MongoDB repository implementation.
- Integrated MongoDB connection using mongo.Connect and client.Ping.
- Introduced repository abstraction so controller/handler logic can stay decoupled from storage implementation.
- Added context propagation from handler -> controller -> repository.
- Added timeout handling in repository methods using context.WithTimeout.
- Improved MongoDB error handling:
  - Update checks MatchedCount
  - Delete checks DeletedCount
- Ensured API behavior improvements:
  - GET all returns [] when no records exist
  - DELETE success returns 204
  - Missing resources return 404
- Added bulk product create from CSV upload endpoint:
  - POST /products/bulk (multipart file)
  - Parses CSV and creates valid product rows
- Fixed duplicate-ID issue in create flow:
  - Before insert, checks existing product by id
  - Returns "product with this id already exists" for duplicates
  - Also handles Mongo duplicate key errors defensively

How I want you to work now:
1. First, restate your understanding of the current backend/go status and week alignment from the attached PDF.
2. Tell me if anything in my summary looks inconsistent or risky.
3. Then ask me how I want to proceed next.

Do not jump into random changes before confirming understanding.

---
