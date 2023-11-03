import 'package:swim_results/model/Heat.dart';
import 'package:swim_results/model/Swimmer.dart';

class Start {
  late int heatId;
  Heat? heat;
  late int swimmerId;
  Swimmer? swimmer;
  late int lane;
  late DateTime? time;

  Start.fromJson(Map json) {
    heatId = json["heatid"];
    if (json.containsKey("heat")) heat = Heat.fromJson(json["heat"]);

    swimmerId = json["swimmerid"];
    if (json.containsKey("swimmer")) {
      swimmer = Swimmer.fromJson(json["swimmer"]);
    }

    lane = json["lane"];
    time = DateTime.tryParse("2000-01-01 ${json["time"]}z");
  }
}
