import { useQuery, useMutation } from "@tanstack/react-query";
import { User, Account } from "@shared/schema";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { apiRequest, queryClient } from "@/lib/queryClient";
import { Loader2, ShieldCheck, CheckCircle2, XCircle, Plus } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

export default function Admin() {
  const { toast } = useToast();
  const [manualId, setManualId] = useState("");
  const { data: users, isLoading: usersLoading } = useQuery<User[]>({ 
    queryKey: ["/api/admin/users"] 
  });
  const { data: accounts, isLoading: accountsLoading } = useQuery<Account[]>({ 
    queryKey: ["/api/admin/accounts"] 
  });

  const hostManualMutation = useMutation({
    mutationFn: async (discordId: string) => {
      await apiRequest("POST", "/api/admin/host-manual", { discordId });
    },
    onSuccess: () => {
      setManualId("");
      queryClient.invalidateQueries({ queryKey: ["/api/admin/accounts"] });
      toast({ title: "Manual host request sent" });
    }
  });

  const authMutation = useMutation({
    mutationFn: async (id: number) => {
      await apiRequest("POST", `/api/admin/users/${id}/auth`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/admin/users"] });
      toast({ title: "User authorized" });
    }
  });

  const validateInstancesMutation = useMutation({
    mutationFn: async () => {
      await apiRequest("POST", "/api/admin/accounts/validate");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/api/admin/accounts"] });
      toast({ title: "Instances validated" });
    }
  });

  if (usersLoading || accountsLoading) return <div className="flex justify-center p-8"><Loader2 className="h-8 w-8 animate-spin" /></div>;

  return (
    <div className="container max-w-7xl mx-auto p-6 md:p-8 space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h2 className="text-4xl font-black tracking-tight text-white mb-2 font-display text-gradient">Admin Panel</h2>
          <p className="text-white/40 text-lg font-medium">System-wide management and monitoring.</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 bg-[#121216] border border-white/5 rounded-xl px-3 py-1.5">
            <Input 
              placeholder="Discord ID" 
              className="bg-transparent border-0 h-8 w-40 text-sm focus-visible:ring-0" 
              value={manualId}
              onChange={(e) => setManualId(e.target.value)}
            />
            <Button 
              size="sm" 
              onClick={() => hostManualMutation.mutate(manualId)}
              disabled={hostManualMutation.isPending || !manualId}
              className="bg-purple-600 hover:bg-purple-500 h-8"
            >
              {hostManualMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4 mr-1" />}
              Host
            </Button>
          </div>
          <Button 
            onClick={() => validateInstancesMutation.mutate()} 
            disabled={validateInstancesMutation.isPending}
            variant="outline"
            className="border-white/10 hover:bg-white/5"
            data-testid="button-validate-instances"
          >
            {validateInstancesMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <ShieldCheck className="h-4 w-4 mr-2" />}
            Validate Instances
          </Button>
        </div>
      </div>

      <Card className="bg-[#121216] border-white/5">
        <CardHeader>
          <CardTitle className="font-display">Users</CardTitle>
          <CardDescription className="text-white/40">All registered users in the system.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader className="border-white/5">
              <TableRow className="hover:bg-transparent border-white/5">
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">User</TableHead>
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">Discord ID</TableHead>
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">Authed</TableHead>
                <TableHead className="text-right text-white/40 uppercase text-[10px] font-black tracking-widest">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {users?.map((user) => (
                <TableRow key={user.id} className="border-white/5 hover:bg-white/[0.02]">
                  <TableCell className="font-bold text-white/90">{user.username}</TableCell>
                  <TableCell className="font-mono text-xs text-white/40">{user.discordId}</TableCell>
                  <TableCell>
                    {user.isAuthed ? (
                      <CheckCircle2 className="h-5 w-5 text-emerald-500" />
                    ) : (
                      <XCircle className="h-5 w-5 text-rose-500" />
                    )}
                  </TableCell>
                  <TableCell className="text-right">
                    {!user.isAuthed && (
                      <Button 
                        size="sm" 
                        variant="outline" 
                        onClick={() => authMutation.mutate(user.id)}
                        disabled={authMutation.isPending}
                        className="border-white/10 hover:bg-purple-500/10 hover:text-purple-400"
                        data-testid={`button-auth-user-${user.id}`}
                      >
                        Auth User
                      </Button>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card className="bg-[#121216] border-white/5">
        <CardHeader>
          <CardTitle className="font-display">Active Instances (Accounts)</CardTitle>
          <CardDescription className="text-white/40">Global view of all Discord accounts linked to the platform.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader className="border-white/5">
              <TableRow className="hover:bg-transparent border-white/5">
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">Account</TableHead>
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">Token (Obfuscated)</TableHead>
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">Status</TableHead>
                <TableHead className="text-white/40 uppercase text-[10px] font-black tracking-widest">PID</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {accounts?.map((acc) => (
                <TableRow key={acc.id} className="border-white/5 hover:bg-white/[0.02]">
                  <TableCell className="font-bold text-white/90">{acc.discordUsername || "Unknown"}</TableCell>
                  <TableCell className="font-mono text-xs text-white/40">
                    {acc.token.substring(0, 10)}********************
                  </TableCell>
                  <TableCell>
                    <span className={`capitalize text-xs font-bold px-2 py-0.5 rounded-full ${acc.status === 'online' ? 'bg-emerald-500/10 text-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.2)]' : 'bg-white/5 text-white/40'}`}>
                      {acc.status}
                    </span>
                  </TableCell>
                  <TableCell className="font-mono text-xs text-white/40">{acc.pid || "-"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
