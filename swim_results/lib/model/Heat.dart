import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Start.dart';

class Heat {
  late int id;
  late int eventId;
  Event? event;
  late int heatNr;
  List<Start> starts = [];

  Heat.fromJson(Map json) {
    id = json["id"];
    eventId = json["eventid"];
    heatNr = json["heatnr"];

    if (json.containsKey("start")) {
      for (var start in json["start"]) {
        starts.add(Start.fromJson(start));
      }
    }

    if (json.containsKey("event")) event = Event.fromJson(json["event"]);
  }
}
