import { Link, useLocation } from "wouter";
import { useUser, useLogout } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { SiDiscord } from "react-icons/si";
import { motion } from "framer-motion";

export function Navbar() {
  const { data: user } = useUser();
  const { mutate: logout } = useLogout();
  const [location] = useLocation();

  const navItems = [
    { label: "Home", path: "/" },
    { label: "Dashboard", path: "/dashboard" },
    { label: "Admin", path: "/admin" },
  ];

  return (
    <motion.nav 
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="fixed top-0 left-0 right-0 z-50 flex justify-center p-6 pointer-events-none"
    >
      <div className="glass h-14 px-6 flex items-center justify-between gap-8 rounded-2xl w-full max-w-5xl pointer-events-auto shadow-2xl shadow-black/50">
        <Link href="/">
          <div className="flex items-center gap-3 cursor-pointer group">
            <div className="w-8 h-8 rounded-lg overflow-hidden border border-white/10 group-hover:scale-110 transition-transform duration-500">
              <img 
                src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg?ex=69867290&is=69852110&hm=27d1f6242d0aa1a2d312aab51b6f55eb30fc2c3d9a9894a572e3e4b80bd00379"
                alt="Logo"
                className="w-full h-full object-cover"
              />
            </div>
            <span className="font-display font-bold text-xl tracking-tighter text-gradient">Hurry</span>
          </div>
        </Link>

        <div className="hidden md:flex items-center gap-1">
          {navItems.map((item) => (
            <Link key={item.path} href={item.path}>
              <Button 
                variant="ghost" 
                className={`h-9 px-4 rounded-xl text-sm font-medium transition-all ${
                  location === item.path ? "bg-white/10 text-white" : "text-white/40 hover:text-white"
                }`}
              >
                {item.label}
              </Button>
            </Link>
          ))}
        </div>

        <div className="flex items-center gap-3">
          {user ? (
            <div className="flex items-center gap-4">
              <div className="hidden sm:flex flex-col items-end">
                <span className="text-xs font-bold text-white leading-none">{user.username}</span>
                <span className="text-[10px] text-white/30 uppercase tracking-widest font-black mt-1">Authorized</span>
              </div>
              <div className="w-9 h-9 rounded-xl border border-white/10 overflow-hidden bg-white/5">
                {user.avatar ? (
                  <img src={`https://cdn.discordapp.com/avatars/${user.discordId}/${user.avatar}.png`} alt={user.username} className="w-full h-full object-cover" />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-xs font-bold text-white/20">
                    {user.username?.[0]?.toUpperCase()}
                  </div>
                )}
              </div>
              <Button 
                variant="ghost" 
                size="sm"
                onClick={() => logout()}
                className="text-white/40 hover:text-red-400 hover:bg-red-500/10 rounded-full text-xs h-8 px-3 ml-2"
              >
                Logout
              </Button>
            </div>
          ) : (
            <a href="/api/auth/discord">
              <Button className="rival-button rival-button-discord h-10 px-6 rounded-xl text-sm font-bold flex items-center gap-2">
                <SiDiscord className="w-4 h-4" />
                Login
              </Button>
            </a>
          )}
        </div>
      </div>
    </motion.nav>
  );
}
