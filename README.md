# Build

docker build -t payroll-checker-app:1.0.0-SNAPSHOT .

# Run

docker run -e GEMINI_API_KEY=XXXX payroll-checker-app:1.0.0-SNAPSHOT