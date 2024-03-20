import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/event_screen/components/EventPage.dart';
import 'package:swim_results/model/Heat.dart';
import 'package:swim_results/model/Result.dart';
import 'package:swim_results/model/Session.dart';

class Event {
  late int id;
  late int sessionId;
  Session? session;
  late int displayNr;
  late String name;
  List<Heat> heats = [];
  List<Result> results = [];

  static const TextStyle _infoStyle = TextStyle(
    fontSize: 16,
  );

  Event.fromJson(Map json) {
    id = json["id"];
    sessionId = json["sessionid"];
    displayNr = json["displaynr"];
    name = json["name"];

    if (json.containsKey("heat")) {
      for (var heat in json["heat"]) {
        heats.add(Heat.fromJson(heat));
      }
    }

    if (json.containsKey("result")) {
      for (var result in json["result"]) {
        results.add(Result.fromJson(result));
      }
    }

    if (json.containsKey("session")) {
      session = Session.fromJson(json["session"]);
    }
  }

  static Widget StartItemForSwimmerPage(BuildContext context, Event event) {
    return GestureDetector(
      onTap: () => Navigator.push(
        context,
        MaterialPageRoute(builder: (ctx) => EventPage(event)),
      ),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          color: colorScheme.surfaceTint,
        ),
        child: Column(
          children: [
            Text(
              event.name,
              style: const TextStyle(
                fontSize: 17,
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.format_list_numbered_rounded,
                        color: colorScheme.primary,
                        size: 24,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          "Heat ${event.heats[0].heatNr.toString()}",
                          style: _infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SvgPicture.asset(
                        "assets/StartBlock_cleaned.svg",
                        height: 20,
                        color: colorScheme.primary,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          event.heats[0].starts[0].lane.toString(),
                          style: _infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
                // if (event.heats.first.starts.first.time != null)
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.timer_outlined,
                        size: 24,
                        color: colorScheme.primary,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          formatTime(event.heats[0].starts[0].time),
                          style: _infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
