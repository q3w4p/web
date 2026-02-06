import { motion } from "framer-motion";
import { Link } from "wouter";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { Button } from "@/components/ui/button";
import { useUser } from "@/hooks/use-auth";
import { SiDiscord } from "react-icons/si";
import { useEffect, useState } from "react";

export default function Home() {
  const { data: user } = useUser();
  const [activeStat, setActiveStat] = useState(0);

  const stats = [
    { label: "Active Bots", value: "2,400+" },
    { label: "Total Users", value: "150k+" },
    { label: "Uptime", value: "99.99%" },
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
            <div className="w-20 h-20 rounded-3xl overflow-hidden border border-white/10 mx-auto shadow-2xl shadow-black mb-6 group">
              <img 
                src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg?ex=69867290&is=69852110&hm=27d1f6242d0aa1a2d312aab51b6f55eb30fc2c3d9a9894a572e3e4b80bd00379" 
                className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                alt="Hurry"
              />
            </div>
            <h1 className="text-6xl md:text-8xl font-bold font-display tracking-tighter mb-4">
              <span className="text-gradient">Hurry</span>
            </h1>
            <p className="text-lg text-white/40 font-medium max-w-lg mx-auto leading-relaxed"></p>
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
