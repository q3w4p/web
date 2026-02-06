import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { HelpCircle, Key, Globe, Eye } from "lucide-react";

export default function GetToken() {
  const steps = [
    {
      title: "Open Discord in Browser",
      description: "Go to discord.com/app and log in to your account. This must be done on a desktop browser.",
      icon: Globe
    },
    {
      title: "Open Developer Tools",
      description: "Press F12 or Ctrl+Shift+I (Cmd+Option+I on Mac) to open the browser developer tools.",
      icon: Eye
    },
    {
      title: "Navigate to Network Tab",
      description: "Click on the 'Network' tab at the top of the developer tools window.",
      icon: Activity
    },
    {
      title: "Filter by XHR/Fetch",
      description: "Refresh the page (F5) and type '/api/v9/users/@me' in the filter box.",
      icon: HelpCircle
    },
    {
      title: "Find the Authorization Header",
      description: "Click on the request named '@me', go to 'Headers' tab, and scroll down to find 'authorization'. The long string next to it is your token.",
      icon: Key
    }
  ];

  return (
    <div className="container max-w-4xl mx-auto p-6 md:p-8 space-y-8 animate-in fade-in duration-500">
      <div>
        <h2 className="text-4xl font-black tracking-tight text-white mb-2 font-display text-gradient">How to Get Your Token</h2>
        <p className="text-white/40 text-lg font-medium">Follow this step-by-step guide to retrieve your Discord user token.</p>
      </div>

      <div className="grid gap-6">
        {steps.map((step, index) => (
          <Card key={index} className="relative overflow-hidden bg-[#121216] border-white/5 hover:border-white/10 transition-colors group">
            <div className="absolute left-0 top-0 bottom-0 w-1 bg-purple-600/50 group-hover:bg-purple-500 transition-colors" />
            <CardHeader className="flex flex-row items-center gap-4">
              <div className="flex items-center justify-center w-10 h-10 rounded-xl bg-purple-600/10 text-purple-400 font-black text-lg border border-purple-500/20">
                {index + 1}
              </div>
              <div className="flex-1">
                <CardTitle className="text-xl font-display text-white/90">{step.title}</CardTitle>
              </div>
              <step.icon className="h-6 w-6 text-white/20 group-hover:text-purple-400 transition-colors" />
            </CardHeader>
            <CardContent>
              <p className="text-white/40 leading-relaxed font-medium">{step.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="bg-rose-500/5 border-rose-500/20 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="text-rose-400 flex items-center gap-3 font-display">
            <ShieldAlert className="h-6 w-6" />
            Security Warning
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-rose-400/70 font-medium leading-relaxed">
            Never share your Discord token with anyone you do not trust. Anyone with your token has full access to your account. 
            Our system uses your token only to host your instances securely and encrypts it at rest.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

import { ShieldAlert, Activity } from "lucide-react";
