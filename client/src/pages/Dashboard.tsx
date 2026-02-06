import { useQuery, useMutation } from "@tanstack/react-query";
import { Account, User } from "@shared/schema";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { apiRequest, queryClient } from "@/lib/queryClient";
import { useToast } from "@/hooks/use-toast";
import { Loader2, Plus, Play, Square, Trash2, CheckCircle2, XCircle, Users, Activity, ShieldCheck, User as UserIcon, MoreVertical, Key } from "lucide-react";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from "@/components/ui/dialog";

export default function Dashboard() {
  const { toast } = useToast();
  const [newToken, setNewToken] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { data: user } = useQuery<User>({ queryKey: ["/api/user"] });
  const { data: accounts, isLoading } = useQuery<Account[]>({ 
    queryKey: ["/api/accounts"] 
  });

  const createMutation = useMutation({
    mutationFn: async (token: string) => {
      await apiRequest("POST", "/api/accounts", { token });
    },
    onSuccess: () => {
      setNewToken("");
      setIsModalOpen(false);
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
      toast({ title: "Account linked successfully" });
    }
  });

  const validateMutation = useMutation({
    mutationFn: async () => {
      await apiRequest("POST", "/api/accounts/validate");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
      toast({ title: "Account validation completed" });
    }
  });

  const startMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiRequest("POST", `/api/accounts/${id}/start`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
    }
  });

  const stopMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiRequest("POST", `/api/accounts/${id}/stop`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiRequest("DELETE", `/api/accounts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
      toast({ title: "Account removed" });
    }
  });

  if (isLoading) return (
    <div className="flex items-center justify-center min-h-[400px]">
      <Loader2 className="h-8 w-8 animate-spin text-purple-500" />
    </div>
  );

  return (
    <div className="container max-w-7xl mx-auto p-6 md:p-8 space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h2 className="text-4xl font-black tracking-tight text-white mb-2 font-display">Dashboard</h2>
          <p className="text-slate-400 text-lg">Your discord infrastructure at a glance.</p>
        </div>
        <div className="flex items-center gap-3">
          <Button 
            onClick={() => validateMutation.mutate()} 
            disabled={validateMutation.isPending || !accounts?.length}
            variant="outline"
            className="border-white/10 hover:bg-white/5"
            data-testid="button-validate-users"
          >
            {validateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <ShieldCheck className="h-4 w-4 mr-2" />}
            Validate Status
          </Button>

          <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
            <DialogTrigger asChild>
              <Button className="bg-purple-600 hover:bg-purple-500 shadow-lg shadow-purple-500/20" data-testid="button-add-account-modal">
                <Plus className="h-4 w-4 mr-2" />
                Add Account
              </Button>
            </DialogTrigger>
            <DialogContent className="bg-[#121216] border-white/10 text-white">
              <DialogHeader>
                <DialogTitle>Link New Account</DialogTitle>
                <DialogDescription className="text-slate-400">
                  Enter your Discord user token to start hosting this account.
                </DialogDescription>
              </DialogHeader>
              <div className="py-4 space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-300">Account Token</label>
                  <div className="relative">
                    <Key className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
                    <Input 
                      type="password"
                      placeholder="MTY3M..." 
                      className="bg-[#0a0a0c] border-white/10 pl-10 focus:ring-purple-500/50"
                      value={newToken} 
                      onChange={(e) => setNewToken(e.target.value)}
                      data-testid="input-token"
                    />
                  </div>
                  <p className="text-[11px] text-slate-500">Your token is encrypted and never shared.</p>
                </div>
              </div>
              <DialogFooter>
                <Button 
                  variant="ghost" 
                  onClick={() => setIsModalOpen(false)}
                  className="text-slate-400 hover:text-white"
                >
                  Cancel
                </Button>
                <Button 
                  onClick={() => createMutation.mutate(newToken)} 
                  disabled={createMutation.isPending || !newToken}
                  className="bg-purple-600 hover:bg-purple-500"
                  data-testid="button-confirm-add"
                >
                  {createMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : "Link Account"}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {[
          { title: "Verification", value: user?.isAuthed ? "Authorized" : "Pending", icon: ShieldCheck, color: user?.isAuthed ? "text-emerald-500" : "text-rose-500", bg: user?.isAuthed ? "bg-emerald-500/10" : "bg-rose-500/10" },
          { title: "Accounts", value: accounts?.length || 0, icon: Users, color: "text-purple-400", bg: "bg-purple-500/10" },
          { title: "Active PIDs", value: accounts?.filter(a => a.status === 'online').length || 0, icon: Activity, color: "text-blue-400", bg: "bg-blue-500/10" }
        ].map((stat, i) => (
          <Card key={i} className="bg-[#121216] border-white/5 hover:border-white/10 transition-colors hover-elevate">
            <CardHeader className="flex flex-row items-center justify-between pb-2 gap-1">
              <span className="text-xs font-bold uppercase tracking-widest text-slate-500">{stat.title}</span>
              <div className={`p-2 rounded-lg ${stat.bg} ${stat.color}`}>
                <stat.icon className="h-4 w-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-black text-white font-display">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {accounts?.map((acc) => (
          <Card key={acc.id} className="bg-[#121216] border-white/5 group hover:border-purple-500/30 transition-all duration-300 hover:shadow-2xl hover:shadow-purple-500/5">
            <CardHeader className="pb-4">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                  <div className="relative">
                    <Avatar className="h-14 w-14 border-2 border-[#1a1a20]">
                      <AvatarImage src={acc.discordAvatar || undefined} />
                      <AvatarFallback className="bg-gradient-to-br from-purple-600 to-blue-600 text-white"><UserIcon /></AvatarFallback>
                    </Avatar>
                    <div className={`absolute bottom-0 right-0 h-4 w-4 rounded-full border-4 border-[#121216] ${acc.status === 'online' ? 'bg-emerald-500 shadow-[0_0_10px_rgba(16,185,129,0.5)]' : 'bg-slate-600'}`} />
                  </div>
                  <div className="flex flex-col">
                    <h3 className="text-lg font-bold text-white group-hover:text-purple-400 transition-colors truncate max-w-[140px] font-display">
                      {acc.discordUsername || "New Instance"}
                    </h3>
                    <div className="flex items-center gap-1 mt-0.5">
                      <Badge variant="outline" className="text-[10px] font-mono border-white/5 text-slate-500 py-0 px-1.5">
                        {acc.token.substring(0, 6)}...
                      </Badge>
                    </div>
                  </div>
                </div>
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="text-slate-600 hover:text-rose-500 hover:bg-rose-500/10"
                  onClick={() => deleteMutation.mutate(acc.id)}
                  data-testid={`button-delete-${acc.id}`}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-3">
                <div className="bg-[#0a0a0c] p-3 rounded-xl border border-white/5">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-slate-500 block mb-1">Guilds</span>
                  <span className="text-lg font-black text-white font-display">{acc.guildsCount}</span>
                </div>
                <div className="bg-[#0a0a0c] p-3 rounded-xl border border-white/5">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-slate-500 block mb-1">Friends</span>
                  <span className="text-lg font-black text-white font-display">{acc.friendsCount}</span>
                </div>
              </div>

              {acc.pid && (
                <div className="flex items-center justify-between px-3 py-2 rounded-lg bg-emerald-500/5 border border-emerald-500/10">
                  <span className="text-[10px] font-bold uppercase tracking-wider text-emerald-500/60">Active PID</span>
                  <span className="font-mono text-xs text-emerald-400 font-bold">{acc.pid}</span>
                </div>
              )}

              <div className="flex gap-2 pt-2">
                {acc.status === 'online' ? (
                  <Button 
                    variant="outline" 
                    className="flex-1 bg-rose-500/5 border-rose-500/20 text-rose-500 hover:bg-rose-500 hover:text-white transition-all duration-300 font-bold"
                    onClick={() => stopMutation.mutate(acc.id)}
                    disabled={stopMutation.isPending}
                    data-testid={`button-stop-${acc.id}`}
                  >
                    <Square className="h-4 w-4 mr-2 fill-current" />
                    Kill PID
                  </Button>
                ) : (
                  <Button 
                    className="flex-1 bg-purple-600 hover:bg-purple-500 shadow-lg shadow-purple-500/20 font-bold"
                    onClick={() => startMutation.mutate(acc.id)}
                    disabled={startMutation.isPending}
                    data-testid={`button-start-${acc.id}`}
                  >
                    <Play className="h-4 w-4 mr-2 fill-current" />
                    Deploy Instance
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))}

        {accounts?.length === 0 && (
          <button 
            onClick={() => setIsModalOpen(true)}
            className="flex flex-col items-center justify-center p-12 bg-[#121216]/50 border-2 border-dashed border-white/5 rounded-2xl hover:border-purple-500/30 hover:bg-purple-500/5 transition-all group"
          >
            <div className="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
              <Plus className="h-8 w-8 text-slate-600 group-hover:text-purple-500" />
            </div>
            <h3 className="text-lg font-bold text-slate-300 group-hover:text-white">Link Account</h3>
            <p className="text-sm text-slate-500">Host your first Discord instance.</p>
          </button>
        )}
      </div>
    </div>
  );
}
