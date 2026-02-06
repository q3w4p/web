import { motion } from "framer-motion";
import { Link } from "wouter";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Button } from "@/components/ui/button";
import { useUser } from "@/hooks/use-auth";
import { SiDiscord } from "react-icons/si";
import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";

export default function Home() {
  const { data: user } = useUser();
  const [activeStat, setActiveStat] = useState(0);

  const { data: statsData } = useQuery<{ activeBots: number, totalUsers: number, uptime: string }>({
    queryKey: ["/api/stats"],
  });

  const stats = [
    { label: "Active Bots", value: statsData ? `${statsData.activeBots}+` : "..." },
    { label: "Total Users", value: statsData ? `${statsData.totalUsers}` : "..." },
    { label: "Uptime", value: statsData?.uptime || "99.9%" },
  ];

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveStat((prev) => (prev + 1) % stats.length);
    }, 2000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen flex flex-col relative overflow-hidden text-white bg-fluid">
      <Navbar />

      {/* Hero Section */}
      <main className="flex-1 flex items-center justify-center relative pt-20">
        <div className="max-w-4xl mx-auto px-6 py-20 text-center relative z-10 load-up">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, ease: "easeOut" }}
            className="mb-10"
          >
            <div className="w-24 h-24 rounded-3xl overflow-hidden border border-white/10 mx-auto shadow-2xl shadow-black mb-6 group">
              <img 
                src="https://cdn.discordapp.com/attachments/1461273269058666619/1469322858755788891/1.jpeg?ex=69873d0c&is=6985eb8c&hm=3b9cd281dcaa3dc6be224f2ca083f0131e0d512db11f50f1d65ec8c34011bc48" 
                className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                alt="Hurry"
              />
            </div>
            <h1 className="text-6xl md:text-8xl font-bold font-display tracking-tighter mb-4">
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-500 to-red-500 animate-gradient-x">Hurry</span>
            </h1>
            <p className="text-lg text-white/60 font-medium max-w-lg mx-auto leading-relaxed">
              premium hosting for discord
            </p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="flex flex-col items-center gap-12"
          >
            {user ? (
              <Link href="/dashboard">
                <Button className="rival-button rival-button-primary h-14 px-10 text-lg rounded-xl font-bold">
                  Dashboard
                </Button>
              </Link>
            ) : (
              <a href="/api/auth/discord">
                <Button className="rival-button rival-button-discord h-14 px-10 text-lg rounded-xl font-bold flex items-center gap-3">
                  <SiDiscord className="w-5 h-5" />
                  Add to Discord
                </Button>
              </a>
            )}

            <div className="flex flex-wrap items-center justify-center gap-10 md:gap-16">
              {stats.map((stat, i) => (
                <div
                  key={i}
                  className="flex flex-col items-center cursor-default"
                >
                  <span className={`text-3xl font-bold font-display mb-1 stat-underline ${activeStat === i ? 'active' : ''}`}>{stat.value}</span>
                  <span className="text-[10px] uppercase tracking-widest text-white/20 font-black">{stat.label}</span>
                </div>
              ))}
            </div>
          </motion.div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
