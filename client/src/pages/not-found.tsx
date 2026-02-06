import { Link } from "wouter";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { AlertTriangle } from "lucide-react";

export default function NotFound() {
  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-fluid text-white relative">
      <Card className="w-full max-w-md mx-4 glass-card border-white/10">
        <CardContent className="pt-6">
          <div className="flex mb-4 gap-2 text-red-400">
            <AlertTriangle className="h-8 w-8" />
            <h1 className="text-2xl font-bold font-display text-white">404 Page Not Found</h1>
          </div>

          <p className="mt-4 text-white/60 leading-relaxed">
            The page you are looking for does not exist. It might have been moved or deleted.
          </p>

          <Link href="/">
            <Button className="w-full mt-6 bg-white text-black hover:bg-white/90">
              Return Home
            </Button>
          </Link>
        </CardContent>
      </Card>
    </div>
  );
}
