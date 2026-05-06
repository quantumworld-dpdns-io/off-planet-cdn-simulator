import { serve } from "https://deno.land/std@0.177.0/http/server.ts";

serve(async (req: Request) => {
  const body = await req.json();
  console.log("Node alert webhook received:", JSON.stringify(body));
  return new Response(JSON.stringify({ received: true }), {
    headers: { "Content-Type": "application/json" },
  });
});
