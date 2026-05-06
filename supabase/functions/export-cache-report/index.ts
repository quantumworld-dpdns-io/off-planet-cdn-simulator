// export-cache-report — Supabase Edge Function
//
// TODO: Full implementation will:
//   1. Accept an org_id + optional date range from the request body.
//   2. Connect to the read-replica via Supabase client (service-role key).
//   3. Run an analytical query joining cache_objects, cache_object_versions,
//      content_requests, and telemetry_events for the requested period.
//   4. Export the result set using DuckDB-Wasm running inside Deno, serialised
//      to Arrow IPC format (.arrows) for efficient streaming to the caller.
//   5. Stream the Arrow IPC bytes back as application/vnd.apache.arrow.stream.
//
// This stub returns a placeholder until the DuckDB export pipeline is wired up.

import { serve } from "https://deno.land/std@0.177.0/http/server.ts";

serve(async (_req: Request) => {
  return new Response(
    JSON.stringify({ status: "not_implemented" }),
    {
      status: 501,
      headers: { "Content-Type": "application/json" },
    },
  );
});
