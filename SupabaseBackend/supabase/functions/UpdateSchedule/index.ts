import axiod from "https://deno.land/x/axiod@0.26.2/mod.ts";
import { createClient } from "https://esm.sh/@supabase/supabase-js@2.38.3";
import { DOMParser, Element } from "https://deno.land/x/deno_dom@v0.1.41-alpha-artifacts/deno-dom-wasm.ts";
import { format, parse } from "https://deno.land/std@0.204.0/datetime/mod.ts";
import { corsHeaders } from "../_shared/cors.ts";
import { evaluteRegex } from "../_shared/regexPro.ts";

const supabase = createClient("https://qeudknoyuvjztxvgbmou.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU")

const swimmerIdSet: number[] = []
let maxHeatId: number
let maxResultId: number

enum Gender
{
  M = "M",
  W = "W",
  X = "X"
}

interface SessionInfo
{
  day: string
  displaynr: number
  warmupstart?: string
  sessionstart?: string
}

interface EventInfo
{
  displaynr: number
  name: string
}

interface Swimmer 
{
  id?: number
  firstname: string
  lastname: string
  birthyear: number
  clubid: number
  gender: Gender
  isrelay?: boolean
}

interface Start
{
  heatid: number
  swimmerid: number
  lane: number
  time?: string
}

interface Ageclass
{
  resultid: number
  name: string
  position?: number
  timetofirst?: string
}

interface Result 
{
  id?: number
  eventid: number
  swimmerid: number
  time?: string
  splits?: string
  finapoints?: number
  additionalinfo?: string
}

const parseSessionInfo = (heading: string): SessionInfo =>
{
  heading = heading.trim();

  const regex = /\d{2}\.\d{2}\.\d{4}|\d{2}:\d{2}|\d+/g
  const results = heading.match(regex)

  if (!results)
    throw new Error("Invalid session string format")

  const [day, month, year] = results[0].split(".");
  if (results.length < 4)
  {
    return {
      day: [year, month, day].join("-"),
      displaynr: +results[1],
    }
  }
  return {
    day: [year, month, day].join("-"),
    displaynr: Number.parseInt(results[1]),
    warmupstart: results[2] + ":00",
    sessionstart: results[3] + ":00"
  }
}

const parseEventInfo = (heading: string): EventInfo => 
{
  const l = heading.split(" - ")

  return {
    displaynr: Number.parseInt(l[0]),
    name: l[1]
  }
}

const parseTime = (time: string) =>
{
  if (time.length == 0)
    return undefined

  if (!time.includes(":"))
    return format(parse(time, "ss.SS"), "HH:mm:ss.SS")

  if (time.includes("h"))
  {
    time = time.replace("h", ":")
    return format(parse(time, "H:mm:ss.SS"), "HH:mm:ss.SS")
  }

  return format(parse(time, "m:ss.SS"), "HH:mm:ss.SS")
}

const firstLetterUppercase = (s: string) => s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()

const insertSwimmerFromStartOrResult = async (swimmerId: number, divElement: Element) =>
{
  const hrefs = divElement.getElementsByTagName("a")
  const matches = getAttr(hrefs[1], "href").match(/\d*$/)
  if (!matches) throw new Error("Couldn't parse clubId")

  const swimmer = {} as Swimmer
  swimmer.id = swimmerId
  swimmer.lastname = firstLetterUppercase(hrefs[0].innerText.split(" ")[0])
  swimmer.firstname = firstLetterUppercase(hrefs[1].innerText.split(" ")[0])
  swimmer.clubid = Number.parseInt(matches[0])
  const birthAndGender = divElement.getElementsByClassName("myresults_content_divtable_details")[0].innerText.match(/\d+|[A-Z]/g)
  if (!birthAndGender) throw new Error("No gender or birthyear information found")

  if (birthAndGender.length == 2)
  {
    swimmer.birthyear = Number.parseInt(birthAndGender[0])
    swimmer.gender = birthAndGender[1] as Gender
  }
  else
  {
    swimmer.isrelay = true
    swimmer.gender = birthAndGender[0] as Gender
  }

  await supabase.from("swimmer").insert(swimmer)
  swimmerIdSet.push(swimmerId)
}

const insertStartInfo = async (meetId: number, startId: number, eventId: number) =>
{
  const data = await axiod.get(`https://myresults.eu/de-DE/Meets/Today-Upcoming/${meetId}/Starts/${startId}`);
  const soup = new DOMParser().parseFromString(data.data, "text/html");
  const divElements = soup?.getElementsByClassName("myresults_content_divtablerow");
  if (!divElements)
    return;

  const startCnt = divElements.filter((element) => !element.classList.contains("myresults_content_divtablerow_header")).length
  const heatCnt = divElements?.length - startCnt;

  const [dbStarts, dbHeats] = await Promise.all([
    supabase.from("start").select("*, heat!inner(*)", { count: "exact" }).eq("heat.eventid", eventId),
    supabase.from("heat").select("*", { count: "exact" }).eq("eventid", eventId)
  ]);

  if (startCnt == dbStarts.count && heatCnt == dbHeats.count)
    return;

  console.log("\t\t\t- Inserting Starts");

  await supabase.from("heat").delete().eq("eventid", eventId)

  let heatNr = 0
  const heats = []
  const starts: Start[] = []


  for (const divElement of divElements)
  {
    // Heat-Item
    if (divElement.classList.contains("myresults_content_divtablerow_header"))
    {
      maxHeatId++
      heatNr++
      heats.push({
        "id": maxHeatId,
        "eventid": eventId,
        "heatnr": heatNr
      })
    }
    // Start-Item
    else
    {
      const matches = getAttr(divElement.getElementsByTagName("a")[0], "href").match(/\d*$/)
      if (!matches)
        throw new Error("Couldn't parse swimmerId")
      const swimmerId = Number.parseInt(matches[0])
      const startTime = parseTime(divElement.getElementsByClassName("hidden-xs col-sm-2 col-md-1 text-right myresults_content_divtable_right")[0].innerText)
      const lane = Number.parseInt(divElement.getElementsByClassName("col-xs-1")[0].innerText)

      if (!swimmerIdSet.includes(swimmerId))
        await insertSwimmerFromStartOrResult(swimmerId, divElement)

      starts.push({
        heatid: maxHeatId,
        swimmerid: swimmerId,
        lane: lane,
        time: startTime,
      })
    }
  }

  await supabase.from("heat").insert(heats)
  await supabase.from("start").insert(starts)
}

const insertResultInfo = async (meetId: number, msecmResultId: number, eventId: number) =>
{
  const data = await axiod.get(`https://myresults.eu/de-DE/Meets/Today-Upcoming/${meetId}/Results/${msecmResultId}`);
  const soup = new DOMParser().parseFromString(data.data, "text/html");
  const divElements = soup?.getElementsByClassName("myresults_content_divtablerow");
  if (!divElements)
    return;

  const resultCnt = divElements.filter((element) => !element.classList.contains("myresults_content_divtablerow_header")).length

  const { count: dbResultCnt } = await supabase.from("ageclass").select("*, result!inner(*)", { count: "exact" }).eq("result.eventid", eventId)

  if (resultCnt == dbResultCnt)
    return;

  console.log("\t\t\t- Inserting Results");

  await supabase.from("result").delete().eq("eventid", eventId)

  const swimmerIdToDbResultId = new Map()
  let ageClassName = ""
  const results: Result[] = []
  const ageclasses: Ageclass[] = []


  for (const divElement of divElements)
  {
    // Ageclass-Item
    if (divElement.classList.contains("myresults_content_divtablerow_header"))
    {
      ageClassName = divElement.innerText.trim()
    }
    // Result-Item
    else
    {
      const matches = getAttr(divElement.getElementsByTagName("a")[0], "href").match(/\d*$/)
      if (!matches)
        throw new Error("Couldn't parse swimmerId")
      const swimmerId = Number.parseInt(matches[0])
      if (!swimmerIdSet.includes(swimmerId))
        await insertSwimmerFromStartOrResult(swimmerId, divElement)

      const resultInfoString = divElement.getElementsByClassName("myresults_content_divtable_points")[0].innerText.trim()
      const resultInfo = evaluteRegex(/(?<timeToFirst>\+\d+\.\d+)|(?<finaPoints>\d+)|(?<additionalInfo>[\S]+)/g, resultInfoString)
      // if (resultInfo == undefined)
      // throw new Error(`Couldn't parse resultInfo: "${resultInfoString}", swimmerId="${swimmerId}", msecmResultId=${msecmResultId}`)
      let dbResultId = swimmerIdToDbResultId.get(swimmerId)

      if (dbResultId == undefined)
      {
        dbResultId = ++maxResultId
        swimmerIdToDbResultId.set(swimmerId, dbResultId)

        const finaPointsString = resultInfo.get("finaPoints")
        const finaPoints = finaPointsString ? Number.parseInt(finaPointsString) : undefined
        const additionalInfo = resultInfo.get("additionalInfo")
        const splits = divElement.getElementsByClassName("myresults_content_divtable_details")[3].innerText.trim()
        const time = parseTime(divElement.getElementsByClassName("hidden-xs col-sm-2 col-md-1 text-right myresults_content_divtable_right")[0].innerText)

        results.push(
          {
            id: dbResultId,
            swimmerid: swimmerId,
            eventid: eventId,
            time: time,
            splits: splits,
            finapoints: finaPoints,
            additionalinfo: additionalInfo
          }
        )
      }
      const timeToFirst = resultInfo.get("timeToFirst")
      const position = Number.parseInt(divElement.getElementsByClassName("col-xs-1")[0].innerText)
      ageclasses.push({
        name: ageClassName,
        resultid: dbResultId,
        position: position,
        timetofirst: timeToFirst
      })
    }
  }

  let { error } = await supabase.from("result").insert(results)
  if (error)
  {
    throw error
  }

  ({ error } = await supabase.from("ageclass").insert(ageclasses))
  if (error)
    throw error
}

const updateSchedule = async (meetId: number) =>
{
  const data = await axiod.get("https://myresults.eu/de-DE/Meets/Today-Upcoming/" + meetId + "/Schedule");
  const soup = new DOMParser().parseFromString(data.data, "text/html");
  const divElements = soup?.getElementsByClassName("myresults_content_divtablerow");
  if (!divElements)
    return;

  const eventCnt = divElements.filter((event) => !event.classList.contains("myresults_content_divtablerow_header")).length
  const sessionCnt = divElements?.length - eventCnt;

  const [dbEvents, dbSessions] = await Promise.all([
    supabase.from("event").select("*, session!inner(*)", { count: "exact" }).eq("session.meetid", meetId),
    supabase.from("session").select("*", { count: "exact" }).eq("meetid", meetId)
  ]);

  const scheduleUpToDate = eventCnt == dbEvents.count && sessionCnt == dbSessions.count

  if (!scheduleUpToDate)
    await supabase.from("session").delete().eq("meetid", meetId)

  let sessionId: number | null = null;
  const tasks = [];

  for (const divElement of divElements)
  {
    // Session-Item
    if (divElement.classList.contains("myresults_content_divtablerow_header"))
    {
      const sessionInfo = parseSessionInfo(divElement.innerText);

      if (scheduleUpToDate)
        sessionId = dbSessions.data?.find((session) => session["displaynr"] == sessionInfo.displaynr && session["day"] == sessionInfo.day)["id"];

      else
      {
        const { data, error } = await supabase.from("session").upsert({
          "meetid": meetId,
          ...sessionInfo
        }).select();

        if (!data || error) throw error;

        sessionId = data[0]["id"];
      }
    }

    else
    {
      const eventInfo = parseEventInfo(divElement.getElementsByClassName("col-xs-6")[0].childNodes[2].textContent.trim());
      if (!eventInfo.displaynr || !eventInfo.name)
        continue

      let eventId: number;

      if (scheduleUpToDate)
      {
        eventId = dbEvents.data?.find((event) => event["sessionid"] == sessionId && event["name"] == eventInfo.name && event["displaynr"] == eventInfo.displaynr)["id"];
      }
      else
      {
        const { data, error } = await supabase.from("event").upsert({
          "sessionid": sessionId,
          ...eventInfo
        }).select();

        if (!data || error) throw error;

        eventId = data[0]["id"];
      }

      const hrefs = divElement.getElementsByClassName("myresults_content_link myresults_content_divtablecol")

      if (!hrefs || !hrefs.length)
        continue;

      const matches = /\d*$/g.exec(getAttr(hrefs[0], "href"))
      if (!matches)
        continue

      const startResultId = Number.parseInt(matches[0])

      if (hrefs.length >= 1)
        tasks.push(insertStartInfo(meetId, startResultId, eventId));
      if (hrefs.length == 2)
        tasks.push(insertResultInfo(meetId, startResultId, eventId));
    }
  }
  await Promise.all(tasks)
}

const getAttr = (element: Element, attribute: string) =>
{
  const v = element.attributes.getNamedItem(attribute)?.value
  if (!v)
    throw new Error(`Couldn't get attribute ${attribute} from ${element.attributes}`)
  return v
}

// await updateSchedule(2021)

Deno.serve(async (req) =>
{
  const { data: maxHeatIdData } = await supabase.rpc("maxheatid").select()
  const { data: maxResultIdData } = await supabase.rpc("maxresultid").select()
  const { data: swimmerIdSetData } = await supabase.from("swimmer").select("id")

  if (maxHeatIdData == null)
    throw new Error("Couldn't get maxHeatId")
  else if (maxResultIdData == null)
    throw new Error("Couldn't get maxResultId")
  else if (swimmerIdSetData == null)
    throw new Error("Couldn't get swimmerIds")

  maxHeatId = Number.parseInt(maxHeatIdData.toString())
  maxResultId = Number.parseInt(maxResultIdData.toString())

  swimmerIdSetData.forEach((value) => swimmerIdSet.push(value["id"]))

  const { meetId } = await req.json()

  const start = new Date().getTime();
  await updateSchedule(meetId);

  const data = {
    message: `Updated meet: ${meetId}`,
    time: `${(new Date().getTime() - start)}ms`
  }

  return new Response(
    JSON.stringify(data),
    { headers: { "Content-Type": "application/json", ...corsHeaders } },
  )
})
