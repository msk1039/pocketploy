"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Database, Shield, Zap, ArrowRight } from "lucide-react";
import { useEffect } from "react";

export default function Home() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  // Redirect to dashboard if already authenticated
  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.push("/dashboard");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Header */}
      <header className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex h-16 items-center justify-between">
            <div className="flex items-center space-x-2">
              <Database className="h-6 w-6 text-blue-600" />
              <span className="text-xl font-bold text-gray-900">pocketploy</span>
            </div>
            <div className="flex items-center space-x-4">
              <Link href="/login">
                <Button variant="ghost">Log in</Button>
              </Link>
              <Link href="/signup">
                <Button>Sign up</Button>
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="py-20 sm:py-32">
          <div className="text-center">
            <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
              Lightweight{" "}
              <span className="text-blue-600">PocketBase</span>{" "}
              Hosting
            </h1>
            <p className="mt-6 text-lg leading-8 text-gray-600 max-w-2xl mx-auto">
              Create and manage isolated PocketBase instances dynamically.
              Get your own database up and running in seconds.
            </p>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Link href="/signup">
                <Button size="lg" className="group">
                  Get started
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
              <Link href="/login">
                <Button size="lg" variant="outline">
                  Log in
                </Button>
              </Link>
            </div>
          </div>

          {/* Features */}
          <div className="mt-20 grid grid-cols-1 gap-8 sm:grid-cols-3">
            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center space-y-4">
                  <div className="rounded-full bg-blue-100 p-3">
                    <Zap className="h-6 w-6 text-blue-600" />
                  </div>
                  <h3 className="font-semibold text-lg">Lightning Fast</h3>
                  <p className="text-sm text-gray-600">
                    Deploy your PocketBase instance in seconds with our streamlined process
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center space-y-4">
                  <div className="rounded-full bg-green-100 p-3">
                    <Shield className="h-6 w-6 text-green-600" />
                  </div>
                  <h3 className="font-semibold text-lg">Secure & Isolated</h3>
                  <p className="text-sm text-gray-600">
                    Each instance runs in its own isolated environment with JWT authentication
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center space-y-4">
                  <div className="rounded-full bg-purple-100 p-3">
                    <Database className="h-6 w-6 text-purple-600" />
                  </div>
                  <h3 className="font-semibold text-lg">Easy Management</h3>
                  <p className="text-sm text-gray-600">
                    Manage all your PocketBase instances from a single dashboard
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* CTA Section */}
          <div className="mt-20 rounded-2xl bg-blue-600 px-6 py-16 text-center">
            <h2 className="text-3xl font-bold text-white">
              Ready to get started?
            </h2>
            <p className="mt-4 text-lg text-blue-100">
              Create your free account and deploy your first PocketBase instance today
            </p>
            <div className="mt-8">
              <Link href="/signup">
                <Button size="lg" variant="secondary" className="group">
                  Create account
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t mt-20">
        <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
          <p className="text-center text-sm text-gray-600">
            © 2025 pocketploy. Phase 1 - Authentication Complete.
          </p>
        </div>
      </footer>
    </div>
  );
}
