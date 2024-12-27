# Build

docker build -t payroll-checker-app:1.0.0-SNAPSHOT .

# Run

docker run -e GEMINI_API_KEY=XXXX -p 8080:8080 payroll-checker-app:1.0.0-SNAPSHOT