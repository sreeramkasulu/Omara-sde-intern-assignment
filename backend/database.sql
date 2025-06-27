-- Strategic Insight Analyst Database Schema

-- Users Table
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY, -- Firebase Auth UID
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Documents Table
CREATE TABLE documents (
    id VARCHAR(255) PRIMARY KEY, -- Unique ID for the document
    user_id VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    storage_path VARCHAR(255) NOT NULL, -- Path to the original file in GCS
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Document Chunks Table (for LLM context)
CREATE TABLE document_chunks (
    id VARCHAR(255) PRIMARY KEY, -- Unique ID for the chunk
    document_id VARCHAR(255) NOT NULL,
    chunk_index INT NOT NULL, -- Order of the chunk within the document
    content TEXT NOT NULL,
    embedding JSONB, -- Optional: For future semantic search/retrieval
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    UNIQUE (document_id, chunk_index) -- Ensures unique ordering of chunks per document
);

-- Chat History Table (for user-LLM interactions per document)
CREATE TABLE chat_history (
    id VARCHAR(255) PRIMARY KEY, -- Unique ID for the chat message
    document_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    message_type VARCHAR(4) NOT NULL CHECK (message_type IN ('user', 'ai')), -- 'user' for query, 'ai' for response
    message_content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_document_chunks_document_id ON document_chunks(document_id);
CREATE INDEX idx_chat_history_document_id ON chat_history(document_id);
CREATE INDEX idx_chat_history_user_id ON chat_history(user_id);

