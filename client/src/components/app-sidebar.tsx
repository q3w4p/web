import { Sidebar, SidebarContent, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar";
import { LayoutDashboard, Shield, HelpCircle, LogOut, Terminal, Cpu } from "lucide-react";
import { Link, useLocation } from "wouter";
import { useQuery } from "@tanstack/react-query";
import { User } from "@shared/schema";

const items = [
  {
    title: "Overview",
    url: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    title: "Guide",
    url: "/get-token",
    icon: HelpCircle,
  },
];

const adminItems = [
  {
    title: "Admin Panel",
    url: "/admin",
    icon: Shield,
  },
];

export function AppSidebar() {
  const [location] = useLocation();
  const { data: user } = useQuery<User>({ queryKey: ["/api/user"] });

  return (
    <Sidebar className="border-r border-white/5 bg-[#0a0a0c]">
      <SidebarContent className="p-4">
        <div className="flex items-center gap-3 px-3 py-6 mb-4">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-600 shadow-lg shadow-purple-500/20">
            <Cpu className="h-6 w-6 text-white" />
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-black text-white uppercase tracking-tighter leading-none">BotHost</span>
            <span className="text-[10px] font-bold text-purple-400 uppercase tracking-widest mt-1">v2.0 Beta</span>
          </div>
        </div>

        <SidebarGroup>
          <SidebarGroupLabel className="text-[10px] uppercase tracking-[0.2em] font-black text-white/30 px-3 mb-2">Main Menu</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton 
                    asChild 
                    isActive={location === item.url}
                    className="h-11 rounded-xl data-[active=true]:bg-white/5 data-[active=true]:text-purple-400 hover:bg-white/5 transition-all group"
                  >
                    <Link href={item.url}>
                      <item.icon className={`h-4 w-4 ${location === item.url ? 'text-purple-400' : 'text-white/40 group-hover:text-white/90'}`} />
                      <span className="font-semibold">{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {user?.isAdmin && (
          <SidebarGroup className="mt-4">
            <SidebarGroupLabel className="text-[10px] uppercase tracking-[0.2em] font-black text-white/30 px-3 mb-2">Internal</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {adminItems.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton 
                      asChild 
                      isActive={location === item.url}
                      className="h-11 rounded-xl data-[active=true]:bg-white/5 data-[active=true]:text-purple-400 hover:bg-white/5 transition-all group"
                    >
                      <Link href={item.url}>
                        <item.icon className={`h-4 w-4 ${location === item.url ? 'text-purple-400' : 'text-white/40 group-hover:text-white/90'}`} />
                        <span className="font-semibold">{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}
        
        <div className="mt-auto pt-6">
          <div className="p-4 rounded-2xl bg-gradient-to-br from-purple-600/10 to-blue-600/10 border border-white/5">
            <Terminal className="h-5 w-5 text-purple-500 mb-2" />
            <h4 className="text-xs font-bold text-white mb-1">Status: Stable</h4>
            <p className="text-[10px] text-slate-500">Latency: 24ms</p>
          </div>
        </div>
      </SidebarContent>
    </Sidebar>
  );
}
