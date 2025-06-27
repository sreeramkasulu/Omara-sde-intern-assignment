import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import Link from "next/link"

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-bold text-gray-900 mb-6">
            Strategic Insight Analyst
          </h1>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Upload business documents and leverage AI to extract, summarize, and compare strategic insights. 
            Transform your documents into actionable intelligence.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <Card className="text-center">
            <CardHeader>
              <CardTitle className="text-2xl text-blue-600">📄 Upload Documents</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-lg">
                Securely upload PDF and text files containing business information, 
                annual reports, and market analyses.
              </CardDescription>
            </CardContent>
          </Card>

          <Card className="text-center">
            <CardHeader>
              <CardTitle className="text-2xl text-green-600">🤖 AI Analysis</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-lg">
                Ask specific questions about your documents and get strategic insights 
                powered by advanced AI models.
              </CardDescription>
            </CardContent>
          </Card>

          <Card className="text-center">
            <CardHeader>
              <CardTitle className="text-2xl text-purple-600">📊 Compare & Summarize</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription className="text-lg">
                Compare multiple documents, identify trends, and generate 
                comprehensive strategic summaries.
              </CardDescription>
            </CardContent>
          </Card>
        </div>

        <div className="text-center">
          <Link href="/dashboard">
            <Button size="lg" className="text-lg px-8 py-4">
              Get Started
            </Button>
          </Link>
        </div>

        <div className="mt-16 text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-8">Key Features</h2>
          <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
            <div className="text-left">
              <h3 className="text-xl font-semibold mb-4">🔐 Secure Authentication</h3>
              <p className="text-gray-600">
                Email/password login system ensures your documents remain private and secure.
              </p>
            </div>
            <div className="text-left">
              <h3 className="text-xl font-semibold mb-4">💬 Interactive Chat</h3>
              <p className="text-gray-600">
                Chat-like interface for asking questions and getting detailed insights about your documents.
              </p>
            </div>
            <div className="text-left">
              <h3 className="text-xl font-semibold mb-4">📈 Strategic Analysis</h3>
              <p className="text-gray-600">
                Extract key strategic initiatives, financial highlights, and competitive advantages.
              </p>
            </div>
            <div className="text-left">
              <h3 className="text-xl font-semibold mb-4">🔄 Document Comparison</h3>
              <p className="text-gray-600">
                Compare multiple documents to identify patterns and strategic differences.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

