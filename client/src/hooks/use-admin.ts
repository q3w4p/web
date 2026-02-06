import { useQuery } from "@tanstack/react-query";
import { api } from "@shared/routes";

export function useAdminUsers() {
  return useQuery({
    queryKey: [api.admin.users.path],
    queryFn: async () => {
      const res = await fetch(api.admin.users.path, { credentials: "include" });
      if (!res.ok) throw new Error("Failed to fetch users");
      return api.admin.users.responses[200].parse(await res.json());
    },
  });
}

export function useAdminBots() {
  return useQuery({
    queryKey: [api.admin.bots.path],
    queryFn: async () => {
      const res = await fetch(api.admin.bots.path, { credentials: "include" });
      if (!res.ok) throw new Error("Failed to fetch all bots");
      return api.admin.bots.responses[200].parse(await res.json());
    },
  });
}
