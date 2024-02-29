create
or replace function updateTodaysMeets () returns int as $$
declare m record;
begin
  for m in
    (select * 
    from meet
    where startdate <= current_date and enddate >= current_date)
  loop
    perform
      net.http_post (
        'https://qeudknoyuvjztxvgbmou.supabase.co/functions/v1/UpdateSchedule',
        format('{"meetId": %I}', m.id)::JSONB,
        headers := '{
        "Content-Type": "application/json",
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU"
        }'::JSONB,
        timeout_milliseconds := 60000
      );

  end loop;
  return 1;
end;
$$ language plpgsql;
