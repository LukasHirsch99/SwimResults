from bs4 import BeautifulSoup
import re
from databaseModels import Meet
from datetime import datetime
# from supabase import create_client, Client
from time import time
import logging
import httpx
import os

logging.getLogger("httpx").setLevel(logging.WARNING)

url: str = "https://qeudknoyuvjztxvgbmou.supabase.co"
key: str = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY2OTQ3NTQyMCwiZXhwIjoxOTg1MDUxNDIwfQ.B3KPV0TWNbrdY0_5n4XLfOehhXA1AFNTZWi23lTTiSU"

# os.environ["SUPABASE_URL"] = url
# os.environ["SUPABASE_KEY"] = key
print(os.environ["SUPABASE_URL"])
print(os.environ["SUPABASE_KEY"])

# supabase: Client = create_client(url, key)

# print(
#     await Supabase.atable("test")
#     .insert({"date": None})
#     .execute()
# )

# print(await Supabase.atable("session").select("*, event(*)", count="exact").eq("meetid", 2021).execute().count)
# await Supabase.atable("session").delete().eq("meetid", 2021).execute()

# print(re.findall(r"\d{2}\.\d{2}\.\d{4}|\d{2}:\d{2}|\d+", "Donnerstag 24.11.2022 - 1. Abschnitt - Einschwimmen 15:00, Beginn 16:30"))

# print(
#     await Supabase.atable("result")
#     .select("*, ageclass!inner(*), swimmer!inner(*, club!inner(*)), event!inner(*, session!inner(*))")
#     .eq("swimmerid", 106200)
#     .eq("event.session.meetid", 2021)
#     .execute()
#     .data[0]
# )
# dbEventCnt = (
#         await Supabase.atable("session")
#         .select("*, event!inner(*)", count="exact")
#         .eq("meetid", 1936)
#         .execute()
#         .count
#     )
# print(dbEventCnt)
# print(supabase.rpc("maxresultid", {}).execute().data)

exit()
s = requests.Session()

r = s.get("https://myresults.eu/de-DE/Meets/Today-Upcoming/2053/Overview")
soup = BeautifulSoup(r.content, "lxml")

divElements = soup.findAll(
    "div", class_=re.compile("^row myresults_content_divtablerow")
)

temp = (
    soup.find("span", class_="myresults_content_divtable_details", string="Ort")
    .parent.contents[3]
    .text.strip()
)
print('"' + str(temp) + '"')
exit()
start = time()
eventCnt = len(
    [
        x
        for x in divElements
        if "myresults_content_divtablerow_header" not in x.attrs["class"]
    ]
)
print(time() - start)
print(eventCnt)

start = time()

div_element: BeautifulSoup = soup.findAll(
    "div",
    class_=re.compile(
        r"row myresults_content_divtablerow myresults_content_divtablerow_[^h].*"
    ),
)

print(time() - start)
print(len(div_element))

# temp = div_element.find("div", class_="col-xs-6").contents[2].text.strip()
