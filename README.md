# Strategic Insight Analyst

## Overview

This project is a **Strategic Insight Analyst** web application that allows users to upload business documents such as annual reports, press releases, and market analyses. The application leverages AI models to extract, summarize, and compare strategic insights from these documents.

The project was developed as part of an SDE Internship Assignment, and I have successfully implemented both the backend and frontend for document processing and analysis.

## Features

- **User Authentication (Email/Password Login):** Users can log in via email and password using Firebase Authentication.
- **Document Upload:** Users can upload PDF or `.txt` files for analysis.
- **AI-Powered Insights:** The system uses a large language model (LLM) to provide insights from uploaded documents.
- **Interactive Interface:** Users can interact with the system to ask specific strategic queries, which the LLM will respond to based on the uploaded documents.

## Task Requirements

The task was to design, implement, and deploy a full-stack application that meets the following requirements:

- **User Authentication** with secure login (Firebase Authentication).
- **Document Management** including uploading, storage, and listing of documents.
- **Document Processing** for text extraction from PDFs and `.txt` files.
- **AI Integration** to generate strategic insights and answers based on uploaded documents.
- **Deployment** of both frontend and backend to publicly accessible URLs.

## Tech Stack

- **Frontend:** Next.js with ShadCN UI for components.
- **Backend:** GoLang for the server, Firebase for authentication, and Google Cloud Storage for document file storage.
- **Database:** SQL-based database (PostgreSQL).
- **LLM:** Google Gemini for generating insights from documents.

## Deployment

Both the frontend and backend of the application are deployed and accessible at the following live URL:

- **Live Application:** [Strategic Insight Analyst Dashboard](https://smwicopi.manus.space/)

## Project Satisfaction

I am satisfied with the implementation and completion of the project. I have met all of the core requirements, including secure user authentication, document management, AI-powered insight generation, and successful deployment of the application.

## How to Run the Project Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/sreeramkasulu/Omara-sde-intern-assignment.git
   ```

2. Install dependencies:
   - For the frontend:
     ```bash
     cd frontend
     npm install
     ```
   - For the backend:
     ```bash
     cd backend
     go mod tidy
     ```

3. Set up environment variables for both the frontend and backend (e.g., API keys, Firebase credentials, database URLs).

4. Run both frontend and backend servers:
   - Frontend: `npm run dev`
   - Backend: `go run main.go`

5. Open the application in your browser at `http://localhost:3000`.

## License

This project is licensed under the MIT License.

