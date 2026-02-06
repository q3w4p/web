import { useEffect } from "react";
import { useLocation } from "wouter";
import { useUser } from "@/hooks/use-auth";
import { useBots } from "@/hooks/use-bots";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { BotCard } from "@/components/BotCard";
import { CreateBotDialog } from "@/components/CreateBotDialog";
import { Loader2, ServerCrash, Bot } from "lucide-react";
import { motion } from "framer-motion";

export default function Dashboard() {
  const [_, setLocation] = useLocation();
  const { data: user, isLoading: userLoading } = useUser();
  const { data: bots, isLoading: botsLoading, error: botsError } = useBots();

  useEffect(() => {
    if (!userLoading && !user) {
      setLocation("/");
    }
  }, [user, userLoading, setLocation]);

  if (userLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="w-8 h-8 text-violet-500 animate-spin" />
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="min-h-screen flex flex-col bg-fluid text-white relative">
      <Navbar />

      <main className="flex-1 max-w-7xl mx-auto w-full px-6 py-24 relative z-10">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-12">
          <div>
            <h1 className="text-3xl font-bold font-display mb-2">My Bots</h1>
            <p className="text-white/60">Manage your deployed instances and monitor status.</p>
          </div>
          <CreateBotDialog />
        </div>

        {botsLoading ? (
          <div className="flex flex-col items-center justify-center py-20 text-white/40">
            <Loader2 className="w-10 h-10 animate-spin mb-4" />
            <p>Loading your bots...</p>
          </div>
        ) : botsError ? (
          <div className="glass-card p-8 flex flex-col items-center justify-center text-center">
            <ServerCrash className="w-12 h-12 text-red-400 mb-4" />
            <h3 className="text-xl font-bold mb-2">Failed to load bots</h3>
            <p className="text-white/50 mb-6">Something went wrong while fetching your data.</p>
            <button onClick={() => window.location.reload()} className="text-violet-400 hover:underline">
              Try refreshing
            </button>
          </div>
        ) : bots?.length === 0 ? (
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="glass-card p-12 flex flex-col items-center justify-center text-center min-h-[400px]"
          >
            <div className="w-20 h-20 rounded-2xl bg-white/5 flex items-center justify-center mb-6">
              <Bot className="w-10 h-10 text-white/20" />
            </div>
            <h3 className="text-2xl font-bold mb-3">No bots deployed yet</h3>
            <p className="text-white/50 max-w-md mb-8">
              You haven't hosted any bots yet. Click the button above to deploy your first Discord bot instance.
            </p>
            <CreateBotDialog />
          </motion.div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {bots?.map((bot) => (
              <BotCard key={bot.id} bot={bot} />
            ))}
          </div>
        )}
      </main>

      <Footer />
    </div>
  );
}
