create or replace function getSwimmersByNameForMeet (meetingid integer, swimmername character varying) returns table (
  id integer,
  name character varying,
  birthyear integer,
  clubid integer,
  gender public.gender,
  clubname character varying,
  nationality character varying
) as $$

begin

return query select distinct sw.*, c.name clubname, c.nationality from session s
join event e on e.sessionid = s.id
join heat h on h.eventid = e.id
join start st on st.heatid = h.id
join swimmer sw on sw.id = st.swimmerid
join club c on c.id = sw.clubid

where s.meetid = meetingid and sw.name ilike swimmername;
end;

$$ language plpgsql;
