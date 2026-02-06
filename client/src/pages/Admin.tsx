import { useEffect } from "react";
import { useLocation } from "wouter";
import { useUser } from "@/hooks/use-auth";
import { useAdminUsers, useAdminBots } from "@/hooks/use-admin";
import { Navbar } from "@/components/Navbar";
import { Loader2, Users, Bot, ShieldAlert } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { motion } from "framer-motion";

export default function Admin() {
  const [_, setLocation] = useLocation();
  const { data: user, isLoading: userLoading } = useUser();
  const { data: allUsers, isLoading: usersLoading } = useAdminUsers();
  const { data: allBots, isLoading: botsLoading } = useAdminBots();

  useEffect(() => {
    if (!userLoading) {
      if (!user) {
        setLocation("/");
      } else if (!user.isAdmin) {
        setLocation("/dashboard");
      }
    }
  }, [user, userLoading, setLocation]);

  if (userLoading || usersLoading || botsLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center text-white">
        <Loader2 className="w-8 h-8 text-violet-500 animate-spin" />
      </div>
    );
  }

  if (!user?.isAdmin) return null;

  const stats = [
    { label: "Total Users", value: allUsers?.length || 0, icon: Users, color: "text-blue-400" },
    { label: "Active Bots", value: allBots?.length || 0, icon: Bot, color: "text-violet-400" },
    { label: "Online Instances", value: allBots?.filter(b => b.status === 'online').length || 0, icon: ShieldAlert, color: "text-green-400" },
  ];

  return (
    <div className="min-h-screen flex flex-col bg-fluid text-white relative">
      <Navbar />

      <main className="flex-1 max-w-7xl mx-auto w-full px-6 py-24 relative z-10">
        <div className="mb-12">
          <h1 className="text-3xl font-bold font-display mb-2">Admin Dashboard</h1>
          <p className="text-white/60">System overview and user management.</p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          {stats.map((stat, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.1 }}
              className="glass-card p-6 flex items-center justify-between"
            >
              <div>
                <p className="text-sm font-medium text-white/50 mb-1">{stat.label}</p>
                <p className="text-3xl font-bold">{stat.value}</p>
              </div>
              <div className={`p-4 rounded-xl bg-white/5 ${stat.color}`}>
                <stat.icon className="w-6 h-6" />
              </div>
            </motion.div>
          ))}
        </div>

        {/* Users Table */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="glass-card overflow-hidden"
        >
          <div className="p-6 border-b border-white/10">
            <h2 className="text-xl font-bold">User Directory</h2>
          </div>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader className="bg-white/5 hover:bg-white/5">
                <TableRow className="border-white/10 hover:bg-transparent">
                  <TableHead className="text-white/60">User</TableHead>
                  <TableHead className="text-white/60">Discord ID</TableHead>
                  <TableHead className="text-white/60">Joined</TableHead>
                  <TableHead className="text-white/60 text-right">Role</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {allUsers?.map((u) => (
                  <TableRow key={u.id} className="border-white/10 hover:bg-white/5 transition-colors">
                    <TableCell className="font-medium text-white">
                      <div className="flex items-center gap-3">
                        {u.avatar ? (
                          <img 
                            src={`https://cdn.discordapp.com/avatars/${u.discordId}/${u.avatar}.png`}
                            alt=""
                            className="w-8 h-8 rounded-full"
                          />
                        ) : (
                          <div className="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center">
                            <span className="text-xs">{u.username[0]}</span>
                          </div>
                        )}
                        {u.username}
                      </div>
                    </TableCell>
                    <TableCell className="text-white/60 font-mono text-xs">{u.discordId}</TableCell>
                    <TableCell className="text-white/60">
                      {new Date(u.createdAt || '').toLocaleDateString()}
                    </TableCell>
                    <TableCell className="text-right">
                      {u.isAdmin ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-violet-500/20 text-violet-300">
                          Admin
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-white/10 text-white/60">
                          User
                        </span>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </motion.div>
      </main>
    </div>
  );
}
