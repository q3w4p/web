import { motion } from "framer-motion";
import { Bot, Power, Trash2, Activity, Terminal } from "lucide-react";
import { Button } from "@/components/ui/button";
import { type Bot as BotType } from "@shared/schema";
import { useDeleteBot, useBotAction } from "@/hooks/use-bots";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

interface BotCardProps {
  bot: BotType;
}

export function BotCard({ bot }: BotCardProps) {
  const deleteBot = useDeleteBot();
  const botAction = useBotAction();

  const isOnline = bot.status === "online";
  const isPending = botAction.isPending || deleteBot.isPending;

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="glass-card p-6 relative group overflow-hidden"
    >
      {/* Background glow effect */}
      <div className={`absolute top-0 right-0 w-32 h-32 bg-gradient-to-br ${isOnline ? 'from-green-500/20' : 'from-red-500/20'} to-transparent blur-3xl rounded-full -mr-16 -mt-16 transition-colors duration-500`} />

      <div className="flex justify-between items-start mb-6 relative z-10">
        <div className="flex items-center gap-3">
          <div className={`p-3 rounded-xl ${isOnline ? 'bg-green-500/20 text-green-400' : 'bg-white/5 text-white/60'}`}>
            <Bot className="w-6 h-6" />
          </div>
          <div>
            <h3 className="font-bold text-lg text-white">{bot.name}</h3>
            <div className="flex items-center gap-2 mt-1">
              <span className={`relative flex h-2 w-2`}>
                <span className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${isOnline ? 'bg-green-400' : 'bg-red-400'}`}></span>
                <span className={`relative inline-flex rounded-full h-2 w-2 ${isOnline ? 'bg-green-500' : 'bg-red-500'}`}></span>
              </span>
              <span className={`text-xs font-medium uppercase tracking-wider ${isOnline ? 'text-green-400' : 'text-white/40'}`}>
                {bot.status}
              </span>
            </div>
          </div>
        </div>

        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button 
              variant="ghost" 
              size="icon" 
              className="text-white/40 hover:text-red-400 hover:bg-red-500/10"
              disabled={isPending}
            >
              <Trash2 className="w-4 h-4" />
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent className="glass border-white/10 text-white">
            <AlertDialogHeader>
              <AlertDialogTitle>Delete Bot?</AlertDialogTitle>
              <AlertDialogDescription className="text-white/60">
                Are you sure you want to delete <span className="text-white font-medium">{bot.name}</span>? This action cannot be undone.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel className="bg-transparent border-white/10 hover:bg-white/10 text-white hover:text-white">Cancel</AlertDialogCancel>
              <AlertDialogAction 
                onClick={() => deleteBot.mutate(bot.id)}
                className="bg-red-500 hover:bg-red-600 text-white border-none"
              >
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>

      <div className="space-y-3 relative z-10">
        <div className="p-3 rounded-lg bg-black/40 border border-white/5 font-mono text-xs text-white/60 truncate">
          Token: ••••••••••••••••••••••••
        </div>
        
        <div className="grid grid-cols-2 gap-3">
          <div className="p-3 rounded-lg bg-white/5 border border-white/5 flex items-center gap-3">
            <Activity className="w-4 h-4 text-violet-400" />
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-wider text-white/40">Uptime</span>
              <span className="text-sm font-medium text-white">99.9%</span>
            </div>
          </div>
          <div className="p-3 rounded-lg bg-white/5 border border-white/5 flex items-center gap-3">
            <Terminal className="w-4 h-4 text-fuchsia-400" />
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-wider text-white/40">Memory</span>
              <span className="text-sm font-medium text-white">45MB</span>
            </div>
          </div>
        </div>

        <Button 
          className={`w-full ${
            isOnline 
              ? 'bg-red-500/20 hover:bg-red-500/30 text-red-400 border border-red-500/50' 
              : 'bg-green-500/20 hover:bg-green-500/30 text-green-400 border border-green-500/50'
          }`}
          onClick={() => botAction.mutate({ id: bot.id, action: isOnline ? 'stop' : 'start' })}
          disabled={isPending}
        >
          <Power className="w-4 h-4 mr-2" />
          {isPending ? 'Processing...' : (isOnline ? 'Stop Bot' : 'Start Bot')}
        </Button>
      </div>
    </motion.div>
  );
}
