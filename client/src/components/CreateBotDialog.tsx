import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { insertBotSchema } from "@shared/schema";
import { useCreateBot } from "@/hooks/use-bots";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Plus, Bot } from "lucide-react";

export function CreateBotDialog() {
  const [open, setOpen] = useState(false);
  const createBot = useCreateBot();

  const form = useForm<z.infer<typeof insertBotSchema>>({
    resolver: zodResolver(insertBotSchema),
    defaultValues: {
      name: "",
      token: "",
    },
  });

  function onSubmit(values: z.infer<typeof insertBotSchema>) {
    createBot.mutate(values, {
      onSuccess: () => {
        setOpen(false);
        form.reset();
      },
    });
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="bg-gradient-to-r from-violet-600 to-fuchsia-600 hover:from-violet-500 hover:to-fuchsia-500 text-white shadow-lg shadow-violet-500/25 border-none">
          <Plus className="w-4 h-4 mr-2" />
          Host New Bot
        </Button>
      </DialogTrigger>
      <DialogContent className="glass border-white/10 text-white sm:max-w-[425px]">
        <DialogHeader>
          <div className="mx-auto bg-white/10 p-3 rounded-xl mb-4">
            <Bot className="w-8 h-8 text-violet-400" />
          </div>
          <DialogTitle className="text-center text-xl">Host a new Bot</DialogTitle>
          <DialogDescription className="text-center text-white/60">
            Enter your bot's details below to start hosting.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 mt-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-white/80">Bot Name</FormLabel>
                  <FormControl>
                    <Input 
                      placeholder="My Awesome Bot" 
                      className="bg-black/20 border-white/10 text-white placeholder:text-white/20 focus:border-violet-500"
                      {...field} 
                    />
                  </FormControl>
                  <FormMessage className="text-red-400" />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="token"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-white/80">Bot Token</FormLabel>
                  <FormControl>
                    <Input 
                      type="password"
                      placeholder="MTA..." 
                      className="bg-black/20 border-white/10 text-white placeholder:text-white/20 focus:border-violet-500 font-mono text-sm"
                      {...field} 
                    />
                  </FormControl>
                  <FormMessage className="text-red-400" />
                </FormItem>
              )}
            />
            <Button 
              type="submit" 
              className="w-full bg-violet-600 hover:bg-violet-500 text-white"
              disabled={createBot.isPending}
            >
              {createBot.isPending ? "Creating..." : "Create Bot"}
            </Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
