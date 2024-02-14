import axiod from "https://deno.land/x/axiod@0.26.2/mod.ts";
import { createClient } from "https://esm.sh/@supabase/supabase-js@2.38.3";
import {
  DOMParser,
  Element,
} from "https://deno.land/x/deno_dom@v0.1.41-alpha-artifacts/deno-dom-wasm.ts";
import { format, parse } from "https://deno.land/std@0.204.0/datetime/mod.ts";
import { corsHeaders } from "../_shared/cors.ts";

const supabase = createClient(
  "https://qeudknoyuvjztxvgbmou.supabase.co",
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU",
);

const clubIdSet: number[] = [];
const meetIdSet: number[] = [];

interface MeetDate {
  startdate: string;
  enddate: string;
}

interface Club {
  id: number;
  name: string;
  nationality?: string;
}

interface Meet {
  id: number;
  name: string;
  image?: string;
  invitations?: string[];
  deadline: string;
  address: string;
  startdate: string;
  enddate: string;
  googlemapslink?: string;
}

const getAttr = (element: Element, attribute: string) => {
  const v = element.attributes.getNamedItem(attribute)?.value;
  if (!v) {
    throw new Error(
      `Couldn't get attribute ${attribute} from ${element.attributes}`,
    );
  }
  return v;
};

const parseDate = (s: string): MeetDate => {
  // 01.-05.08.2020
  let m = s.match(
    /^(?<firstDay>\d{2}).-(?<lastDay>\d{2}).(?<month>\d{2}).(?<year>\d{4})/,
  )?.groups;

  if (m) {
    return {
      startdate: `${m["year"]}-${m["month"]}-${m["firstDay"]}`,
      enddate: `${m["year"]}-${m["month"]}-${m["lastDay"]}`,
    };
  }

  // 03.10.2020
  m = s.match(/^(?<day>\d{2}).(?<month>\d{2}).(?<year>\d{4})/)?.groups;
  if (m) {
    return {
      startdate: `${m["year"]}-${m["month"]}-${m["day"]}`,
      enddate: `${m["year"]}-${m["month"]}-${m["day"]}`,
    };
  }

  // 29.02.-01.03.2020
  m = s.match(
    /^(?<firstDay>\d{2}).(?<firstMonth>\d{2}).-(?<lastDay>\d{2}).(?<lastMonth>\d{2}).(?<year>\d{4})/,
  )?.groups;
  if (m) {
    return {
      startdate: `${m["year"]}-${m["firstMonth"]}-${m["firstDay"]}`,
      enddate: `${m["year"]}-${m["lastMonth"]}-${m["lastDay"]}`,
    };
  }

  throw new Error("Couldn't parse date");
};

const insertClubs = async (meetId: number) => {
  console.log("\t- Importing clubs");

  const data = await axiod.get(
    `https://myresults.eu/de-DE/Meets/Today-Upcoming/${meetId}/Overview/Statistics`,
  );
  const soup = new DOMParser().parseFromString(data.data, "text/html");
  const clubElements = soup?.getElementsByClassName(
    "myresults_content_divtablerow",
  );
  if (!clubElements) return;

  clubElements.splice(0, 1);
  clubElements.splice(clubElements.length - 2, 2);

  const clubs: Club[] = [];
  for (const clubElement of clubElements) {
    const matches = clubElement
      .getElementsByTagName("a")[0]
      .attributes.getNamedItem("href")
      ?.value.match(/\d*$/);
    if (!matches) throw new Error("Couldn't parse swimmerId");
    const clubId = Number.parseInt(matches[0]);

    if (clubIdSet.includes(clubId)) continue;
    clubIdSet.push(clubId);

    const name = clubElement.getElementsByTagName("a")[0].innerText.trim();
    const flagLink = clubElement
      .getElementsByTagName("img")[0]
      .attributes.getNamedItem("src")?.value;
    const nationality = flagLink
      ? `https://myresults.eu${flagLink}`
      : undefined;
    clubs.push({
      id: clubId,
      name: name,
      nationality: nationality,
    });
  }
  await supabase.from("club").insert(clubs);
};

const getAllMeetInfo = async (meetId: number): Promise<Meet | null> => {
  let data = await axiod.get(
    `https://myresults.eu/de-DE/Meets/Today-Upcoming/${meetId}/Overview`,
  );
  let soup = new DOMParser().parseFromString(data.data, "text/html");
  if (!soup) return null;

  const details = soup.getElementsByClassName(
    "row myresults_content_divtablerow",
  );
  if (!details) return null;

  if (!details[6].innerText.includes("Austria")) return null;

  const name = soup
    .getElementsByClassName(
      "row myresults_content_divtablerow myresults_content_divtablerow_header",
    )[0]
    .innerText.trim();
  const imgTag = soup.getElementsByClassName("img-responsive center-block")[0];
  const image = imgTag
    ? "https://myresults.eu" + getAttr(imgTag, "src")
    : undefined;
  const date = parseDate(
    details[3].childNodes[1].childNodes[0].textContent.trim(),
  );
  const deadline = format(
    parse(
      details[4].childNodes[1].childNodes[0].textContent.trim(),
      "dd.MM.yyyy HH:mm",
    ),
    "yyyy-MM-dd HH:mm:ss",
  );
  const address = details[6].childNodes[1].childNodes[0].textContent.trim();

  const msecmLink = details[13].childNodes[1].childNodes[1].textContent.trim();
  let invitations: string[] = [];
  let googlemapslink: string | undefined;

  if (!msecmLink.includes("msecm")) invitations = [];
  else {
    data = await axiod.get(msecmLink);
    soup = new DOMParser().parseFromString(data.data, "text/html");
    googlemapslink = soup
      ?.getElementsByClassName("text-right")[0]
      ?.getElementsByTagName("a")[0]
      .attributes.getNamedItem("href")?.value;
    const invLinks = soup
      ?.getElementsByTagName("a")
      .filter((v) => v.attributes.getNamedItem("href")?.value.includes(".pdf"));
    invLinks?.forEach((v) => invitations.push(getAttr(v, "href")));
  }

  return {
    id: meetId,
    name: name,
    deadline: deadline,
    address: address,
    startdate: date.startdate,
    enddate: date.enddate,
    image: image,
    invitations: invitations,
    googlemapslink: googlemapslink,
  };
};

const insertUpcomingMeets = async (updateRecent = false) => {
  if (updateRecent) console.log("Updating Recent Meets");
  else console.log("Updating Upcoming Meets");

  const data = await axiod.get(
    `https://myresults.eu/de-DE/Meets/${
      updateRecent ? "Recent" : "Today-Upcoming"
    }`,
  );
  const soup = new DOMParser().parseFromString(data.data, "text/html");

  const meetLinks = soup
    ?.getElementsByTagName("a")
    .filter((v) => getAttr(v, "href").includes("Overview"));
  if (!meetLinks) throw new Error("Didn't find any new upcoming meets");

  for (const meetLink of meetLinks) {
    const matches = /\d+/g.exec(getAttr(meetLink, "href"));
    if (!matches) throw new Error("Couldn't parse meetId");

    const meetId = Number.parseInt(matches[0]);

    const meet = await getAllMeetInfo(meetId);
    if (!meet) continue;

    if (!meetIdSet.includes(meetId)) {
      meetIdSet.push(meetId);
      await supabase.from("meet").insert(meet);
    } else continue;

    console.log(`Inserting new meet ${meetId}`);
    await insertClubs(meetId);
    console.log(
      (
        await supabase.functions.invoke("UpdateSchedule", {
          body: { meetId: meetId },
        })
      ).error,
    );
  }
};

// await insertUpcomingMeets()

Deno.serve(async (req: Request) => {
  const { data: meetIdSetData } = await supabase.from("meet").select("id");
  const { data: clubIdSetData } = await supabase.from("club").select("id");

  if (clubIdSetData == null) throw new Error("Couldn't get clubIdSet");
  else if (meetIdSetData == null) throw new Error("Couldn't get meetIdSet");

  meetIdSetData.forEach((value) => meetIdSet.push(value["id"]));
  clubIdSetData.forEach((value) => clubIdSet.push(value["id"]));

  const { updateRecent } = await req.json();

  await insertUpcomingMeets(updateRecent);
  const data = {
    message: updateRecent ? `Updated recent Meets` : `Updated new Meets`,
  };

  return new Response(JSON.stringify(data), {
    headers: { "Content-Type": "application/json", ...corsHeaders },
  });
});

