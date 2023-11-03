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

  static final TextStyle _infoStyle = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w300,
    color: colorScheme.onPrimary,
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
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.primary,
      ),
      child: GestureDetector(
        onTap: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (ctx) => EventPage(event)),
        ),
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.only(top: 5),
              child: Text(
                event.name,
                style: TextStyle(
                  color: colorScheme.onPrimary,
                  fontSize: 17,
                  fontWeight: FontWeight.w300,
                ),
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
                        Icons.format_list_bulleted_rounded,
                        color: colorScheme.onPrimary,
                        size: 25,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          event.heats[0].heatNr.toString(),
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
                        color: colorScheme.onPrimary,
                        height: 22,
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
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.timer_outlined,
                        color: colorScheme.onPrimary,
                        size: 25,
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
