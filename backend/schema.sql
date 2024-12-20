create type gender as enum ('M', 'W', 'X');

create table ageclass
(
    id   serial primary key,
    name varchar not null
);

create table club
(
    id          serial primary key,
    name        varchar not null,
    nationality varchar
);

create table meet
(
    id             serial primary key,
    name           varchar   not null,
    image          varchar,
    invitations    character varying[],
    deadline       timestamp not null,
    address        varchar   not null,
    startdate      date      not null,
    enddate        date      not null,
    googlemapslink varchar,
    msecmid        integer
);

create table session
(
    id           serial unique not null,
    meetid       integer not null
        references meet
            on delete cascade,
    day          date    not null,
    warmupstart  time,
    sessionstart time,
    displaynr    integer not null,
    primary key (meetid, displaynr)
);

create table event
(
    id        serial unique not null,
    sessionid integer not null
        references session (id)
            on delete cascade,
    displaynr integer not null,
    name      varchar not null,
    primary key (sessionid, displaynr, name)
);

create table heat
(
    id      serial unique not null,
    eventid integer not null
        references event (id)
            on delete cascade,
    heatnr  integer not null,
    primary key (eventid, heatnr)
);

create table swimmer
(
    id        serial primary key,
    birthyear integer,
    clubid    integer               not null
        references club,
    gender    gender                not null,
    firstname text                  not null,
    lastname  text                  not null,
    isrelay   boolean default false not null
);

create table result
(
    id             serial primary key,
    swimmerid      integer not null
        references swimmer,
    time           time(2),
    splits         json,
    finapoints     integer,
    additionalinfo varchar,
    penalty        boolean,
    reactiontime   real
);

create table ageclass_to_result
(
    ageclassid integer not null
        references ageclass,
    resultid   integer not null
        references result (id),
    eventid    integer not null
        constraint ageclass_to_result_eventid_fk
            references event (id) on delete cascade,
    primary key (ageclassid, resultid, eventid)
);

create table start
(
    heatid    integer not null
        references heat (id)
            on delete cascade,
    swimmerid integer not null
        references swimmer,
    lane      integer not null,
    time      time(2),
    primary key (heatid, lane)
);


-- Trigger to delete results which are not needed anymore
CREATE OR REPLACE FUNCTION delete_results_on_ageclass_to_result_delete()
    RETURNS TRIGGER AS
$$
BEGIN
    DELETE from result where id = OLD.resultid;
    RETURN OLD;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER on_ageclass_to_result_trigger
    AFTER DELETE ON ageclass_to_result
    FOR EACH ROW
EXECUTE FUNCTION delete_results_on_ageclass_to_result_delete();
