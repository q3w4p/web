import { useQuery, useMutation } from "@tanstack/react-query";
import { Account, User } from "@shared/schema";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { apiRequest, queryClient } from "@/lib/queryClient";
import { useToast } from "@/hooks/use-toast";
import { Loader2, Plus, Play, Square, Trash2, CheckCircle2, XCircle, Users, Activity, ShieldCheck, User as UserIcon, MoreVertical, Key, Search } from "lucide-react";
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
      const res = await apiRequest("POST", "/api/accounts", { token });
      return res.json();
    },
    onSuccess: () => {
      setNewToken("");
      setIsModalOpen(false);
      queryClient.invalidateQueries({ queryKey: ["/api/accounts"] });
      // Trigger validation immediately after adding to get username
      validateMutation.mutate();
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
    <div className="container max-w-7xl mx-auto p-4 md:p-6 space-y-6 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-white mb-1 font-display">Hosted Accounts</h2>
          <p className="text-slate-400 text-sm">Manage your hosted Discord accounts</p>
        </div>
        <div className="flex items-center gap-2">
          <Button 
            onClick={() => validateMutation.mutate()} 
            disabled={validateMutation.isPending || !accounts?.length}
            variant="outline"
            className="border-amber-500/20 text-amber-500 hover:bg-amber-500/10 h-9"
            data-testid="button-validate-users"
          >
            {validateMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <CheckCircle2 className="h-4 w-4 mr-2" />}
            Validate All
          </Button>

          <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
            <DialogTrigger asChild>
              <Button className="bg-purple-600/20 text-purple-400 border border-purple-500/30 hover:bg-purple-600/30 h-9" data-testid="button-add-account-modal">
                <Plus className="h-4 w-4 mr-2" />
                Host New
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

      <div className="grid gap-4 md:grid-cols-3">
        {[
          { title: "Accounts Hosted", value: accounts?.length || "-", icon: Users, color: "text-purple-400", bg: "bg-purple-500/5" },
          { title: "Hosting Limit", value: 5, icon: ShieldCheck, color: "text-purple-400", bg: "bg-purple-500/5" },
          { title: "Slots Available", value: Math.max(0, 5 - (accounts?.length || 0)) || "-", icon: Activity, color: "text-purple-400", bg: "bg-purple-500/5" }
        ].map((stat, i) => (
          <Card key={i} className="bg-[#0c0c0e] border-white/5 h-28 flex flex-col items-center justify-center">
            <div className="text-2xl font-bold text-purple-400 mb-1">{stat.value}</div>
            <div className="text-[10px] font-bold uppercase tracking-widest text-slate-600">{stat.title}</div>
          </Card>
        ))}
      </div>

      <Card className="bg-[#0c0c0e] border-white/5 overflow-hidden">
        <div className="p-4 border-b border-white/5 flex items-center justify-between">
          <h3 className="text-sm font-bold text-white">Your Accounts</h3>
        </div>
        <div className="p-4 space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-600" />
            <Input 
              placeholder="Search by username or ID..." 
              className="bg-[#08080a] border-white/5 pl-10 h-10 text-sm"
            />
          </div>
          
          <div className="grid gap-3">
            {accounts?.map((acc) => (
              <div key={acc.id} className="flex items-center justify-between p-3 bg-[#08080a] border border-white/5 rounded-lg group">
                <div className="flex items-center gap-3">
                  <Avatar className="h-10 w-10 border border-white/5">
                    <AvatarImage src={acc.discordAvatar || undefined} />
                    <AvatarFallback className="bg-purple-900/20 text-purple-400 text-xs"><UserIcon className="h-4 w-4" /></AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="text-sm font-bold text-white group-hover:text-purple-400 transition-colors">
                      {acc.discordUsername || "New Instance"}
                    </div>
                    <div className="flex items-center gap-2 mt-0.5">
                      <div className={`h-1.5 w-1.5 rounded-full ${acc.status === 'online' ? 'bg-emerald-500' : 'bg-slate-600'}`} />
                      <span className="text-[10px] text-slate-500 font-mono">{acc.token.substring(0, 8)}...</span>
                      <span className="text-[10px] text-slate-500">•</span>
                      <span className="text-[10px] text-slate-500">{acc.guildsCount} Guilds</span>
                      <span className="text-[10px] text-slate-500">•</span>
                      <span className="text-[10px] text-slate-500">{acc.friendsCount} Friends</span>
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center gap-2">
                  {acc.status === 'online' ? (
                    <Button 
                      variant="ghost" 
                      size="sm"
                      className="h-8 px-3 text-rose-500 hover:bg-rose-500/10 text-xs font-bold"
                      onClick={() => stopMutation.mutate(acc.id)}
                      disabled={stopMutation.isPending}
                    >
                      Disconnect
                    </Button>
                  ) : (
                    <Button 
                      variant="ghost" 
                      size="sm"
                      className="h-8 px-3 text-purple-400 hover:bg-purple-500/10 text-xs font-bold"
                      onClick={() => startMutation.mutate(acc.id)}
                      disabled={startMutation.isPending}
                    >
                      Deploy
                    </Button>
                  )}
                  <Button 
                    variant="ghost" 
                    size="icon" 
                    className="h-8 w-8 text-slate-600 hover:text-rose-500"
                    onClick={() => deleteMutation.mutate(acc.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}

            {accounts?.length === 0 && !isLoading && (
              <div className="flex flex-col items-center justify-center py-12 text-slate-600">
                <div className="bg-white/5 p-4 rounded-full mb-4">
                  <Activity className="h-8 w-8" />
                </div>
                <p className="text-sm">No accounts hosted yet.</p>
              </div>
            )}
          </div>
        </div>
      </Card>
    </div>
  );
}
