
import 'package:swim_results/model/Result.dart';

class AgeClass {
  late int resultId;
  Result? result;
  late String name;
  late int? position;
  late String? timeToFirst;

  AgeClass.fromJson(Map json) {
    resultId = json["resultid"];
    if (json.containsKey("result")) {
      result = Result.fromJson(json["result"]);
    }
    name = json["name"];
    position = json["position"];
    timeToFirst = json["timetofirst"];
  }
}