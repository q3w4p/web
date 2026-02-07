import { Switch, Route } from "wouter";
import { queryClient } from "./lib/queryClient";
import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/toaster";
import { TooltipProvider } from "@/components/ui/tooltip";
import Home from "@/pages/Home";
import Dashboard from "@/pages/Dashboard";
import Admin from "@/pages/Admin";
import GetToken from "@/pages/GetToken";
import NotFound from "@/pages/not-found";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { useQuery } from "@tanstack/react-query";
import { User } from "@shared/schema";
import { Button } from "@/components/ui/button";
import { LogOut, User as UserIcon, Menu } from "lucide-react";
import { apiRequest } from "./lib/queryClient";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";

import { useUser } from "@/hooks/use-auth";
import { Redirect } from "wouter";

function ProtectedRoute({ component: Component, adminOnly = false }: { component: React.ComponentType, adminOnly?: boolean }) {
  const { data: user, isLoading } = useUser();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
      </div>
    );
  }

  if (!user || user.id === 0) {
    return <Redirect to="/" />;
  }

  if (adminOnly && !user.isAdmin) {
    return <Redirect to="/" />;
  }

  return <Component />;
}

function Router() {
  return (
    <div className="min-h-screen flex flex-col relative overflow-hidden text-white bg-transparent">
      <Navbar />
      <main className="flex-1 relative pt-24 pb-12 overflow-y-auto">
        <Switch>
          <Route path="/" component={Home} />
          <Route path="/dashboard">
            <ProtectedRoute component={Dashboard} />
          </Route>
          <Route path="/admin">
            <ProtectedRoute component={Admin} adminOnly />
          </Route>
          <Route path="/get-token">
            <ProtectedRoute component={GetToken} />
          </Route>
          <Route component={NotFound} />
        </Switch>
      </main>
      <Footer />
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <Router />
        <Toaster />
      </TooltipProvider>
    </QueryClientProvider>
  );
}

export default App;
