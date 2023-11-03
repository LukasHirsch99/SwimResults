import 'package:swim_results/components/globals.dart';
import 'package:swim_results/model/Club.dart';

class Swimmer {
  late int id;
  late String name;
  int? birthyear;
  int? clubId;
  Club? club;
  String? gender;

  Swimmer.fromJson(Map json) {
    id = json["id"];
    name = json["name"];
    birthyear = json["birthyear"];
    clubId = json["clubid"];
    if (json.containsKey("club")) club = Club.fromJson(json["club"]);
    gender = json["gender"];
  }

  Swimmer.fromSwimmerSearchPage(Map json) {
    id = json["id"];
    name = json["name"];
    birthyear = json["birthyear"];
    clubId = json["clubid"];
    gender = json["gender"];
    club = Club.fromSwimmerSearchPage(json);
  }
}
