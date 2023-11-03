import 'package:swim_results/model/AgeClass.dart';
import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Swimmer.dart';

class Result {
  late int id;
  late int eventId;
  Event? event;
  late int swimmerId;
  Swimmer? swimmer;
  late DateTime? time;
  late String? splits;
  late int? finaPoints;
  late String? additionalInfo;

  List<AgeClass> ageClasses = [];

  Result.fromJson(Map json) {
    id = json["id"];
    eventId = json["eventid"];
    if (json.containsKey("event")) event = Event.fromJson(json["event"]);

    time = DateTime.tryParse("2000-01-01 ${json["time"]}z");
    splits = json["splits"];
    finaPoints = json["finapoints"];
    additionalInfo = json["additionalinfo"];

    if (json.containsKey("ageclass")) {
      for (var ageclass in json["ageclass"]) {
        AgeClass ac = AgeClass.fromJson(ageclass);
        ac.result = this;
        ageClasses.add(ac);
      }
    }

    swimmerId = json["swimmerid"];
    if (json.containsKey("swimmer")) {
      swimmer = Swimmer.fromJson(json["swimmer"]);
    }
  }
}
