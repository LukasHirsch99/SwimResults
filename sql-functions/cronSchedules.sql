select
  cron.schedule (
    'updateTodaysMeets',
    '* * * * *', -- every minute
    $$
      select updatetodaysmeets();
    $$
  );

select
  cron.schedule (
    'insertnewmeets',
  --  ┌───────────── min (0 - 59)
  --  │ ┌────────────── hour (0 - 23)
  --  │ │ ┌─────────────── day of month (1 - 31)
  --  │ │ │ ┌──────────────── month (1 - 12)
  --  │ │ │ │ ┌───────────────── day of week (0 - 6) (0 to 6 are Sunday to
  --  │ │ │ │ │                  Saturday, or use names; 7 is also Sunday)
  --  │ │ │ │ │
  --  │ │ │ │ │
  --  * * * * *
    '10 0 * * *', -- every day at 00:10
    $$
      select
      net.http_post(
          url:='https://qeudknoyuvjztxvgbmou.supabase.co/functions/v1/InsertUpcomingMeets',
          headers:='{"Content-Type": "application/json", "Authorization": "Bearer YOUR_ANON_KEY"}'::jsonb,
          body:=concat('{"time": "', now(), '"}')::jsonb,
          timeout_milliseconds := 600000
      );
    $$
  );
