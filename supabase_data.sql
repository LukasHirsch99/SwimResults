--
-- PostgreSQL database dump
--

-- Dumped from database version 15.1 (Ubuntu 15.1-1.pgdg20.04+1)
-- Dumped by pg_dump version 16.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: gender; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.gender AS ENUM (
    'M',
    'W',
    'X'
);


ALTER TYPE public.gender OWNER TO postgres;

--
-- Name: maxids_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.maxids_type AS (
	maxsessionid integer,
	maxeventid integer,
	maxheatid integer,
	maxresultid integer
);


ALTER TYPE public.maxids_type OWNER TO postgres;

--
-- Name: getswimmersbynameformeet(integer, character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.getswimmersbynameformeet(meetingid integer, swimmername character varying) RETURNS TABLE(id integer, name character varying, birthyear integer, clubid integer, gender public.gender, firstname text, lastname text, isrelay boolean, clubname character varying, nationality character varying)
    LANGUAGE plpgsql
    AS $$

begin

return query select distinct sw.*, c.name clubname, c.nationality from session s
join event e on e.sessionid = s.id
join heat h on h.eventid = e.id
join start st on st.heatid = h.id
join swimmer sw on sw.id = st.swimmerid
join club c on c.id = sw.clubid

where
  s.meetid = meetingid and (sw.firstname ilike swimmername || '%' or sw.lastname ilike swimmername || '%');
end;

$$;


ALTER FUNCTION public.getswimmersbynameformeet(meetingid integer, swimmername character varying) OWNER TO postgres;

--
-- Name: install_available_extensions_and_test(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.install_available_extensions_and_test() RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE extension_name TEXT;
allowed_extentions TEXT[] := string_to_array(current_setting('supautils.privileged_extensions'), ',');
BEGIN 
  FOREACH extension_name IN ARRAY allowed_extentions 
  LOOP
    SELECT trim(extension_name) INTO extension_name;
    /* skip below extensions check for now */
    CONTINUE WHEN extension_name = 'pgroonga' OR  extension_name = 'pgroonga_database' OR extension_name = 'pgsodium';
    CONTINUE WHEN extension_name = 'plpgsql' OR  extension_name = 'plpgsql_check' OR extension_name = 'pgtap';
    CONTINUE WHEN extension_name = 'supabase_vault' OR extension_name = 'wrappers';
    RAISE notice 'START TEST FOR: %', extension_name;
    EXECUTE format('DROP EXTENSION IF EXISTS %s CASCADE', quote_ident(extension_name));
    EXECUTE format('CREATE EXTENSION %s CASCADE', quote_ident(extension_name));
    RAISE notice 'END TEST FOR: %', extension_name;
  END LOOP;
    RAISE notice 'EXTENSION TESTS COMPLETED..';
    return true;
END;
$$;


ALTER FUNCTION public.install_available_extensions_and_test() OWNER TO postgres;

--
-- Name: maxheatid(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.maxheatid() RETURNS integer
    LANGUAGE sql
    AS $$
  SELECT COALESCE(MAX(id), 0) FROM heat;
$$;


ALTER FUNCTION public.maxheatid() OWNER TO postgres;

--
-- Name: maxid(character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.maxid(tablename character varying) RETURNS integer
    LANGUAGE plpgsql
    AS $$begin
if lower(tablename) = 'session' then
  return (select max(id) from session);
elsif lower(tablename) = 'event' then
  return (select max(id) from event);
elsif lower(tablename) = 'heat' then
  return (select max(id) from heat);
elsif lower(tablename) = 'result' then
  return (select max(id) from result);
else
  return -1;
end if;
end$$;


ALTER FUNCTION public.maxid(tablename character varying) OWNER TO postgres;

--
-- Name: maxids(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.maxids() RETURNS public.maxids_type
    LANGUAGE plpgsql
    AS $$
DECLARE 
  ret maxids_type;
BEGIN
  SELECT COALESCE(MAX(id), 0) into ret.maxsessionid FROM session;
  SELECT COALESCE(MAX(id), 0) into ret.maxeventid FROM event;
  SELECT COALESCE(MAX(id), 0) into ret.maxheatid FROM heat;
  SELECT COALESCE(MAX(id), 0) into ret.maxresultid FROM result;

  return ret;
END
$$;


ALTER FUNCTION public.maxids() OWNER TO postgres;

--
-- Name: maxresultid(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.maxresultid() RETURNS integer
    LANGUAGE sql
    AS $$
  SELECT COALESCE(MAX(id), 0) FROM result;
$$;


ALTER FUNCTION public.maxresultid() OWNER TO postgres;

--
-- Name: nextheatid(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.nextheatid() RETURNS integer
    LANGUAGE sql
    AS $$
  Select nextval(pg_get_serial_sequence('heat', 'id'));
$$;


ALTER FUNCTION public.nextheatid() OWNER TO postgres;

--
-- Name: nextresultid(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.nextresultid() RETURNS integer
    LANGUAGE sql
    AS $$
  Select nextval(pg_get_serial_sequence('result', 'id'));
$$;


ALTER FUNCTION public.nextresultid() OWNER TO postgres;

--
-- Name: resultcountforevent(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.resultcountforevent(p_eventid integer) RETURNS integer
    LANGUAGE sql
    AS $$
    select count(*) from result join ageclass a on result.ageclassid = a.id and a.eventid = p_eventid;
$$;


ALTER FUNCTION public.resultcountforevent(p_eventid integer) OWNER TO postgres;

--
-- Name: startcountforevent(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.startcountforevent(p_eventid integer) RETURNS integer
    LANGUAGE sql
    AS $$
    select count(*) from start join heat h on start.heatid = h.id and h.eventid = p_eventid;
$$;


ALTER FUNCTION public.startcountforevent(p_eventid integer) OWNER TO postgres;

--
-- Name: updatetodaysmeets(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.updatetodaysmeets() RETURNS integer
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.updatetodaysmeets() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ageclass; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ageclass (
    resultid integer NOT NULL,
    name character varying NOT NULL,
    "position" integer,
    timetofirst text
);


ALTER TABLE public.ageclass OWNER TO postgres;

--
-- Name: club; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.club (
    id integer NOT NULL,
    name character varying NOT NULL,
    nationality character varying
);


ALTER TABLE public.club OWNER TO postgres;

--
-- Name: event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event (
    id integer NOT NULL,
    sessionid integer NOT NULL,
    displaynr integer NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.event OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.event_id_seq OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;


--
-- Name: event_id_seq1; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.event ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.event_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: heat; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.heat (
    id integer NOT NULL,
    eventid integer NOT NULL,
    heatnr integer NOT NULL
);


ALTER TABLE public.heat OWNER TO postgres;

--
-- Name: meet; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.meet (
    id integer NOT NULL,
    name character varying NOT NULL,
    image character varying,
    invitations character varying[],
    deadline timestamp without time zone NOT NULL,
    address character varying NOT NULL,
    startdate date NOT NULL,
    enddate date NOT NULL,
    googlemapslink character varying,
    msecmid integer
);


ALTER TABLE public.meet OWNER TO postgres;

--
-- Name: result; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.result (
    eventid integer NOT NULL,
    swimmerid integer NOT NULL,
    "time" time(2) without time zone,
    splits character varying,
    finapoints integer,
    additionalinfo character varying,
    id integer NOT NULL,
    penalty boolean,
    reactiontime double precision
);


ALTER TABLE public.result OWNER TO postgres;

--
-- Name: session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.session (
    id integer NOT NULL,
    meetid integer NOT NULL,
    day date NOT NULL,
    warmupstart time without time zone,
    sessionstart time without time zone,
    displaynr integer NOT NULL
);


ALTER TABLE public.session OWNER TO postgres;

--
-- Name: session_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.session_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.session_id_seq OWNER TO postgres;

--
-- Name: session_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.session_id_seq OWNED BY public.session.id;


--
-- Name: start; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.start (
    heatid integer NOT NULL,
    swimmerid integer NOT NULL,
    lane integer NOT NULL,
    "time" time(2) without time zone
);


ALTER TABLE public.start OWNER TO postgres;

--
-- Name: swimmer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.swimmer (
    id integer NOT NULL,
    birthyear integer,
    clubid integer NOT NULL,
    gender public.gender NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    isrelay boolean DEFAULT false NOT NULL
);


ALTER TABLE public.swimmer OWNER TO postgres;

--
-- Name: COLUMN swimmer.isrelay; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.swimmer.isrelay IS 'When the swimmer is actually a relay, this value is true';


--
-- Name: session id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session ALTER COLUMN id SET DEFAULT nextval('public.session_id_seq'::regclass);


--
-- Name: ageclass ageclass_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ageclass
    ADD CONSTRAINT ageclass_pkey PRIMARY KEY (resultid, name);


--
-- Name: club club_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.club
    ADD CONSTRAINT club_pkey PRIMARY KEY (id);


--
-- Name: event event_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_id_key UNIQUE (id);


--
-- Name: event event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_pkey PRIMARY KEY (sessionid, displaynr, name);


--
-- Name: heat heat_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.heat
    ADD CONSTRAINT heat_id_key UNIQUE (id);


--
-- Name: heat heat_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.heat
    ADD CONSTRAINT heat_pkey PRIMARY KEY (eventid, heatnr);


--
-- Name: meet meet_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.meet
    ADD CONSTRAINT meet_pkey PRIMARY KEY (id);


--
-- Name: result result_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.result
    ADD CONSTRAINT result_id_key UNIQUE (id);


--
-- Name: result result_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.result
    ADD CONSTRAINT result_pkey PRIMARY KEY (eventid, swimmerid);


--
-- Name: session session_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_id_key UNIQUE (id);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (meetid, displaynr);


--
-- Name: start start_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.start
    ADD CONSTRAINT start_pkey PRIMARY KEY (heatid, lane);


--
-- Name: swimmer swimmer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.swimmer
    ADD CONSTRAINT swimmer_pkey PRIMARY KEY (id);


--
-- Name: ageclass ageclass_resultid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ageclass
    ADD CONSTRAINT ageclass_resultid_fkey FOREIGN KEY (resultid) REFERENCES public.result(id) ON DELETE CASCADE;


--
-- Name: event event_sessionid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_sessionid_fkey FOREIGN KEY (sessionid) REFERENCES public.session(id) ON DELETE CASCADE;


--
-- Name: heat heat_eventid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.heat
    ADD CONSTRAINT heat_eventid_fkey FOREIGN KEY (eventid) REFERENCES public.event(id) ON DELETE CASCADE;


--
-- Name: result result_eventid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.result
    ADD CONSTRAINT result_eventid_fkey FOREIGN KEY (eventid) REFERENCES public.event(id) ON DELETE CASCADE;


--
-- Name: result result_swimmerid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.result
    ADD CONSTRAINT result_swimmerid_fkey FOREIGN KEY (swimmerid) REFERENCES public.swimmer(id);


--
-- Name: session session_meetid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_meetid_fkey FOREIGN KEY (meetid) REFERENCES public.meet(id) ON DELETE CASCADE;


--
-- Name: start start_heatid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.start
    ADD CONSTRAINT start_heatid_fkey FOREIGN KEY (heatid) REFERENCES public.heat(id) ON DELETE CASCADE;


--
-- Name: start start_swimmerid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.start
    ADD CONSTRAINT start_swimmerid_fkey FOREIGN KEY (swimmerid) REFERENCES public.swimmer(id);


--
-- Name: swimmer swimmer_clubid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.swimmer
    ADD CONSTRAINT swimmer_clubid_fkey FOREIGN KEY (clubid) REFERENCES public.club(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;
GRANT USAGE ON SCHEMA public TO anon;
GRANT USAGE ON SCHEMA public TO authenticated;
GRANT USAGE ON SCHEMA public TO service_role;


--
-- Name: FUNCTION getswimmersbynameformeet(meetingid integer, swimmername character varying); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.getswimmersbynameformeet(meetingid integer, swimmername character varying) TO anon;
GRANT ALL ON FUNCTION public.getswimmersbynameformeet(meetingid integer, swimmername character varying) TO authenticated;
GRANT ALL ON FUNCTION public.getswimmersbynameformeet(meetingid integer, swimmername character varying) TO service_role;


--
-- Name: FUNCTION install_available_extensions_and_test(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.install_available_extensions_and_test() TO anon;
GRANT ALL ON FUNCTION public.install_available_extensions_and_test() TO authenticated;
GRANT ALL ON FUNCTION public.install_available_extensions_and_test() TO service_role;


--
-- Name: FUNCTION maxheatid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.maxheatid() TO anon;
GRANT ALL ON FUNCTION public.maxheatid() TO authenticated;
GRANT ALL ON FUNCTION public.maxheatid() TO service_role;


--
-- Name: FUNCTION maxid(tablename character varying); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.maxid(tablename character varying) TO anon;
GRANT ALL ON FUNCTION public.maxid(tablename character varying) TO authenticated;
GRANT ALL ON FUNCTION public.maxid(tablename character varying) TO service_role;


--
-- Name: FUNCTION maxids(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.maxids() TO anon;
GRANT ALL ON FUNCTION public.maxids() TO authenticated;
GRANT ALL ON FUNCTION public.maxids() TO service_role;


--
-- Name: FUNCTION maxresultid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.maxresultid() TO anon;
GRANT ALL ON FUNCTION public.maxresultid() TO authenticated;
GRANT ALL ON FUNCTION public.maxresultid() TO service_role;


--
-- Name: FUNCTION nextheatid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.nextheatid() TO anon;
GRANT ALL ON FUNCTION public.nextheatid() TO authenticated;
GRANT ALL ON FUNCTION public.nextheatid() TO service_role;


--
-- Name: FUNCTION nextresultid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.nextresultid() TO anon;
GRANT ALL ON FUNCTION public.nextresultid() TO authenticated;
GRANT ALL ON FUNCTION public.nextresultid() TO service_role;


--
-- Name: FUNCTION resultcountforevent(p_eventid integer); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.resultcountforevent(p_eventid integer) TO anon;
GRANT ALL ON FUNCTION public.resultcountforevent(p_eventid integer) TO authenticated;
GRANT ALL ON FUNCTION public.resultcountforevent(p_eventid integer) TO service_role;


--
-- Name: FUNCTION startcountforevent(p_eventid integer); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.startcountforevent(p_eventid integer) TO anon;
GRANT ALL ON FUNCTION public.startcountforevent(p_eventid integer) TO authenticated;
GRANT ALL ON FUNCTION public.startcountforevent(p_eventid integer) TO service_role;


--
-- Name: FUNCTION updatetodaysmeets(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.updatetodaysmeets() TO anon;
GRANT ALL ON FUNCTION public.updatetodaysmeets() TO authenticated;
GRANT ALL ON FUNCTION public.updatetodaysmeets() TO service_role;


--
-- Name: TABLE ageclass; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.ageclass TO anon;
GRANT ALL ON TABLE public.ageclass TO authenticated;
GRANT ALL ON TABLE public.ageclass TO service_role;


--
-- Name: TABLE club; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.club TO anon;
GRANT ALL ON TABLE public.club TO authenticated;
GRANT ALL ON TABLE public.club TO service_role;


--
-- Name: TABLE event; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.event TO anon;
GRANT ALL ON TABLE public.event TO authenticated;
GRANT ALL ON TABLE public.event TO service_role;


--
-- Name: SEQUENCE event_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.event_id_seq TO anon;
GRANT ALL ON SEQUENCE public.event_id_seq TO authenticated;
GRANT ALL ON SEQUENCE public.event_id_seq TO service_role;


--
-- Name: SEQUENCE event_id_seq1; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.event_id_seq1 TO anon;
GRANT ALL ON SEQUENCE public.event_id_seq1 TO authenticated;
GRANT ALL ON SEQUENCE public.event_id_seq1 TO service_role;


--
-- Name: TABLE heat; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.heat TO anon;
GRANT ALL ON TABLE public.heat TO authenticated;
GRANT ALL ON TABLE public.heat TO service_role;


--
-- Name: TABLE meet; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.meet TO anon;
GRANT ALL ON TABLE public.meet TO authenticated;
GRANT ALL ON TABLE public.meet TO service_role;


--
-- Name: TABLE result; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.result TO anon;
GRANT ALL ON TABLE public.result TO authenticated;
GRANT ALL ON TABLE public.result TO service_role;


--
-- Name: TABLE session; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.session TO anon;
GRANT ALL ON TABLE public.session TO authenticated;
GRANT ALL ON TABLE public.session TO service_role;


--
-- Name: SEQUENCE session_id_seq; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.session_id_seq TO anon;
GRANT ALL ON SEQUENCE public.session_id_seq TO authenticated;
GRANT ALL ON SEQUENCE public.session_id_seq TO service_role;


--
-- Name: TABLE start; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.start TO anon;
GRANT ALL ON TABLE public.start TO authenticated;
GRANT ALL ON TABLE public.start TO service_role;


--
-- Name: TABLE swimmer; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.swimmer TO anon;
GRANT ALL ON TABLE public.swimmer TO authenticated;
GRANT ALL ON TABLE public.swimmer TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON SEQUENCES TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON FUNCTIONS TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO service_role;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: supabase_admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES TO anon;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES TO authenticated;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_admin IN SCHEMA public GRANT ALL ON TABLES TO service_role;


--
-- PostgreSQL database dump complete
--

