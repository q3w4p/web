import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api, buildUrl, type InsertBot, type UpdateBotRequest } from "@shared/routes";
import { useToast } from "@/hooks/use-toast";

export function useBots() {
  return useQuery({
    queryKey: [api.bots.list.path],
    queryFn: async () => {
      const res = await fetch(api.bots.list.path, { credentials: "include" });
      if (!res.ok) throw new Error("Failed to fetch bots");
      return api.bots.list.responses[200].parse(await res.json());
    },
  });
}

export function useBot(id: number) {
  return useQuery({
    queryKey: [api.bots.get.path, id],
    queryFn: async () => {
      const url = buildUrl(api.bots.get.path, { id });
      const res = await fetch(url, { credentials: "include" });
      if (!res.ok) throw new Error("Failed to fetch bot");
      return api.bots.get.responses[200].parse(await res.json());
    },
    enabled: !!id,
  });
}

export function useCreateBot() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: async (data: InsertBot) => {
      const res = await fetch(api.bots.create.path, {
        method: api.bots.create.method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
        credentials: "include",
      });
      
      if (!res.ok) {
        if (res.status === 400) {
          const error = api.bots.create.responses[400].parse(await res.json());
          throw new Error(error.message);
        }
        throw new Error("Failed to create bot");
      }
      return api.bots.create.responses[201].parse(await res.json());
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [api.bots.list.path] });
      toast({
        title: "Success",
        description: "Bot created successfully",
      });
    },
    onError: (error) => {
      toast({
        title: "Error",
        description: error.message,
        variant: "destructive",
      });
    },
  });
}

export function useDeleteBot() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: async (id: number) => {
      const url = buildUrl(api.bots.delete.path, { id });
      const res = await fetch(url, { 
        method: api.bots.delete.method,
        credentials: "include" 
      });
      if (!res.ok) throw new Error("Failed to delete bot");
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [api.bots.list.path] });
      toast({
        title: "Deleted",
        description: "Bot removed successfully",
      });
    },
  });
}

export function useBotAction() {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: async ({ id, action }: { id: number, action: 'start' | 'stop' }) => {
      const endpoint = action === 'start' ? api.bots.start : api.bots.stop;
      const url = buildUrl(endpoint.path, { id });
      
      const res = await fetch(url, { 
        method: endpoint.method,
        credentials: "include"
      });

      if (!res.ok) throw new Error(`Failed to ${action} bot`);
      return endpoint.responses[200].parse(await res.json());
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: [api.bots.list.path] });
      toast({
        title: variables.action === 'start' ? "Bot Started" : "Bot Stopped",
        description: `Bot successfully ${variables.action}ed`,
      });
    },
  });
}
