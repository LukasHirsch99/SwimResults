from datetime import datetime
import re
import asyncio
import time
import httpx
from bs4 import BeautifulSoup
from databaseModels import Meet, Gender, Club
import schedule
# from supabase import create_client, Client
from aiosupabase import Supabase as supabase

import logging

logging.getLogger("httpx").setLevel(logging.WARNING)

client = httpx.AsyncClient()

url: str = "https://qeudknoyuvjztxvgbmou.supabase.co"
key: str = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY2OTQ3NTQyMCwiZXhwIjoxOTg1MDUxNDIwfQ.B3KPV0TWNbrdY0_5n4XLfOehhXA1AFNTZWi23lTTiSU"
# supabase: Client = create_client(url, key)
supabase.configure(url=url, key=key)

meetIdSet = set([d["id"] for d in supabase.table("meet").select("id").execute().data])
swimmerIdSet = set(
    [d["id"] for d in supabase.table("swimmer").select("id").execute().data]
)
clubSet: set[Club] = set(
    [Club.from_dict(d) for d in supabase.table("club").select("*").execute().data]
)

maxHeatId = supabase.rpc("maxheatid", {}).execute().data
maxResultId = supabase.rpc("maxresultid", {}).execute().data

todaysMeets: list[Meet] = None


def parseDate(dateRange: str):
    # 03.10.2020
    # 01.-05.08.2020
    # 29.02.-01.03.2020
    sl = dateRange.split(".")
    month = sl[len(sl) - 2]
    year = sl[len(sl) - 1]

    if dateRange.find("-") >= 0:
        sl = dateRange.split("-")
        if sl[0].count(".") == 2:
            sd = sl[0] + year
            ed = sl[1]
        else:
            sd = sl[0] + month + "." + year
            ed = sl[1]
    else:
        sd = dateRange
        ed = dateRange

    ed = datetime.strptime(ed, "%d.%m.%Y")
    sd = datetime.strptime(sd, "%d.%m.%Y")

    return sd, ed


def getAllMeetInfo(meetId: str):
    r = httpx.get(f"https://myresults.eu/de-DE/Meets/Today-Upcoming/{meetId}/Overview")
    soup = BeautifulSoup(r.content, "lxml")

    if (
        soup.find("span", class_="myresults_content_divtable_details", string="Ort")
        .parent.contents[3]
        .text.strip()
        .find("Austria")
        == -1
    ):
        return None

    m = Meet(None, None, None, None, None, None, None, None, None)

    m.id = meetId

    m.name = soup.find(
        "div",
        class_="row myresults_content_divtablerow myresults_content_divtablerow_header",
    ).text.strip()

    imageTag = soup.find(
        "img",
        src=re.compile("/images/competition_logo/"),
    )

    if imageTag != None:
        m.image = "https://myresults.eu" + imageTag.attrs["src"]
    else:
        m.image = "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ac/No_image_available.svg/300px-No_image_available.svg.png"

    msecmLink = (
        soup.find(
            "span", class_="myresults_content_divtable_details", string="Internet"
        )
        .findPreviousSibling("a")
        .attrs["href"]
        .strip()
    )

    dateString = (
        soup.find("span", class_="myresults_content_divtable_details", string="Datum")
        .parent.contents[0]
        .text.strip()
    )

    m.startdate, m.enddate = parseDate(dateString)

    deadlineString = (
        soup.find(
            "span", class_="myresults_content_divtable_details", string="Meldeschluß"
        )
        .parent.contents[0]
        .text.strip()
    )
    m.deadline = datetime.strptime(deadlineString, "%d.%m.%Y %H:%M")

    m.address = (
        soup.find("span", class_="myresults_content_divtable_details", string="Ort")
        .parent.contents[0]
        .text.strip()
    )

    if msecmLink.find("msecm") == -1:
        m.invitations = []
        m.googlemapslink = None

    else:
        r = httpx.get(msecmLink)
        soup = BeautifulSoup(r.content, "lxml")

        m.googlemapslink = soup.find(
            "a", href=re.compile(r"^http://maps.google.com/maps")
        ).attrs["href"]

        m.invitations = [
            "https://msecm.at" + invLink.attrs["href"]
            for invLink in soup.findAll("a", class_="", href=re.compile(r"^/events/"))
        ]

    return m


def cleanText(str: str) -> str:
    return " ".join(str.split())


def parseTime(t: str):
    if t == "":
        return None

    elif t.find(":") == -1:
        return datetime.strptime(t, "%S.%f").time().isoformat()

    elif t.find("h") != -1:
        return datetime.strptime(t, "%Hh%M:%S.%f").time().isoformat()

    elif t.find(":") != -1:
        return datetime.strptime(t, "%M:%S.%f").time().isoformat()

    return t


def addSwimmerInfoFromStartOrResult(swimmerId, divElement):
    global swimmerIdSet
    name = divElement.findAll("a")[0].text
    clubId = int(divElement.findAll("a")[1].attrs["href"].split("/")[-1])
    birthAndGenderInfo = divElement.findAll(
        "span", class_="myresults_content_divtable_details"
    )[0].text.split(" ")

    if len(birthAndGenderInfo) == 2:
        birthyear = int(birthAndGenderInfo[0])
        gender = (
            Gender.female
            if birthAndGenderInfo[1] == "W"
            else Gender.male
            if birthAndGenderInfo[1] == "M"
            else Gender.mixed
        )
    elif len(birthAndGenderInfo) == 1:
        gender = (
            Gender.female
            if birthAndGenderInfo[0] == "W"
            else Gender.male
            if birthAndGenderInfo[0] == "M"
            else Gender.mixed
        )
        birthyear = -1
    else:
        print(name + "(" + str(swimmerId) + "): Error getting gender or birth info")
        print(divElement.findAll("span", class_="myresults_content_divtable_details"))
        return

    supabase.table("swimmer").insert(
        {
            "id": swimmerId,
            "name": name,
            "birthyear": birthyear,
            "clubid": clubId,
            "gender": gender.value,
        }
    ).execute()
    swimmerIdSet.add(swimmerId)


async def insertStartInfo(meetId: int, startId: int, eventId: int):
    global maxHeatId

    r = await client.get(
        f"https://myresults.eu/de-DE/Meets/Today-Upcoming/{meetId}/Starts/{startId}"
    )
    soup = BeautifulSoup(r.content, "lxml")

    divElements = soup.findAll(
        "div", class_=re.compile("^row myresults_content_divtablerow")
    )

    startCnt = len(
        [
            x
            for x in divElements
            if "myresults_content_divtablerow_header" not in x.attrs["class"]
        ]
    )
    heatCnt = len(divElements) - startCnt

    dbStartCnt = (
        supabase.atable("start")
        .select("*, heat!inner(*)", count="exact")
        .eq("heat.eventid", eventId)
        .execute()
    )

    dbHeatCnt = (
        supabase.atable("heat")
        .select("*", count="exact")
        .eq("eventid", eventId)
        .execute()
    )

    dbStartCnt = (await dbStartCnt).count
    dbHeatCnt = (await dbHeatCnt).count

    if startCnt == dbStartCnt and heatCnt == dbHeatCnt:
        return

    await supabase.atable("heat").delete().eq("eventid", eventId).execute()
    
    heatNr = 1
    heats = []
    starts = []

    print("\t\t\t- Inserting starts")

    for divElement in divElements:
        # Heat-item
        if "myresults_content_divtablerow_header" in divElement["class"]:
            maxHeatId += 1
            heats.append(
                {
                    "id": maxHeatId,
                    "eventid": eventId,
                    "heatnr": heatNr,
                }
            )
            heatNr += 1

        # Start-item
        else:
            swimmerId = int(re.findall(r"\d*$", divElement.find("a").attrs["href"])[0])
            startTime = parseTime(
                divElement.find(
                    "div",
                    class_="hidden-xs col-sm-2 col-md-1 text-right myresults_content_divtable_right",
                ).text
            )

            lane = divElement.find("div", class_="col-xs-1").text

            if swimmerId not in swimmerIdSet:
                addSwimmerInfoFromStartOrResult(swimmerId, divElement)

            starts.append(
                {
                    "heatid": maxHeatId,
                    "swimmerid": swimmerId,
                    "lane": lane,
                    "time": startTime,
                }
            )

    await supabase.atable("heat").insert(heats).execute()
    await supabase.atable("start").insert(starts).execute()


async def insertResultInfo(meetId: int, resultId: int, eventId: int):
    global maxResultId

    r = await client.get(
        f"https://myresults.eu/de-DE/Meets/Today-Upcoming/{meetId}/Results/{resultId}"
    )
    soup = BeautifulSoup(r.content, "lxml")

    divElements = soup.findAll(
        "div", class_=re.compile("^row myresults_content_divtablerow")
    )

    resultCnt = len(
        [
            x
            for x in divElements
            if "myresults_content_divtablerow_header" not in x.attrs["class"]
        ]
    )

    dbResultCnt = (
        supabase.atable("ageclass")
        .select("*, result!inner(*)", count="exact")
        .eq("result.eventid", eventId)
        .execute()
    )

    dbResultCnt = (await dbResultCnt).count

    # Nothing changed since last time
    if resultCnt == dbResultCnt:
        return

    await supabase.atable("result").delete().eq("eventid", eventId).execute()

    print("\t\t\t- Inserting results")

    swimmerIdToDBResultId = {}
    ageClassName: str
    results = []
    ageclasses = []

    for divElement in divElements:
        # Ageclass-item
        if "myresults_content_divtablerow_header" in divElement["class"]:
            ageClassName = divElement.text.strip()

        # Result-item
        else:
            swimmerId = int(divElement.find("a").attrs["href"].split("/")[-1])

            if swimmerId not in swimmerIdSet:
                addSwimmerInfoFromStartOrResult(swimmerId, divElement)

            additionalTimeInfo = divElement.find(
                "div", class_="myresults_content_divtable_points"
            ).text

            timeToFirst = re.search(r"(\+\d+.\d+)", additionalTimeInfo)

            if timeToFirst != None:
                timeToFirst = timeToFirst.group()

            finaPoints = re.search(r"\d+$", additionalTimeInfo)
            if finaPoints != None:
                finaPoints = finaPoints.group()

            if timeToFirst == None and finaPoints == None:
                additionalInfo = additionalTimeInfo
            else:
                additionalInfo = re.search(r"[A-z]+\s", additionalTimeInfo)
                if additionalInfo != None:
                    additionalInfo = additionalInfo.group()

            dbResultId = swimmerIdToDBResultId.get(swimmerId)

            if dbResultId == None:
                dbResultId = (maxResultId := maxResultId + 1)
                swimmerIdToDBResultId[swimmerId] = dbResultId

                splits = divElement.findAll(
                    "span", class_="myresults_content_divtable_details"
                )[3].text

                t = parseTime(
                    divElement.find(
                        "div",
                        class_="hidden-xs col-sm-2 col-md-1 text-right myresults_content_divtable_right",
                    ).text
                )

                results.append(
                    {
                        "id": dbResultId,
                        # PK
                        "swimmerid": swimmerId,
                        "eventid": eventId,
                        "time": t,
                        "splits": splits,
                        "finapoints": finaPoints,
                        "additionalinfo": additionalInfo,
                    }
                )

            position = divElement.find("div", class_="col-xs-1").text
            position = None if position == "" else int(position.replace(".", ""))

            ageclasses.append(
                {
                    "resultid": dbResultId,
                    "name": ageClassName,
                    "position": position,
                    "timetofirst": timeToFirst,
                }
            )

    await supabase.atable("result").insert(results).execute()
    await supabase.atable("ageclass").insert(ageclasses).execute()


def insertClubs(meetId: int):
    global clubSet
    print("\t- Importing clubs")
    r = httpx.get(
        f"https://myresults.eu/de-DE/Meets/Today-Upcoming/{meetId}/Overview/Statistics"
    )

    soup = BeautifulSoup(r.content, "lxml")

    clubs = soup.findAll("div", class_="myresults_content_divtablerow")
    clubs.pop(0)
    clubs.pop(len(clubs) - 1)
    clubs.pop(len(clubs) - 1)
    clubNames = [
        club.find("div", class_="col-xs-11 col-sm-5").text.strip() for club in clubs
    ]
    clubIds = [int(club.find("a").attrs["href"].split("/")[-1]) for club in clubs]

    clubFlags = [
        "https://myresults.eu"
        + club.find("img", class_="myresults_img_16").attrs["src"]
        for club in clubs
    ]

    insertData = []

    for clubId, name, flag in zip(clubIds, clubNames, clubFlags):
        if Club(clubId, name, flag) not in clubSet:
            insertData.append({"id": clubId, "name": name, "nationality": flag})
            clubSet.add(Club(clubId, name, flag))

    supabase.table("club").upsert(insertData).execute()


async def updateSchedule(meetId: int):
    def parseSessionInfo(text: str):
        # Donnerstag 24.11.2022 - 1. Abschnitt - Einschwimmen 15:00, Beginn 16:30
        # Sonntag 07.08.2022 - 7. Abschnitt - Einschwimmen 08:00, Beginn 09:00
        results = re.findall(r"\d{2}\.\d{2}\.\d{4}|\d{2}:\d{2}|\d+", text)

        day = datetime.strptime(results[0], "%d.%m.%Y")
        displayNr: int = results[1]
        if len(results) < 4:
            warmupStart = None
            sessionStart = None
        else:
            warmupStart = results[2] + ":00"
            sessionStart = results[3] + ":00"
        return (day.ctime(), displayNr, warmupStart, sessionStart)

    def parseSessionItemInfo(text: str) -> tuple[int, str]:
        # 40 - 1500m Freistil Damen langsame LÃ¤ufe
        text = text.split(" - ")
        return (int(text[0]), text[1])

    r = await client.get(
        f"https://myresults.eu/de-DE/Meets/Today-Upcoming/{meetId}/Schedule"
    )

    soup = BeautifulSoup(r.content, "lxml")

    divElements = soup.findAll(
        "div", class_=re.compile("^row myresults_content_divtablerow")
    )

    eventCnt = len(
        [
            x
            for x in divElements
            if "myresults_content_divtablerow_header" not in x.attrs["class"]
        ]
    )
    sessionCnt = len(divElements) - eventCnt

    dbEventCnt = (
        supabase.atable("event")
        .select("*, session!inner(*)", count="exact")
        .eq("session.meetid", meetId)
        .execute()
    )

    dbSessionCnt = (
        supabase.atable("session")
        .select("*", count="exact")
        .eq("meetid", meetId)
        .execute()
    )

    dbEventCnt = await dbEventCnt
    dbSessionCnt = await dbSessionCnt
    dbEventCnt = dbEventCnt.count
    dbSessionCnt = dbSessionCnt.count

    if eventCnt != dbEventCnt or sessionCnt != dbSessionCnt:
        supabase.table("session").delete().eq("meetid", meetId).execute().data

    eventDisplayNr: int = None
    sessionDisplayNr: int = None

    tasks = []

    for divElement in divElements:
        # Session-item
        if "myresults_content_divtablerow_header" in divElement["class"]:
            sessionInfo = divElement.text
            day, sessionDisplayNr, warmupStart, sessionStart = parseSessionInfo(
                sessionInfo
            )

            print("\t- Updating session: " + str(sessionDisplayNr))

            sessionId = (
                supabase.table("session")
                .upsert(
                    {
                        "meetid": meetId,
                        "displaynr": sessionDisplayNr,
                        "day": day,
                        "warmupstart": warmupStart,
                        "sessionstart": sessionStart,
                    }
                )
                .execute()
                .data[0]["id"]
            )

        # Event-item
        elif "myresults_content_divtablerow_header" not in divElement["class"]:
            eventInfoString = (
                divElement.find("div", class_="col-xs-6").contents[2].text.strip()
            )

            eventDisplayNr, eventName = parseSessionItemInfo(eventInfoString)

            print("\t\t- Updating event: " + eventName)

            eventId = (
                supabase.table("event")
                .upsert(
                    {
                        "sessionid": sessionId,
                        "name": eventName,
                        "displaynr": eventDisplayNr,
                    }
                )
                .execute()
                .data[0]["id"]
            )

            hrefs = divElement.findAll(
                "a", class_="myresults_content_link myresults_content_divtablecol"
            )

            # Insert starts and/or results
            if len(hrefs) > 0:
                startResultId = int(hrefs[0]["href"].split("/")[-1])

                if len(hrefs) == 2:
                    tasks.append(insertStartInfo(meetId, startResultId, eventId))
                    tasks.append(insertResultInfo(meetId, startResultId, eventId))

                elif hrefs[0]["href"].find("Starts") != -1:
                    tasks.append(insertStartInfo(meetId, startResultId, eventId))

                elif hrefs[0]["href"].find("Results") != -1:
                    tasks.append(insertResultInfo(meetId, startResultId, eventId))

    await asyncio.gather(*tasks)


async def insertUpcomingMeets():
    global meetIdSet

    r = await client.get("https://myresults.eu/de-AT/Meets/Today-Upcoming")
    soup = BeautifulSoup(r.content, "lxml")

    meetIds = [
        int(re.findall(r"[0-9]+", meetLink.attrs["href"])[0])
        for meetLink in soup.findAll("a", href=re.compile(r".*\/Meets\/.+\/Overview"))
    ]

    for meetId in meetIds:
        # Add meet if it doesn't exist in database
        print("Meet: " + str(meetId))
        if meetId not in meetIdSet:
            m = getAllMeetInfo(meetId)

            if m == None:
                continue

            meetIdSet.add(meetId)

            supabase.table("meet").insert(
                {
                    "id": meetId,
                    "name": m.name,
                    "image": m.image,
                    "invitations": m.invitations,
                    "deadline": m.deadline.ctime(),
                    "address": m.address,
                    "googlemapslink": m.googlemapslink,
                    "startdate": m.startdate.ctime(),
                    "enddate": m.enddate.ctime(),
                }
            ).execute()

        # Update all data from meet even if it exists already
        insertClubs(meetId)
        await updateSchedule(meetId)


def setTodaysMeets():
    global todaysMeets
    todaysMeetsDict = (
        supabase.table("meet")
        .select("*")
        .lte("startdate", time.ctime())
        .gte("enddate", time.ctime())
        .order("startdate")
        .execute()
        .data
    )

    todaysMeets = [Meet.from_dict(d) for d in todaysMeetsDict]


async def updateTodaysMeets():
    global todaysMeets

    if todaysMeets == None:
        return

    for meet in todaysMeets:
        print(f"Updating meet: {meet.name}, id: {meet.id}")
        await updateSchedule(meet.id)


if __name__ == "__main__":
    # start = time.time()
    # asyncio.run(updateSchedule(2007))
    # print(time.time() - start)
    # exit()

    # asyncio.run(insertUpcomingMeets())
    setTodaysMeets()
    asyncio.run(updateTodaysMeets())

    schedule.every().day.at("00:01").do(insertUpcomingMeets)
    schedule.every().day.at("00:10").do(setTodaysMeets)
    schedule.every(1).minutes.do(updateTodaysMeets)

    while True:
        schedule.run_pending()
        time.sleep(10)
