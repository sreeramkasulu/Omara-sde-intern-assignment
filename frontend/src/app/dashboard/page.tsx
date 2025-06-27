"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Upload, FileText, MessageCircle, Trash2 } from "lucide-react"

interface Document {
  id: string
  user_id: string
  file_name: string
  storage_path: string
  uploaded_at: string
}

interface ChatMessage {
  id: string
  document_id: string
  user_id: string
  message_type: "user" | "ai"
  message_content: string
  timestamp: string
}

const API_BASE_URL = "https://8080-ipkf1nf88rwalpt7bhibe-8212350b.manusvm.computer/api"

export default function Dashboard() {
  const [documents, setDocuments] = useState<Document[]>([])
  const [selectedDocument, setSelectedDocument] = useState<Document | null>(null)
  const [chatHistory, setChatHistory] = useState<ChatMessage[]>([])
  const [query, setQuery] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [userId] = useState("ddf650f5-2147-4e2a-9fdc-524c0d321994") // Using the created user ID

  useEffect(() => {
    fetchDocuments()
  }, [])

  useEffect(() => {
    if (selectedDocument) {
      fetchChatHistory(selectedDocument.id)
    }
  }, [selectedDocument])

  const fetchDocuments = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/documents?user_id=${userId}`)
      if (response.ok) {
        const docs = await response.json()
        setDocuments(docs || [])
      }
    } catch (error) {
      console.error("Failed to fetch documents:", error)
    }
  }

  const fetchChatHistory = async (documentId: string) => {
    try {
      const response = await fetch(`${API_BASE_URL}/documents/${documentId}/chat-history`)
      if (response.ok) {
        const history = await response.json()
        setChatHistory(history || [])
      }
    } catch (error) {
      console.error("Failed to fetch chat history:", error)
    }
  }

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append("file", file)
    formData.append("user_id", userId)

    try {
      setIsLoading(true)
      const response = await fetch(`${API_BASE_URL}/documents/upload`, {
        method: "POST",
        body: formData,
      })

      if (response.ok) {
        await fetchDocuments()
        event.target.value = "" // Reset file input
      } else {
        alert("Failed to upload document")
      }
    } catch (error) {
      console.error("Upload error:", error)
      alert("Failed to upload document")
    } finally {
      setIsLoading(false)
    }
  }

  const handleAnalyze = async () => {
    if (!selectedDocument || !query.trim()) return

    try {
      setIsLoading(true)
      const response = await fetch(
        `${API_BASE_URL}/documents/${selectedDocument.id}/analyze?user_id=${userId}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ query }),
        }
      )

      if (response.ok) {
        await fetchChatHistory(selectedDocument.id)
        setQuery("")
      } else {
        alert("Failed to analyze document")
      }
    } catch (error) {
      console.error("Analysis error:", error)
      alert("Failed to analyze document")
    } finally {
      setIsLoading(false)
    }
  }

  const handleDeleteDocument = async (documentId: string) => {
    if (!confirm("Are you sure you want to delete this document?")) return

    try {
      const response = await fetch(`${API_BASE_URL}/documents/${documentId}`, {
        method: "DELETE",
      })

      if (response.ok) {
        await fetchDocuments()
        if (selectedDocument?.id === documentId) {
          setSelectedDocument(null)
          setChatHistory([])
        }
      } else {
        alert("Failed to delete document")
      }
    } catch (error) {
      console.error("Delete error:", error)
      alert("Failed to delete document")
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Strategic Insight Dashboard</h1>
          <p className="text-gray-600">Upload documents and analyze them with AI-powered insights</p>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Document Management */}
          <div className="lg:col-span-1">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  Documents
                </CardTitle>
                <CardDescription>Upload and manage your business documents</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="file-upload">Upload Document</Label>
                  <Input
                    id="file-upload"
                    type="file"
                    accept=".pdf,.txt"
                    onChange={handleFileUpload}
                    disabled={isLoading}
                    className="mt-1"
                  />
                  <p className="text-sm text-gray-500 mt-1">PDF and TXT files only</p>
                </div>

                <div className="space-y-2">
                  <Label>Your Documents</Label>
                  {documents.length === 0 ? (
                    <p className="text-sm text-gray-500">No documents uploaded yet</p>
                  ) : (
                    <div className="space-y-2">
                      {documents.map((doc) => (
                        <div
                          key={doc.id}
                          className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                            selectedDocument?.id === doc.id
                              ? "border-blue-500 bg-blue-50"
                              : "border-gray-200 hover:border-gray-300"
                          }`}
                          onClick={() => setSelectedDocument(doc)}
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex-1 min-w-0">
                              <p className="text-sm font-medium text-gray-900 truncate">
                                {doc.file_name}
                              </p>
                              <p className="text-xs text-gray-500">
                                {new Date(doc.uploaded_at).toLocaleDateString()}
                              </p>
                            </div>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={(e) => {
                                e.stopPropagation()
                                handleDeleteDocument(doc.id)
                              }}
                              className="text-red-500 hover:text-red-700"
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Chat Interface */}
          <div className="lg:col-span-2">
            <Card className="h-[600px] flex flex-col">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MessageCircle className="h-5 w-5" />
                  AI Analysis Chat
                </CardTitle>
                <CardDescription>
                  {selectedDocument
                    ? `Analyzing: ${selectedDocument.file_name}`
                    : "Select a document to start analyzing"}
                </CardDescription>
              </CardHeader>
              <CardContent className="flex-1 flex flex-col">
                {!selectedDocument ? (
                  <div className="flex-1 flex items-center justify-center text-gray-500">
                    <div className="text-center">
                      <Upload className="h-12 w-12 mx-auto mb-4 text-gray-300" />
                      <p>Select a document from the left panel to start analyzing</p>
                    </div>
                  </div>
                ) : (
                  <>
                    {/* Chat History */}
                    <div className="flex-1 overflow-y-auto space-y-4 mb-4 p-4 bg-gray-50 rounded-lg">
                      {chatHistory.length === 0 ? (
                        <p className="text-gray-500 text-center">
                          No conversation yet. Ask a question about this document!
                        </p>
                      ) : (
                        chatHistory.map((message) => (
                          <div
                            key={message.id}
                            className={`flex ${
                              message.message_type === "user" ? "justify-end" : "justify-start"
                            }`}
                          >
                            <div
                              className={`max-w-[80%] p-3 rounded-lg ${
                                message.message_type === "user"
                                  ? "bg-blue-500 text-white"
                                  : "bg-white border"
                              }`}
                            >
                              <p className="text-sm">{message.message_content}</p>
                              <p
                                className={`text-xs mt-1 ${
                                  message.message_type === "user"
                                    ? "text-blue-100"
                                    : "text-gray-500"
                                }`}
                              >
                                {new Date(message.timestamp).toLocaleTimeString()}
                              </p>
                            </div>
                          </div>
                        ))
                      )}
                    </div>

                    {/* Query Input */}
                    <div className="space-y-2">
                      <Label htmlFor="query">Ask a question about this document</Label>
                      <div className="flex gap-2">
                        <Textarea
                          id="query"
                          placeholder="e.g., Summarize the key strategic initiatives mentioned in this report..."
                          value={query}
                          onChange={(e) => setQuery(e.target.value)}
                          disabled={isLoading}
                          className="flex-1"
                          rows={2}
                        />
                        <Button
                          onClick={handleAnalyze}
                          disabled={isLoading || !query.trim()}
                          className="self-end"
                        >
                          {isLoading ? "Analyzing..." : "Analyze"}
                        </Button>
                      </div>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}

