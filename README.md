# My Backend

## Deployment

### Prerequisites
- Install Go, Docker, and Docker Compose.
- Ensure your local machine has ports 8080, 5432, and 6379 free. If not, modify the ports in the `docker-compose.yml` file.

### Steps

1. **Start PostgreSQL and Redis**
   ```sh
   docker-compose up
   ```

2. **Build Docker Image**
   ```sh
   docker build -t backend:1.0.0 .
   ```

3. **Migrate Database and Run Application**
   ```sh
   docker build -t backend:1.0.0 .
   docker run backend:1.0.0 migrate
   docker run -d -p 8080:8080 backend:1.0.0 api
   ```

4. **Run Tests**
    - I'm using `mockery` to generate mocks in the `tests` folder (use the `mock.sh` script to generate mocks).
    - I've added a command to run test cases in the Dockerfile (`RUN go test ./... -v`), but it's commented out because it takes a lot of time to run and I don't want to run it every time I build the image.
    - I think it should be done locally (on your local machine) or in a CI/CD pipeline (for real-world applications).

5. **Import Postman Collection**
    - Import the Postman collection using the JSON file in the `docs` folder.

6. **Get JWT Token**
    - Use the `users/login` endpoint to get a JWT token.

7. **Generate Test Data**
    - Use the `products/generate-test-data` endpoint to generate test data.
    - It will run in the background and take about 3 minutes.
    - Check the progress in the container logs or rerun the API to know if the data generation is in progress (look for "Another test data generation process is running" in the log).

## Explanation

This application is a simple REST API that allows you to manage products, users, categories, and visualize statistics as well. Below is an overview of its design and implementation:

- **Technologies Used**:
    - **Cobra**: To create a command-line application that can run `migrate` and `api` commands separately.
    - **Docker**: To package the application to run anywhere that has Docker installed.
    - **Onion Architecture**: Separates the application into layers. It's easy to test and maintain when the codebase becomes large but hard to understand for newcomers at the beginning.
    - **Tech Stack**:
        - **PostgreSQL**: Database.
        - **Go**: Programming language.
        - **Fiber**: HTTP server.
        - **Go Bun**: ORM.
        - **Zerolog**: Logging.
        - **Redis**: Cache and distributed lock.
        - **JWT**: Authentication.

- **Database Schema Changes**:
    - Changed `product-category` to a 1-1 relationship by adding `category_id` to the `products` table. Removed the `product_category` relation table.
    - Created an index for the `product_id` column in the `products` table to get dashboard data faster.
    - Added some indexes for the `products` table to search faster.

- **Features Implemented**:
    - Due to time constraints, I focused on implementing product features, the dashboard (the hardest parts) and login features.
    - Didn't implement any validation for the request. It can be done easy with Validator package with tag.
    - Other features like CRUD for users, wishlists, and categories were not implemented.
    - **Future Plans**:
        - Implement the login feature using `argon2` to hash passwords.
        - To implement role-based access control, after login, I will generate the correct JWT based on the user’s role.
        - Update authorization to check whether the user is an admin or not. Currently, I just check if the user is logged in or not based on JWT expiry time.

## Real-time Data Processing (Focus on real-time problem solving)

- **Dashboard Feature**:
    - Implemented a dashboard that shows the number of products in each category already.
    - Created an index for the `category_id` column in the `products` table to get dashboard data faster, then cached the result in Redis (in a single key).
    - Whenever a new product is created, I delete the cache.
    - **Potential Improvement**:
        - It can be improved by storing four keys for each category and updating the cache instead of deleting it.

## Log Aggregation Optimization (Focus on efficient log handling)

- Using `zerolog` for logging. The configuration is defined in the `setupLogger` function in `root.go`. It overrides the global logger from `zerolog`.
- Can integrate with third-party monitoring tools like Sentry (without a separate adapter and call to Sentry for each failure).
- I've written a function to send logs to Sentry in case of failure.

## Product Search Optimization (Focus on data search performance)

1. **Test Data Generation**:
    - To generate test data, I'm using a worker pool and sending dummy data to workers. After exceeding the batch size, I insert them into the database.
    - I must set a small batch size (1000) to avoid timeout errors when inserting into PostgreSQL.
    - **Potential Improvements**:
        - It can be improved by using a task queue like Go Machinery background worker to process the data in the background to achieve fault tolerance.
        - I also thought about using a `COPY` query to get the data from a CSV file and insert them into the `products` table. MAYBE it can be faster than using batch insert.

2. **Indexes Created**:
    - For `description` and `name` columns, I'm using GIN indexes to search faster because of `ILIKE` (case-insensitive):
      ```sql
      CREATE EXTENSION IF NOT EXISTS pg_trgm;
      CREATE INDEX products_description_idx ON products USING GIN (description gin_trgm_ops);
      CREATE INDEX products_name_idx ON products USING GIN (name gin_trgm_ops);
      CREATE INDEX products_created_at_idx ON products (created_at);
      CREATE INDEX products_category_id_idx ON products (category_id);
      CREATE INDEX products_status_idx ON products (status);
      ```

3. **Performance**:
    - With 10 million records, my `GET /products` endpoint takes about 1 seconds to get the data with pagination (I've set the default page size to 20).