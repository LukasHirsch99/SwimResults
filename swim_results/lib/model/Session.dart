import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Meet.dart';

class Session {
  late int id;
  late int meetId;
  Meet? meet;
  late DateTime day;
  late DateTime? warmupStart;
  late DateTime? sessionStart;
  late int displayNr;
  List<Event> events = [];

  static final DateFormat _dayFormatter = DateFormat("EEEE dd.MM.yyyy");
  static final DateFormat _startFormatter = DateFormat("HH:mm");

  Session.fromJson(Map json) {
    id = json["id"];
    meetId = json["meetid"];
    day = DateTime.parse(json["day"]);
    warmupStart = DateTime.parse(json["day"] + " " + json["warmupstart"]);
    sessionStart = DateTime.parse(json["day"] + " " + json["sessionstart"]);
    displayNr = json["displaynr"];

    if (json.containsKey("event")) {
      for (var event in json["event"]) {
        Event e = Event.fromJson(event);
        e.session = this;
        events.add(e);
      }
    }

    if (json.containsKey("meet")) meet = Meet.fromJson(json["meet"]);
  }

  static Widget SessionItemContainer({
    required Session session,
    required Widget Function(BuildContext, Event) eventWidget,
  }) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Column(
        children: [
          Text(
            _dayFormatter.format(session.day),
            style: TextStyle(
              color: colorScheme.onSecondary,
              fontSize: 20,
              fontWeight: FontWeight.w400,
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              Text(
                "${_startFormatter.format(session.warmupStart!)} Warmup",
                style: TextStyle(
                    fontSize: 17,
                    color: colorScheme.onSecondary,
                    fontWeight: FontWeight.w300),
              ),
              Text(
                "${_startFormatter.format(session.sessionStart!)} Start",
                style: TextStyle(
                    fontSize: 17,
                    color: colorScheme.onSecondary,
                    fontWeight: FontWeight.w300),
              )
            ],
          ),
          ListView.builder(
            shrinkWrap: true,
            itemCount: session.events.length,
            padding: const EdgeInsets.all(0),
            physics: const NeverScrollableScrollPhysics(),
            itemBuilder: (BuildContext context, int eventIdx) {
              return Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 10,
                  vertical: 5,
                ),
                child: eventWidget(
                  context,
                  session.events[eventIdx],
                ),
              );
            },
          ),
        ],
      ),
    );
  }
}
